```
{
    "url": "codeigniter-mysql-proxy",
    "time": "2014/03/29 22:11",
    "tag": "PHP,CodeIgniter"
}
```

# 一、目标
当前服务器只做了主从，未配置读写分离，读写分离的功能就只有交给程序来实现，本文主要谈谈Codeigniter怎么实现读写分离，并且需要满足以下两点：

**1、读写分离对开发应该透明。**

网上有方案通过手动load多个DB来实现读写分离，这样的分离跟业务关联太紧，增加了开发难度也不利于维护，我们要做的是默认读重库，写则写主库，读写分离对开发者透明

**2、配置简单。**

保留现有的配置方式，通过增加一个数组来配置读写分离，不影响原有使用方式。

# 二、实现思路
- 1、要实现读写分离最简单的思路就是在最终执行查询的地方根据查询语句判断是插入主库还是读取从库，所以需要找到该函数。
- 2、应该只连接一次数据库，下次操作该链接应当可复用。也就是连一次重库后所有的读操作都可用，不需再次连接，主库同理。所以我们可以将链接放在CI超级对象中。
- 3、主从的判断是根据最终执行的SQL语句来判断的，所以数据库配置中的自动链接autoinit参数就不用设置为true了，如果默认连接了而又不需要操作该库就浪费资源了。
- 4、模型中可以使用$this->db来直接操作查询，不需要其他调整。
- 5、不直接修改system下的文件

# 三、实现读写分离
CI的DB类固定为读取system下的文件，我们可以通过适当的重写来实现。首先是Loader.php，其中的database方法用来加载数据库对象，固定引用了system/database/DB.php文件，我们判断下是否存在自定义DB.php文件，存在则引入。

## 3.1 重写Loader.php
```
public function database($params = '', $return = FALSE, $active_record = NULL)
{
    $CI =& get_instance();
    if (class_exists('CI_DB') AND $return == FALSE AND $active_record == NULL AND isset($CI->db) AND is_object($CI->db)) {
        return FALSE;
    }
    if(file_exists(APPPATH.'core/database/DB.php')) {
        require_once(APPPATH.'core/database/DB.php');
    } else {
        require_once(BASEPATH.'database/DB.php');
    }
    if ($return === TRUE) {
        return DB($params, $active_record);
    }
    $CI->db = '';
    $CI->db =& DB($params, $active_record);
}
/* End of file MY_Loader.php */
/* Location: ./application/core/MY_Loader.php */
```
接着我们在application/core下创建database/DB.php，该文件只有一个DB方法，用来读取配置文件并进行初始化工作。同样有两处地方需要重写下：

## 3.2 重写DB.php
```
//DB_driver.php为所有驱动方式的父类，最终执行查询的方法在该文件中
//第一处修改为判断自定义的DB_driver.php是否存在，存在则引入
if(file_exists(APPPATH.'core/database/DB_driver.php')) {
    require_once(APPPATH.'core/database/DB_driver.php');
} else {
    require_once(BASEPATH.'database/DB_driver.php');
}
  
//第二处 $params['dbdriver'].'_driver.php' 该文件可不调整，实际未修改该文件，为了方便调试也加了
//mysql驱动对应system/database/drivers/mysql/mysql_driver.php，mysql的最后执行方法在这里，
//包括数据库打开和关闭、查询等,可以该文件增加相应日志查看读写分离是否有效
if(file_exists(APPPATH.'core/database/drivers/'.$params['dbdriver'].'/'.$params['dbdriver'].'_driver.php')) {
    require_once(APPPATH.'core/database/drivers/'.$params['dbdriver'].'/'.$params['dbdriver'].'_driver.php');
} else {
    require_once(BASEPATH.'database/drivers/'.$params['dbdriver'].'/'.$params['dbdriver'].'_driver.php');
}
//将当前group name赋值给param，方便判断
$params['group_name'] = $active_group;
  
/* End of file DB.php */
/* Location: ./application/core/database/DB.php */
```
整个DB.php调整的也基本上是文件的引入，group name的引入是为了方便后面的判断， 不引入则可以通过主机、数据库名称这些来配置。如果想强制关闭autoint，可以在DB.php中删掉下面这段：
```
if ($DB->autoinit == TRUE)
{
    $DB->initialize();
}
```
接下来就是最核心的地方。根据查询语句实现读写分离。DB_driver.php中的simple_query方法可以理解为最后执行SQL语句的方法，我们可以在这里进行数据库链接的判断。

## 3.3 重写DB_driver.php
```
//增加属性，表示当前组
var $active_group;
  
//增加属性，使用强制使用主库
var $db_force_master;
  
//该方法为执行查询的必经之地，我们可以在这里根据类型判断使用哪个链接。
function simple_query($sql)
{
    //load_db_proxy_setting方法这里写在helper中，也可以直接写在该类中，写在helper中则需要在自动加载中加载该helper
    //该方法的作用是根据当前链接group name 和sql读写类型，以及是否强制使用主库判断使用哪个链接。使用主库 OR 重库？
    //主重库的负载均衡，单点故障都可以在这里考虑。也就是根据3个参数返回一个可用的配置数组。
    $proxy_setting = load_db_proxy_setting($this->group_name, $this->is_write_type($sql), $this->db_force_master);
    if(is_array($proxy_setting) && ! empty($proxy_setting)) {
        $proxy_setting_key = key($proxy_setting);
        $this->group_name = $proxy_setting_key;
        //将当前配置重新赋值给类的属性，如果database.php配置的是DSN字符串，则需要在load_db_proxy_setting中做处理
        foreach($proxy_setting[$proxy_setting_key] as $key => $val) {
            $this->$key = $val;
        }
        //定义链接ID为conn_前缀
        $proxy_conn_id = 'conn_'.$proxy_setting_key;
        $CI = & get_instance();
        //赋值给CI超级对象或者直接从CI超级对象中读取
        if(isset($CI->$proxy_conn_id) && is_resource($CI->$proxy_conn_id)) {
            $this->conn_id = $CI->$proxy_conn_id;
        } else {
            $this->conn_id = false;
            $this->initialize();
            $CI->$proxy_conn_id = $this->conn_id;
        }
        //强制只一次有效，下次查询失效，防止一直强制主库
        $this->reset_force_master();
    }
    if ( ! $this->conn_id)
    {
        $this->initialize();
    }
    return $this->_execute($sql);
}
//某些情况会强制使用主库，先执行该方法即可
public function force_master()
{
    $this->db_force_master = TRUE;
}
public function reset_force_master()
{
    $this->db_force_master = FALSE;
}
/* End of file DB_driver.php */
/* Location: ./application/core/database/DB_driver.php */
```
到这里读写分离即基本实现了，但做事情得善始善终，链接的数据库对象需要关闭，可以在公用控制器中执行完毕后关掉连接。DB_driver.php中也有close方法，可以考虑下是否可以在该方法中关闭？这里认为是不行的。

## 3.4 关闭数据库链接
```
class MY_Controller extends CI_Controller
{
    public function __construct()
    {
        parent::__construct();
        $this->load->service('common/helper_service', NULL, 'helper');
        //下面这段为关闭CI超级对象中的数据库对象和数据库链接，db的对象Codeigniter.php中会关闭
        register_shutdown_function(function(){
            foreach(get_object_vars($this) as $key => $val) {
                if(substr($key, 0, 3) == 'db_' && is_object($this->{$key}) && method_exists($this->{$key}, 'close')) {
                    $this->{$key}->close();
                }
                if(substr($key, 0, 5) == 'conn_'  && is_resource($this->{$key})) {
                    $this->db->_close($val);
                    unset($this->{$key});
                }
            }
        });
    }
}
/* End of file MY_Controller.php */
/* Location: ./application/core/MY_Controller.php */
```
模型中的使用，为了使每个model中都可使用$this->db，以及不多次连接数据库，这里也是将链接放在CI超级对象中。这里就算不读写分离也可以这么处理，可以很方便的连接多个DB，具体的model要使用其他库只需要在构造函数中传入group name即可。

## 3.5 模型调整
```
public function __construct($group_name = '')
{
    parent::__construct();
    $this->initDb($group_name);
}
private function initDb($group_name = '')
{
    $db_conn_name = $this->getDbName($group_name);
    $CI = & get_instance();
    if(isset($CI->{$db_conn_name}) && is_object($CI->{$db_conn_name})) {
        $this->db = $CI->{$db_conn_name};
    } else {
        $CI->{$db_conn_name} = $this->db = $this->load->database($group_name, TRUE);
    }
}
private function getDbName($group_name = '')
{
    if($group_name == '') {
        $db_conn_name = 'db';
    } else {
        $db_conn_name = 'db_'.$group_name;
    }
    return $db_conn_name;
}
/* End of file MY_Model.php */
/* Location: ./application/core/MY_Model.php */
```
最后的数据库配置方式，只需要在原有的基础上配置一个数组即可。是使用双主还是一主多从就看这里的配置方式。最开始想到直接在原配置上加键名来处理，但主与从的对应关系还是没有这样子明了，这里的定义方式决定了load_db_proxy_setting的实现方式。

## 3.6 database.php配置
```
$_master_slave_relation = array(
    'default_master' => array('default_slave1', 'default_slave2', 'default_slave3'),
);
/* End of file database.php */
/* Location: ./application/config/database.php */
```
最开始的数据库链接并未放到CI超级对象中，发现load多个模型时每次都会打开链接，所以完成读写分离之后一定要测试，可以在数据库链接打开和关闭的地方查看是否按预期执行(方法对应application/core/database/drivers/mysql/mysql_driver.php中的db_connect和_close)。整个调整过程最重要的两点就是simple_query方法以及构造函数中关闭数据库链接。模型中的调整是为了更方便的链接多个库，未实现读写分离时也是这么调整的，常用的方法独立成一个文件，MY_Model去继承。

实现MYSQL读写分离的中间件挺多，在没有用到这些时可以通过程序上的控制来实现读写分离。当然这里只是实现了读写分离，可以强制使用主库。如果想要更好的分配方式，可以好好想想load_db_proxy_setting中的分配方式。