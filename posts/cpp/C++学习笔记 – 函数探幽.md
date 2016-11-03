url: cpp-primer-plus-function-2
des: 
time: 2015/03/30 22:45
category: cpp
++++++++

# 内联函数

内联函数的编译代码与其他程序代码“内联”起来，也就是说，编译器将使用相应的函数代码替换函数调用。对于内联代码，程序无需跳到另一个位置处执行代码，再跳回来。因此，内联函数的运行速度比常规函数稍快，但代价是需要占用更多内存。

通常的做法是省略函数原型，将整个定义放在本应提供原型的地方。程序员请求将函数作为内联函数时，编译器不一定会满足这种要求。它可能认为该函数过大或注意到函数调用了自己（内联函数不能递归），因此不能将其作为内联函数；而有些编译器没有启用或实现这种特性。
```
inline double square(double x){return x*x;}
```
# 引用变量

引用是已定义变量的别名，主要用途是用作函数的形参。通过将引用变量用作参数，函数将使用原始数据，而不是其副本。这样除指针之外，引用也为函数处理大型结构提供了方便的途径。

**创建引用变量**

C和C++使用&符号来指示变量的地址。C++给&符号赋予了另一种含义，将其用来声明引用。例如，将rodents作为rats变量的别名，可以这样做：
```
int rats;
int & rodents = rats;
```
其中，&不是地址运算符，而是类型标识符的一部分。就像声明中的char*指的是指向char的指针一样，int&指的是指向int的引用。上述引用声明允许将rats和rodents互换 —— 他们指向相同的值和内存单元。

引用接近const指针，必须在创建时进行初始化，一单与某个变量关联起来，就将一直肖忠于它。

按引用传递和按值传递看起来相同，只能通过原型或函数定义才能知道函数是按引用传递的。
```
#include <iostream>
 
void swapr(int &a, int&b);
 
using namespace std;
void main()
{
    int a = 300, b= 500;
    swapr(a, b);
}
 
void swapr(int &a, int&b)
{
    int temp;
    temp = a;
    a = b;
    b=temp;
}
```
上面代码传递a、b的引用，可以对其值进行修改，若限制其修改，可使用常量引用：
```
int func(const int &);
```
**临时变量、引用参数和const**

如果接受引用参数的意图是修改作为参数传递的变量，则创建临时变量将阻止这种意图的实现。解决方式是禁止创建临时变量。如果指定引用为const，不可以修改传递的值，此时生成临时变量不会造成任何影响。

**引用与结构体**

使用结构引用参数的方式与使用基本变量引用相同，只需在声明结构参数时使用引用运算符&即可。如：
```
free_thorws & accumulate(free_throws & target, const free_thorws & source);
```
什么时候使用引用、什么时候使用指针？ 不需要修改参数：

- 如果数据对象很小，如内置数据类型或小型结构，则按值传递。
- 如果数据对象是数组，则使用指针，因为这是唯一选择，并将指针声明为指向const的唯一指针。
- 如果数据对象是较大的结构，则使用const指针或const引用，以提供程序的效率。节省复制结构需要的时间和空间。
- 如果数据对象是类对象，则使用const引用。

需要修改参数：

- 如果数据对象是内置数据类型，则使用指针。
- 如果数据对象是数组，则只能使用指针。
- 如果数据对象是结构，则使用引用或指针。
- 如果数据对象是类对象，则使用引用。

# 默认参数

默认参数之的是当函数调用中省略了实参时自动使用的一个值。设置默认参数只需在函数原型中指定即可，如：
```
char * left(const char * str, int n=1);
```
对于带参数列表的函数，必须从右向左添加默认值。

# 函数重载

函数重载的关键是函数的参数列表 —— 也称为函数特征标。C++允许定义名称相同的函数，条件是他们是特征标不同。
```
#include <iostream>
 
char * left(const char * str, int n);
unsigned long left(unsigned long num, int n);
 
using namespace std;
void main()
{
    cout << left("abc", 2) << endl;
    cout << left(12345, 3) << endl;
}
 
char * left(const char * str, int n)
{
    if(n < 0) n = 0;
    char *rtn = new char[n+1];
    int i = 0;
    while(i < n && str[i]) {
        rtn[i++] = str[i];
    }
    while(i <= n) {
        rtn[i++] = '\0';
    }
    return rtn;
 
}
unsigned long left(unsigned long num, int n)
{
    int len = 1;
    unsigned long t = num;
    if(n <= 0) n=0;
    while(t = t / 10) {
        len ++;
    }
    if(n >= len) {
        return num;
    }
    n = len - n;
    while(n--) {
        num = num / 10;
    }
    return num;
}
```
# 函数模版
```
#include <iostream>
 
template <typename T>
void Swap(T &a, T &b);
 
void main()
{
    using namespace std;
    int i = 10, j = 20;
    Swap(i, j);
    float x = 1.5f, y=2.1f;
    Swap (x, y);
    cout << i << j << x << y;
}
 
template <typename T>
void Swap(T &a, T &b)
{
    T temp;
    temp = a;
    a = b;
    b = temp;
}
```