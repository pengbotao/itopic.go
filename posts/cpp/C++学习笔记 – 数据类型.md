```
{
    "url": "cpp-primer-plus-type",
    "time": "2015/03/26 21:05",
    "tag": "C++",
    "toc": "no"
}
```

# 基本类型

`C++`的基本类型分为两种：一组由存储为整数的值组成，另一组由存储为浮点格式的值组成。

## 整型

整形可分为：`char`、`short`、`int`、`long`、`long long` ，每一种都有有符号和无符号两种。`char`是一种整形，可以表示计算机中的所有的基本符号，用一个字节表示。需注意越界。其中`short`为`short int`的简称，`long`为`long int`的简称。

**C++的整形标准**

- short至少16位；
- int至少和short一样长；
- long至少32位，且至少和int一样长；
- long long至少64位，且至少与long一样长。

- 8位的char取值范围为 -128 ~ 127，无符号为 0 ~ 255
- 16位的int取值范围为 -32768 ~ +32767，无符号为 0 ~ 65535
- 32位的int取值范围为 -2147483648 ~ +2147483647，无符号为 0~4294967295
- 64位的int取值范围为 -9223372036854775808 ~ +9223372036854775807，无符号为 0 ~ 18446744073709551615

## bool类型
字面值`true`和`false`都可以通过提升类型转换为`int`类型，`true`被转换为1，`false`被转换为0。

## char类型
专为存储字符（如字母和数字）而设计的。`C++`对字符用单引号，对字符串使用双引号。

八进制 042 相当于十进制34 cout << dec; 十进制 42 cout << oct; 十六进制 0x42 相当于十进制66，输出为16进制。 
```
cout << hex;
cout << 123456 << endl;//1e240 cout << 1492 << endl; 
```
除非有理由存储为其他类型（如使用特殊的后缀来表示特定的类型，或者值太大，不能存储为`int`），否则c++将整形常量存储为int类型。 后缀l或L表示该整数为`long`常量，u或U表示`unsigned int`常量，ul（可采用任何一种顺序，大小写均可）表示`unsigned long`

C++11提供了表示`long long`的后缀ll或LL，还提供了`unsigned long long`的后缀ull、Ull、uLL、ULL

## const限定符

常量被初始化之后，其值就被固定了，编译器不允许修改常量的值。

## 浮点数

能够表示带小数部分的数字。浮点数分两部分存储：一部分表示值，另一部分用于对值进行放大或缩小。

3.45E6指3.45与1000000相乘的结果，E6表示10的6次方。因此3.45E6表示34560000，6称为指数，3.45成为尾数。指数为负数表示除以10的城防，如8.55E-4表示0.000855。-8.55E4表示-85500，前面的符号用于数值，指数的符号用户缩放。

`float` 、`double` 、`long double`是按照他们可以表示的有效数位和允许的指数最小范围来描述的。有效位是数字中有意义的位。C和C++对有效位的要求是`float`至少32位`double`至少48位，且不少于`float`。`long double`至少和`double`

通常`float`为32位，`double`为64位，`long double`为80 96或128位，另外这三种类型的指数范围为-37 ~ 37

后缀默认情况为都属于`double`类型，如果希望为`float`可使用f或F后缀。对于`long double`类型，可使用l或L后缀。

# 复合类型
复合类型是基于整形和浮点型创建的。影响最为深远的复合类型是类。

## 数组
数组是一种数据格式，能够存储多个同类型的值。
```
typeName arrayName[arraySize]
```
`int months[12] = {1, 2, 3};`其他元素将设置为0

`char name[4] = {'p', 'i', 'g', '\0'};`

不能将一个数组赋值给另一个数组。

## 字符串
存储在内存的连续字节中的一系列字符。C++处理字符串的方式有两种。一种来自C语言，常被称为C-风格字符串，另一种基于string类库的方法。

**C-风格字符串：**
以空字符 \0 结尾，其ASCII为0，用来标记字符串的结尾。
```
char name[4] = {'p', 'i', 'g', '!'};//不是字符串
char name[4] = {'p', 'i', 'g', '\0'};//是字符串
```

用引号括起的字符串称为字符串常量或字符串字面值，隐式包含结尾的空字符，如：
```
char bird[11] = "Mr. Cheeps";
char fish[] = "Bubbles";
```

strlen不计算空字符，`strlen(bird) == 10`
```
strcpy(charr1, charr2);//copy charr2 to charr1
strcat(charr1, charr2);//append contents of charr2 to charr1
 
strncpy(food, "a picnic basket filled with many goodies", 19);
food[10] = '\0';
```
**cin如何确定已完成字符串的输入?**
由于不能通过键盘输入空字符，因此cin需要用别的发那个发来确定字符串的结尾位置。cin使用空白、制表符和换行符来确定字符串的结束位置。

读取一行：`cin.getline()`和`cin.get()`。
这两个函数否读取一行输入，直到达到换行符。不同的是，`cin.getline()`将丢弃换行符，`get`将保留换行符在输入序列中。`cin.getline(fish, 20);`  `cin.get()`读取下一个字符

**string类**

- 可以使用C-风格字符串来初始化string对象，如`string str2 = "test";`
- 可以使用cin、cout来输入或输出string对象
- 可以使用数组表示法来访问存储在string对象中的字符

可以将一个string对象赋值给另一个string对象，string类与C-风格字符串
```
#include <iostream>
#include <string>
 
using namespace std;
 
void main()
{
    char char1[20];
    char char2[20] = "c string";
 
    string str1;
    string str2 = "string";
 
    str1 = str2;
    strcpy(char1, char2);
 
    str1 += " !";
    strcat(char1, " !");
     
    int len1 = strlen(char1);
    int len2 = str1.size();
    cout << char1 << endl << str1 << endl;
}
```
string类读取一行 `getline(cin, str)`;

## 结构和结构数组
```
struct sname
{
    int id;
    char name[20];
 
};
void main()
{
    sname s[1];
    s[0].id = 1;
    strcpy(s[0].name, "pig");
 
    cout << s[0].id << ":" << s[0].name << endl;
}
```
## 共用体

共用体能存储不同的数据类型，但只能同时存储其中的一种类型。

## 枚举
```
enum 枚举名{
    标识符[=整型常数],
    标识符[=整型常数],

    ...
    标识符[=整型常数],
} 枚举变量;
enum spectrum{red = 1, orange, yellow, green, blue};
```
只能将整型赋值给枚举。

## 指针 - 变量的地址

OOP强调的是在运行阶段（而不是编译阶段）进行决策。

**指针：**用于存储值的地址。*运算符被成为间接值或解除引用运算符，将其应用于指针，可以得到该地址处存储的值。
```
int * id;
char *pc = new char;
```
通常情况下，地址需要2个还是4个字节取决于计算机系统。一定要在对指针应用解除引用运算符之前，将指针初始化一个确定的、适当的地址。

C语言可通过`malloc`来分配内存；在C++中仍然可以这样做，但C++还有更好的方法 - `new`运算符。

为一个数据对象（可以是结构，也可以是基本类型）获得并指定分配内存的通用格式如下：
```
typeName * pointer_name = new typeName;
```
new分配的内存块通常与常规变量声明分配的内存块不同。常规变量存储在栈中，而new存储在堆或自由存储区的内存区域中。
使用`new`分配的内存需要使用`delete`来释放，一定要配对使用`new`和`delete`，否则将发生内存泄漏，也就是说，被分配的内存再也无法使用了。如果使用`new`时带方括号，则`delete`也应带上，如果没带，则不用带。

**new和delte遵守规则**

- 不要使用delete来释放不是new分配的内存
- 不要使用delete释放同一个内存卡两次
- 如果使用new [] 为数组分配内存，则应使用delete []来释放
- 如果使用new为一个实体分配内存，则应使用delete来释放
- 对空指针应用delete是安全的。

在编译时给数组分配内存称为静态联编，在运行时动态分配称为动态联编。

**使用new创建动态数组**
为数组分配内存的通过格式
```
typeName * pointer_name = new typeName[num_elements];
```
```
int * psome = new int[10];
delete [] psome;
```

**使用动态数组**
指针指向内存块中的第一个元素，所以`*psome`是第一个元素的值。也可使用数组的方式来访问，如`psome[2]`。数组与指针基本等价是C和C++的优点之一

**C++将数组名解释为地址**
指针变量加1后，增加的量等于它指向的类型的字节数。

对数组使用sizeof运算符得到的是数组的长度，而指针使用sizeof得到的是指针的长度，即使指针指向的是一个数组。

数组名被解释为第一个元素的地址，对数组名应用地址符时，得到的是整个数组的地址。
```
short tell[10];
cout << tell << endl;
cout << &tell << endl;
```
从数字上说，这两个地址相同。但从概率上说，&tell[0]是一个2字节内存块的地址，而&tell是以一个20字节的内存块的地址。因此，表达式tell+1将地址值增加2，而&tell+2将地址加20.

如果结构标识符是结构名，则使用句点运算符；如果标识符是指向结构的指针，则使用箭头运算符(->)

自动存储、静态存储和动态存储

## 数组的替代品

**模板类`vector`**
```
vector<typeName> vt(n_elem);
```
n_elem可以是整形常量，也可以是整形变量

**模版类`array`**
```
array<typeName, n_elem> arr;
```
与创建`vector`不同，`n_elem`不能是变量

数组、`vector`、`array`都可使用标准数组表示法来访问各个元素。