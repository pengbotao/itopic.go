```
{
    "url": "php-clone",
    "time": "2014/03/03 21:23",
    "tag": "PHP"
}
```

周末闲来无事看到了原型模式，其中谈到了浅复制和深复制，想到PHP中的对应赋值、克隆以及克隆是浅复制还是深复制。

先来看看赋值，例如有一个简历类，有身高和体重两个属性：
```
class Resume 
{
    public $height;
    public $weight;
 
    public $workExperience;
}
$ResumeA = new Resume();
$ResumeB = $ResumeA;
```
此时实例化了一个Resume类并赋值给了$ResumeA变量，然后将$ResumeA变量赋值给$ResumeB。PHP手册上有说：

自PHP5起，new运算符自动返回一个引用，一个对象变量已经不再保存整个对象的值。只是保存一个标识符来访问真正的对象内容。 当对象作为参数传递，作为结果返回，或者赋值给另外一个变量，另外一个变量跟原来的不是引用的关系，只是他们都保存着同一个标识符的拷贝，这个标识符指向同一个对象的真正内容。

所以若通过$ResumeB修改height属性，则$ResumeA也会跟着变。如果想要复制一个全新的对象，则可以通过clone来实现，如：
```
$ResumeB = clone $ResumeA;
```
此时将$ResumeA的值拷贝到新的变量$ResumeB中，改变其中一个不影响另一个，修改$ResumeB中height属性，$ResumeA不会跟着改变。
但如果该类引用了其他对象，则所有的引用仍然指向到原来的对象。clone的这种复制方式就是浅复制。**被赋值对象的所有变量都还有与原来对象相同的值，而所有的对其他对象的引用都仍然指向原来的对象。**

如果上面类中workExperience为WorkExperience类的引用，当克隆的时候，克隆前后的workExperience属性还是指向到同一个对象内容。

与浅复制对应的是深复制，**深复制把引用对象的变量指向复制过的新对象，而不是原有的被引用的对象。**

PHP中可以通过两种方式来实现深复制。第一种是__clone魔术方法：
```
public function __clone()
{
    $this->workExperience = new WorkExperience();
}
```
深复制涉及深的层次，通过clone魔术方法实现需要知道有几层然后对每一层依次实现。还有一种是可以通过序列化对象的方式，先将对象序列化之后再反序列化，如：
```
$ResumeB = unserialize(serialize($ResumeA));
```
clone还算常用的拷贝方式，整理的目的只是为了记录一下clone是浅复制，需要注意一下对象的引用。