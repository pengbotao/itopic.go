```
{
    "url": "codeigniter-helpers",
    "time": "2013/11/15 10:15",
    "tag": "PHP,CodeIgniter"
}
```

# 一、helpers、libraries
前面提到过helper、和libraries，主要用来存放一系列辅助函数、辅助类，用来辅助系统实现功能。但helper 和 library 之间到底有什么区别呢？什么时候该用 helper 什么时候该用 library ？这好像是个无聊的问题。。。

来谈下无聊的看法：helper里主要是一些函数， library里主要是一些类的对象。函数表示的是一种行为，类表示的是一种抽象的概念，是一系列属性和方法的集合，对象里的函数称为方法，数据则称为对象的属性。类是有状态的，而函数无状态的，所以函数与类之间最大的区别在于是否有状态，落到更实际的点就是方法与方法之间是否需要共享数据，如果需要共享数据，则可以写成类的方法，如果不需要共享数据，那就可以用写成函数。

很多时候我们会写一些静态类，每个方法都是独立的，这些方法是可以用函数来代替的，而有些时候我们需要先初始化某些参数，然后后面的方法时可以直接用到这些参数，这些不是函数所擅长的。个人相对更喜欢用类来封装，即便有时候时候没有数据的传递，但当需要传输数据时，类会更方便扩展一点。

弄清楚了这些，对 helper 和 library 的理解可能会有一点帮助。接下来谈谈该怎么用 helper 和 library。

上周分享的时候，有些同事说把业务逻辑写在helper中，来看看这样子会有什么不好。首先是很多地方都需要调用get_instance方法来获取CI实例，然后才能调用CI的方法，使用起来不不太方便。接着就是业务逻辑的易变性导致了helper中方法会更有针对性，针对特定的业务逻辑去实现，可能达不到重用的效果。可能有人会说，我会封装的比较好，但我认为这种思维的规律导致了只能是业务的复用，helper函数的复用会降低。同样的library也是类似的道理，所以说不建议 helper 和 library 里调用太多的CI实例，产生的依赖越多，代码就越难复用。其实最开始也考虑过使用library来实现业务逻辑，根据文件夹来区分业务逻辑和类库，也许是因为上面的原因始终觉得不是很好，所以最终才决定增加 service 。

所以用之前需要有这样的意识，**尽量减少与CI的依赖，将问题分析清楚后拆解，适合helper和library的代码才写在这里。**

接下来就是PHP功力的比拼了，不管函数或者类需要遵守一个很重要的原则：单一职责， 即一个类，最好只做一件事，只有一个引起它变化的原因。如果一个函数用一句话描述不清楚，那就需要拆成多个函数；如果一个类有跟它不相干的职责，那就拆成多个类。函数应该短小精悍，如果一个函数太长，意味着关联太多，可拆分的可能性很大。编码的时候应时刻提醒自己，它是不是只做了一件事，还能不能在细分？单一职责是程序员要遵守的首要原则，然后再加上适当的说明以及参数的注释，良好的排版风格，我相信这个函数乱不到哪里去。类更复杂点，怎么样组织代码结构，怎么样更方便调用，这些需要编程知识的积累。
最后，我们可以多回头看看自己的函数或者类，从中吸取经验，让下次写的更好。

# 二、third_party
third_party用来存放系统中引入的第三方类库，类库通常提供的功能比较丰富，相应的学习成本也要高些以及系统中能用到功能有限，所以建议在引入类库时进行适当的封装，让系统中更方便使用，其他人使用时只需关注扩展的方法而无法关注具体的实现。以CI集成Twig模版为例吧。

首先需要下载Twig类库，并放在third_party中，然后在libraries中进行一次封装，示例如下：
```
<?php  if ( ! defined('BASEPATH')) exit('No direct script access allowed');
  
require APPPATH.'third_party/Twig/Autoloader.php';
  
/**
 * Twig模版引擎
 *
 */
class Twig
{
    public $twig;
      
    public $config;
      
    private $data = array();
      
    /**
     * 读取配置文件twig.php并初始化设置
     *
     */
    public function __construct($config)
    {
        $config_default = array(
            'cache_dir' => false,
            'debug' => false,
            'auto_reload' => true,
            'extension' => '.tpl',
        );
        $this->config = array_merge($config_default, $config);
        Twig_Autoloader::register ();
        $loader = new Twig_Loader_Filesystem ($this->config['template_dir']);
        $this->twig = new Twig_Environment ($loader, array (
                'cache' => $this->config['cache_dir'],
                'debug' => $this->config['debug'],
                'auto_reload' => $this->config['auto_reload'],
        ) );
        $CI = & get_instance ();
        $CI->load->helper(array('url'));
        $this->twig->addFunction(new Twig_SimpleFunction('site_url', 'site_url'));
        $this->twig->addFunction(new Twig_SimpleFunction('base_url', 'base_url'));
    }
      
    /**
     * 给变量赋值
     *
     * @param string|array $var
     * @param string $value
     */
    public function assign($var, $value = NULL)
    {
        if(is_array($var)) {
            foreach($val as $key => $val) {
                $this->data[$key] = $val;
            }
        } else {
            $this->data[$var] = $value;
        }
    }
  
    /**
     * 模版渲染
     *
     * @param string $template 模板名
     * @param array $data 变量数组
     * @param string $return true返回 false直接输出页面
     * @return string
     */
    public function render($template, $data = array(), $return = FALSE)
    {
        $template = $this->twig->loadTemplate ( $this->getTemplateName($template) );
        $data = array_merge($this->data, $data);
        if ($return === TRUE) {
            return $template->render ( $data );
        } else {
            return $template->display ( $data );
        }
    }
      
    /**
     * 获取模版名
     *
     * @param string $template
     */
    public function getTemplateName($template)
    {
        $default_ext_len = strlen($this->config['extension']);
        if(substr($template, -$default_ext_len) != $this->config['extension']) {
            $template .= $this->config['extension'];
        }
        return $template;
    }
  
    /**
     * 字符串渲染
     *
     * @param string $string 需要渲染的字符串
     * @param array $data 变量数组
     * @param string $return true返回 false直接输出页面
     * @return string
     */
    public function parse($string, $data = array(), $return = FALSE)
    {
        $string = $this->twig->loadTemplate ( $string );
        $data = array_merge($this->data, $data);
        if ($return === TRUE) {
            return $string->render ( $data );
        } else {
            return $string->display ( $data );
        }
    }
}
  
/* End of file Twig.php */
/* Location: ./application/libraries/Twig.php */
```
模版的操作通常有一些配置的信息，这里通过config下的twig.php进行配置，通过CI load library的方式加载时，与类名同名的配置文件存在时，会自动以数组的方式将参数传入类的构造函数。
```
<?php
// 默认扩展名
$config['extension'] = ".tpl";
  
// 默认模版路劲
$config['template_dir'] = APPPATH . "views/";
  
// 缓存目录
$config['cache_dir'] = APPPATH . "cache/twig/";
  
// 是否开启调试模式
$config['debug'] = false;
  
// 自动刷新
$config['auto_reload'] = true;
  
/* End of file twig.php */
/* Location: ./application/config/twig.php */
```
为了加载base_url site_url等函数到模版，类与CI产生了依赖，分离开可能更好，比如在serice中进行一次封装，增加一些自定义函数等，这样其他地方、其他系统也就很方便复用该类了。