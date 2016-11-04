```
{
    "url": "lua-class",
    "time": "2016/01/17 20:22",
    "tag": "Lua"
}
```

Lua中没有类的概念，但可以通过table、函数以及一些语言特性来模拟类。table的值可以为函数。看下面这段代码：
```
LogUtil = {msg_prefix = string.format("%s ", os.date("%H:%M:%S", os.time())) }
LogUtil.write = function(msg)
    print(LogUtil.msg_prefix .. msg .. "\r\n")
end
 
LogUtil.write("This is a Notice Info")
```
定义一个日志table，table有一个名为msg_prefix的键以及一个名为write的键。msg_prefix指向一个字符串，write指向一个方法，然后调用了write方法。这里看起来就是一个table的操作，但也可以换个角度来看：定义了一个日志类，类里有一个静态成员属性msg_prefix以及一个静态成员方法write。并且调用了write方法。

**这也就是Lua实现类的基础。**

上面代码中，LogUtil为全局变量。如果重写了它的值将会导致方法不可用。所以可以用一个额外的参数来传入对象，这个参数经常为self或this，好比Python类方法中也需要显示的传入self对象。如：
```
LogUtil = {msg_prefix = string.format("%s ", os.date("%H:%M:%S", os.time())) }
 
LogUtil.write = function(self, msg)
    print(self.msg_prefix .. msg .. "\r\n")
end
 
LogA = LogUtil
LogUtil = nil
LogA.write(LogA, "This is a LogA Info")
```
但这样子调用的时候也需要传入self参数，不过，Lua也提供了通过使用冒号操作符来隐藏这个参数的声明。看下面，
```
LogUtil = {msg_prefix = string.format("%s ", os.date("%H:%M:%S", os.time())) }
 
function LogUtil:write(msg)
    print(self.msg_prefix .. msg .. "\r\n")
end
 
LogA = LogUtil
LogA:write("This is a LogA Info")
```
类定义还是被调用，通过冒号方式的调用都会自动将类作为self参数传入，而省去了主动传入的麻烦。还有， function提在前面了，可以看到这里更有函数的味道了。上面两段代码是等价的，归纳下就是：
```
function class.a.b:c = function(param) body end
```
等价于
```
class.a.b.c = function(self, param) body end
```
可能，会被上面的点、冒号、self搞混，其实可以统一使用冒号的方式来定义与调用类方法。

到目前为止还只是一些属性和方法的集合，并没有实例化类，属性的改变会导致全局的变化。接下来看如何实现多个互不影响的类的实例。
# 类的实例
在Lua中有一个神奇的东西叫做metatable(元表)。metatable也是键值对。当访问table中一个不存在的属性时就会触发一些事件。比如当设置了元方法__index时 表示当调用当前table中不存在的属性或方法时会到元方法对应的table中去找，从而实现继承的效果。
```
setmetatable(a, {__index = b})
```
这样，对象a调用任何不存在的成员都会到对象b中查找。术语上，可以将b看作类，a看作对象。
```
LogUtil = {msg_prefix = string.format("%s ", os.date("%H:%M:%S", os.time())) }
 
function LogUtil:write(msg)
    print(self.msg_prefix .. msg .. "\r\n")
end
 
-- 创建实例化的new方法，名字可以任意制定。
function LogUtil:new(o)
    o = o or {}
    setmetatable(o, self)
    self.__index = self
    return o
end
 
LogA = LogUtil:new()
LogA.msg_prefix = ""
LogA:write("This is a LogA Info")
 
LogB = LogUtil:new{msg_prefix = "Hello."}
LogB:write("This is a LogB Info")
 
--[[
-- output:
This is a LogA Info
 
Hello.This is a LogB Info
]]
```
可以看到，A和B的操作是互不影响的，相当于对类进行了实例化。有两点需要注意的地方：

- 1、new方法中的setmetatable(o, self)只是做了一个小小的优化，不需要创建一个额外的表作为LogA或LogB的元表。再简单点说， 元表只是取了LogUtil的一个\_\_index键。对于LogUtil来说，只是多了一个\_\_index键而已。
- 2、LogB实例化的时候参数为一个table，此时可以省略括号。表示在o中已经设置了msg_prefix，不再需要去LogUtil寻找，达到重写的目的。

当执行LogA:write()时先在LogA中查找，LogA没有此方法，此时会到LogUtil中去找。基本上已经也有了继承的概念。

# 类的继承

在前面基础上我们调整下调用，LogA为LogUtil的实例，找不到的方法都到LogUtil中去找。如果想要LogB继承自LogA呢？可以直接在new一下LogA。
```
LogA = LogUtil:new{msg_prefix = 'A:'}
LogA:write("This is a LogA Info")
 
LogB = LogA:new{}
LogB:write("This is a LogB Info")
```
LogA再执行new方法时，此时的self参数为LogA了，所以相当于LogB继承自LogA。

同样， Lua还可以实现多重继承，成员属性私有等，这些就留待以后去深究了。