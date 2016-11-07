```
{
    "url": "codeigniter-mvc",
    "time": "2013/11/01 16:39",
    "tag": "PHP,CodeIgniter"
}
```

# 一、控制器
## 1.1 公用控制器
Codeigniter所有的控制器都必须继承CI_Controller类，但CI_Controller类位于system目录下，不太方便修改。为方便做一些公用的处理，通常情况下我们会在core下创建MY_Controller，用来继承CI_Controller，从而项目中所有的控制器继承MY_Controller。

**那么，MY_Controller 通常会做些什么呢？**

所有的控制器都继承了MY_Controller， MY_Controller常常会加载一些公用帮助函数、公用类库，以及实现一些公用的方法。公用的方法？公有的方法？

看到这些方法会意识到一个问题，如果方法是public的，那是否可以通过浏览器访问到。答案是可以的！这样不该让用户访问到的方法让用户访问到了。那设置protected吧。。。

**备注**:CI_Controller中写public方法不会被访问到，框架限制了CI_Controller中方法通过浏览器访问。

随着项目的不断进展，MY_Controller中的公用方法会越来越多。如果此时要增加后台管理的功能，所有的控制器依然继承MY_Controller，那其中的很多方法可能不适用了。如果后台需要的一些公用方法也写在这里，这里将会变得混乱。

**如何按模块区分不同的控制器？**

有两种处理的方式，第一种是通过不同的公用控制器文件来区分，由控制器去决定继承哪一个公用控制器，当然这里得引入公用文件。还有这种方式是可以通过对象的一个属性来维护，不同的模块赋予该属性不同的对象。如：
```
<?php
if ( ! defined('BASEPATH'))
    exit('No direct script access allowed');
  
class MY_Controller extends CI_Controller
{
  public function __construct($type = NULL)
  {
      parent::__construct();
      switch($type) {
        case 'api' :
          $this->load->library('api_helper', NULL, 'helper');
          break;
        case 'admin' :
          $this->load->library('admin_helper', NULL, 'helper');
            break;
        default :
          $this->load->library('app_helper', NULL, 'helper');
            break
      }
  }
}
  
/* End of file MY_Controller.php */
/* Location: ./application/core/MY_Controller.php */
```
控制器调用MY_Controller构造函数并传入type值，根据不同的type值会加载不同的类库，然后给类定义一个统一的别名，方便处理。具体的library可以处理该模块公用的方法或load公用的资源，相当于该模块的一个公用类。当然处理方式也可以是直接通过路由中的目录名或者控制器名称来控制等等。

这样避免了加载不同的文件，调用方法时只需要通过$this->helper对象调用。在仔细看看，可以发现不同模块的公用类是放在library中，放在library或helper中都可以使用get_intance获取控制器对象，但每次使用都需要获取实例，相对麻烦，如果是模型呢？感觉也不太好。其中的公用方法有一些会跟业务逻辑相关，放在library感觉不太合适。

业务逻辑好像并没有一个好的地方去实现，控制器的私有方法？模型？

先总结下上面的处理方法：

- 1、 不同模块之间可以按需加载以及实现自定义的公用方法，各个模块之间互不影响。如果各模块之间的公用方法比较多，也可以再去继承一个公用的类。
- 2、 公用方法放在library中，调用CI实例不方便。
- 3、 如果不喜欢$this->herlper的调用方法，可以让控制器去继承不同的公用控制器，思路是一样的，只是可能需要手动引入文件。

## 1.2 业务逻辑
前面对公用控制器按模块分发，方便对特定模块的控制，而具体的实现类则是放在library中。那放在library中是否合适呢？以及控制器中更多的业务逻辑该放在哪里？

先说下对CI中几个文件夹的理解

**helpers、libraries**: 存放一系列辅助函数、辅助类，用来辅助控制器、业务逻辑实现功能。他们中的方法应当尽量避免与CI依赖，依赖越紧越难以复用。以邮件发送为例，发送邮件时很多参数是不变的，如编码、协议、端口等，我们可能会在config下进行配置这些参数，然后library封装一个邮件发送的类，并在其中获取CI实例后读取这些参数。此时就出现了与CI实例的依赖，该类就只能在CI框架中使用，其他系统要用到，就只能重写了，没达到复用的目的。如果发送的类只是接收参数，并封装发送方法呢？所以说，尽可能的让helpers、libraries变的简单，职责变得单一。

**controllers**： 控制器目录。控制器主要用来接管程序，起到连接的作用。通常情况下，我们会把业务逻辑写在action中。但随着业务变得复杂，action代码将越来越臃肿，难以维护。

**models**： 模型目录。CI的模型的主要职责就是和数据库打交道，获取数据。很多时候也会把业务逻辑放在模型中，但业务逻辑与模型实际上是两种东西了。模型只是获取数据，业务逻辑可能是把这些数据根据业务需要进行组合，组合方式可能有很多种，放在模型中会让模型难以维护且不利于复用。说个碰到的例子，对数据按一定条件做缓存，获取数据和缓存结果两个流程写在同一个方法中，但同样的数据需要做另一种形式的缓存时发现，获取数据的方法就没法重用了。

**third_party**：第三方类库目录。拿到一个类库后不要直接使用， 可以在library中进行一次封装，让其更适应于系统，其他人使用起来难度也会降低。

可以发现，每个文件夹都有自己的职责，每个模块都有自己的家，都有自己的职能。那业务逻辑该怎么办？

既然这样， 我们也应该给业务逻辑安个家，建立一个唯一的目录用来存放业务逻辑，暂且命名为service。控制器主要负责接收参数并调用service，service来调用模型，各层各尽其责。下面看看怎么实现：

我们可以重写MY_Load，增加service方法，直接通过$this->load->service('user_service');来调用。但业务逻辑很多都需要获取CI实例，这里可以参考模型的方法，core建立一个MY_Service,其他service均继承该类，这样子service里用法就跟控制器里一样了。
```
class MY_Service
{
    public function __construct()
    {
        log_message('debug', "Service Class Initialized");
    }
  
    function __get($key)
    {
        $CI = & get_instance();
        return $CI->$key;
    }
}
```
其实主要思路还是需要有一层用来处理业务逻辑，java中都有这一层。随着对CI的不断熟悉，发觉这里需要这一层，达到解放控制器和模型的目的。和这种类似的做法还有很多，如果系统中有很多地方需要用到web service 或者说cache之类的，其实也可以按照上面的思路单独放在一个文件夹中处理，方便管理。

## 1.3 控制器
接着前面看， 如果做用户登录功能，用户模块会写在user_service类中。需要登录判断则会增加一个login方法，传入用户名和密码并返回bool值。user_service示例代码如下：
```
public function login($username, $password)
{
    $admin = $this->user_model->getUserInfo(array(
           'username' => $username
    ));
    if(! empty($admin)
       && ($admin['password'] == pwd($password, $admin['salt']))) {
        return true;
    }
    return false;
}
```
用户登录时调用该方法， 根据返回值确定是否要设置登录状态并获取用户资料。如果说要增加ajax控制， 则ajax控制器中只需要接受参数并调用service的方法即可。

来分析下上面这个例子，控制器调用了service， service调用model获取数据并判断密码是否相等。如果没有service层会怎么样？控制器调用model中的方法，并在控制器中判断密码是否相等，或者说在model中实现上面方法。

如果model只是返回数据，则ajax和登录页面都需要自己判断密码是否相同。如果model实现判断过程，大部分情况下我们会不由自主的针对问题来写代码，即取数据和判断放在一个过程中。那如果其他地方需要根据用户名获取用户信息就要重写方法了。 如果取数据独立出去呢？也可以的，大部分时候用CI能控制到这里已经很不错了。但这只是一个简单的逻辑，如果更复杂的逻辑，需要调用多个model，那service的功能就更明显了。所以不建议把业务逻辑写在底层的model中。对于控制器，倒不强制把业务逻辑一定写在service中，但至少可以将一些公用部分，以及复杂的业务逻辑抽离。

接下说下控制器中需要注意的事项：

- 1、 控制器只支持一级目录，不可递归下去，因为支持PATH_INFO的路由方式，想一下都不能递归下去。可使用$this->router对象来获取当前文件夹控制器和方法名。
- 2、 上层文件将会先解析。如果控制下有admin.php也有admin文件夹则会解析到admin.php控制器。
- 3、 控制器中的视图不建议在PHP页面load多个，视图的事情应该交给视图去处理， 写在控制器中视图调整时，控制器也要相应调整。
- 4、 对于公用的Action可以按模块或者按请求类型来区分，具体情况根据项目而定。

# 二、模型
MVC中的业务逻辑放在控制器中或者模型里都是不合适的，所以这里对业务逻辑进行了分离，多出一层用来处理业务逻辑，模型就只当作数据访问层，这样子模型将会变得比较轻。CI中并未通过实体对象来传参，参数的传入和返回都由开发者控制，比较灵活。很多情况下都会以数组的方式传入或者返回。

模型的使用也比较简单，这里只提一下使用前想到的几个问题吧。

1、既然已经有了数据访问层了，那我们就应当避免在控制器或者某些类中直接通过SQL查询数据库，所有的数据库访问操作都应当经过模型获取，大多数情况下一个表名对应着一个模型类。

2、模型应当很方便的连接多个数据库，在database配置文件中有谈到多个库的配置问题，根据group_name的不同可以很方便的连接不同的库。如果有主从，还可以考虑到主从的切换使用问题。

3、模型是否还要按模块区分？在控制器中存在公用控制器分发的做法，在模型中这种思维可能不太好，但可以通过继承不同的公用模型类来实现，这些类再继承CI的MY_Model。在某些业务下根据模块来继承可能比较有用，大部分情况可以直接继承MY_Model，MY_Model主要实现数据库的初始化连接以及一些公用方法。

4、数据库提供的操作方式都是比较基础的，需要我们根据自身的需求去组装，但在我们日常操作中很多操作是类似的，如，根据主键获取信息，根据ID获取信息，根据属性获取信息等，可以对这些基础的操作在进行一次封装，更方便使用。因为如果要使用AR的方式来操作数据库，需要记住很多的方法，如我们根据用户名查询：
```
$query = $this->db->from('user')->where(array('username' => 'BobbyPeng'))->get();
return $query->row_array();
```
如果封装后，则只需要记住一个方法即可，如：
```
public function findByAttributes($where = array())
{
    $query = $this->db->from($this->tableName())->where($where)->get();
    return $query->row_array();
}
```
这样子在每个模型中添加一个tableName的方法返回表名后，再通过模型就可以很方便的使用该方法了。

5、上面的方法属于一个公用方法，我们会写在MY_Model中，但这种类似的方法会很多，我们可否将该类型的方法独立于一个文件中？因为这种方法大部分情况下是不会改的，而放在MY_Model中则表示对它的修改开放了，可能会影响到这些方法。如果说该类叫ActiveRecord类，那可以让MY_Model继承ActiveRecord类，而ActiveRecord类继承CI_Model，参考代码见后面。

很多时候类库提供给我们的方法都是比较细的，我们可以封装一下，减少使用难度。关于模型中公用方法的封装一直还在考虑中，下面给出的只是一个针对单表的简单的操作，复杂的操作就得在特定的模型中实现，还有一些公用操作或者说非AR的操作方式可以统一下，看以后是否有机会再来考虑下这个问题。

公用AR封装类，可进行常用的操作，需要赋予db属性为数据库连接对象，并在模型中设置几个方法，如主键、表名
```
<?php if ( ! defined('BASEPATH')) exit('No direct script access allowed');
  
class ActiveRecord extends CI_Model
{
    /**
     * 保存数据
     *
     * @param array $data 需要插入的表数据
     * @return boolean 插入成功返回ID，插入失败返回false
     */
    public function save($data)
    {
        if($this->db->set($data)->insert($this->tableName())) {
            return $this->db->insert_id();
        }
        return FALSE;
    }
  
    /**
     * Replace数据
     * @param array $data
     */
    public function replace($data)
    {
        return $this->db->replace($this->tableName(), $data);
    }
  
    /**
     * 根据主键更新记录
     *
     * @param string $pk 主键值
     * @param array $attributes 更新字段
     * @param array $where 附加where条件
     * @return boolean true更新成功 false更新失败
     */
    public function updateByPk($pk, $attributes, $where = array())
    {
        $where[$this->primaryKey()] = $pk;
        return $this->updateAll($attributes, $where);
    }
  
    /**
     * 更新表记录
     *
     * @param array $attributes
     * @param array $where
     * @return bollean true更新成功 false更新失败
     */
    public function updateAll($attributes, $where = array())
    {
        return $this->db->where($where)->update($this->tableName(), $attributes);
    }
  
    /**
     * 根据主键删除数据
     *
     * @param string $pk 主键值
     * @param array $where 附加删除条件
     * @return boolean true删除成功 false删除失败
     */
    public function deleteByPk($pk, $where = array())
    {
        $where[$this->primaryKey()] = $pk;
        return $this->deleteAll($where);
    }
      
    /**
     * 删除记录
     *
     * @param array $where 删除条件
     * @param int $limit 删除行数
     * @return boolean true删除成功 false删除失败
     */
    public function deleteAll($where = array(), $limit = NULL)
    {
        return $this->db->delete($this->tableName(), $where, $limit);
    }
      
    /**
     * 根据主键检索
     *
     * @param string $pk
     * @param array $where 附加查询条件
     * @return array 返回一维数组，未找到记录则返回空数组
     */
    public function findByPk($pk, $where = array())
    {
        $where[$this->primaryKey()] = $pk;
        $query = $this->db->from($this->tableName())->where($where)->get();
        return $query->row_array();
    }
  
    /**
     * 根据属性获取一行记录
     * @param array $where
     * @return array 返回一维数组，未找到记录则返回空数组
     */
    public function findByAttributes($where = array())
    {
        $query = $this->db->from($this->tableName())->where($where)->limit(1)->get();
        return $query->row_array();
    }
  
    /**
     * 查询记录
     *
     * @param array $where 查询条件，可使用模糊查询，如array('name LIKE' => "pp%") array('stat >' => '1')
     * @param int $limit 返回记录条数
     * @param int $offset 偏移量
     * @param string|array $sort 排序, 当为数组的时候 如：array('id DESC', 'report_date ASC')可以通过第二个参数来控制是否escape
     * @return array 未找到记录返回空数组
     */
    public function findAll($where = array(), $limit = 0, $offset = 0, $sort = NULL)
    {
        $this->db->from($this->tableName())->where($where);
        if($sort !== NULL) {
            if(is_array($sort)){
                foreach($sort as $value){
                    $this->db->order_by($value, '', false);
                }
            } else {
                $this->db->order_by($sort);
            }
        }
        if($limit > 0) {
            $this->db->limit($limit, $offset);
        }
  
        $query = $this->db->get();
  
        return $query->result_array();
    }
  
    /**
     * 统计满足条件的总数
     *
     * @param array $where 统计条件
     * @return int 返回记录条数
     */
    public function count($where = array())
    {
        return $this->db->from($this->tableName())->where($where)->count_all_results();
    }
  
    /**
     * 根据SQL查询， 参数通过$param绑定
     * @param string $sql 查询语句，如SELECT * FROM some_table WHERE id = ? AND status = ? AND author = ?
     * @param array $param array(3, 'live', 'Rick')
     * @return array 未找到记录返回空数组，找到记录返回二维数组
     */
    public function query($sql, $param = array())
    {
        $query = $this->db->query($sql, $param);
        return $query->result_array();
    }
}
  
/* End of file ActiveRecord.php */
/* Location: ./application/core/database/ActiveRecord.php */
```
MY_Model可以继承该类，这样子模型中可以直接调用上面的方法。
```
<?php if ( ! defined('BASEPATH')) exit('No direct script access allowed');
  
require_once APPPATH.'core/database/ActiveRecord.php';
  
class MY_Model extends ActiveRecord
{
    public function __construct($group_name = '')
    {
        $this->initDb($group_name);
        parent::__construct();
    }
  
    protected function initDb($group_name = '')
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
}
  
/* End of file MY_Model.php */
/* Location: ./application/core/MY_Model.php */
```

# 三、视图
CI中视图即application/views/下的模版文件，模版中支持直接使用PHP，所以模版的实现没有太多好说的，说几点从视图想到的吧。

- 1、CI的视图可以在控制器中load多个视图，页面最后将这些内容组合起来后输出。这里load的动作如果交给控制器去做，需要调整模版结构时就需要调整控制器，这不太好。可以在控制器中进行封装或者直接交- 给视图去做，保证每个ACTION都只load一个视图文件。
- 2、让视图做它擅长的事情，不要在PHP代码中直接定义HTML，这样子会让程序和视图都难以维护，如果要加载HTML，可通过load的方式返回。
- 3、虽然有万能的get_instance方法，但在视图中也不要直接去读取数据或者做一些其他模块的事情。视图的作用就是接收数据并展现出来，以确保人尽其责，物尽其用。
- 4、关于是否要使用模版的问题，个人趋向于使用模版。直接使用PHP时，需要自己去处理变量的定义问题、转换为HTML实体的问题，而这些问题模版都可以解决，以及会提供一些更方便的操作。