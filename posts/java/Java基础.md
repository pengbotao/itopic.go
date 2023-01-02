```
{
    "url": "java",
    "time": "2019/07/08 19:45",
    "tag": "Java"
}
```

# 一、8中基本数据类型

| 基本数据类型 | 包装类型            |
| ------------ | ------------------- |
| byte         | java.lang.Byte      |
| short        | java.lang.Short     |
| int          | java.lang.Integer   |
| long         | java.lang.Long      |
| float        | java.lang.Float     |
| double       | java.lang.Double    |
| boolean      | java.lang.Boolean   |
| char         | java.lang.Character |

### **为何要使用包装类型？**

- 基本类型不是类，不能new出来，因此不具备面对对象的功能，无法调用方法。
- 在一个类或接口或方法中定义一个泛型的数据类型，当使用这个类、接口、方法时，要把泛型定义成具体的基本数据类型就必须使用基本数据类型对应的包装类进行定义。  



```
Integer integer1 = new Integer(100);    // 以 int 型变量作为参数创建 Integer 对象
Integer integer2 = new Integer("100");    // 以 String 型变量作为参数创建 Integer 对象

int num = Integer.parseInt(str);    // 将字符串转换为int类型的数值
String s = Integer.toString(i);    // 将int类型的数值转换为字符串


Float float1 = new Float(3.14145);    // 以 double 类型的变量作为参数创建 Float 对象
Float float2 = new Float(6.5);    // 以 float 类型的变量作为参数创建 Float 对象
Float float3 = new Float("3.1415");    // 以 String 类型的变量作为参数创建 Float 对象
float num = Float.parseFloat(str);    // 将字符串转换为 float 类型的数值
String s = Float.toString(f);    // 将 float 类型的数值转换为字符串
```

# 二、数组

## 2.1 数组定义

```java
//方式1
int[] arr = new int[3];
arr[0] = 1;

//方式2
int[] arr = new int[]{1,2,3,4,5};

//方式3
int[] arr = {1,2,3,4,5};

System.out.println(Arrays.toString(arr));
```

## 2.2 数据遍历

```java
int[] x = {1,3,2,4};

for (int i : x) {
    System.out.println(i);
}


String[][][] namelist = {
        {{ "张阳", "李风", "陈飞" }, { "乐乐", "飞飞", "小曼" }},
        {{ "Jack", "Kimi" }, { "Lucy", "Lily", "Rose" }},
        {{ "徐璐璐", "陈海" }, { "李丽丽", "陈海清" }}
};
for (int i = 0; i < namelist.length; i++) {
    for (int j = 0; j < namelist[i].length; j++) {
        for (int k = 0; k < namelist[i][j].length; k++) {
            System.out.println("namelist[" + i + "][" + j + "][" + k + "]=" + namelist[i][j][k]);
        }
    }
}
```

## 2.3 Arrays工具类

```golang
Integer[] x = {1,3,2,4};

//降序，需要是包装类型
Arrays.sort(x, Collections.reverseOrder());
//升序
Arrays.sort(x);
//二分查找
int pos = Arrays.binarySearch(x, 2);
//打印
System.out.println(Arrays.toString(x));
```

# 三、集合

## 3.1 List集合

List 实现了 Collection 接口，它主要有两个常用的实现类：ArrayList 类和 LinkedList 类。

```java
ArrayList<String> list = new ArrayList();
list.add("A");
list.add("B");
list.add("C");
list.remove("C");
list.set(1, "BB");

for (String v : list) {
    System.out.println(v);
}

for (int i = 0; i<list.size(); i++) {
    System.out.println(list.get(i));
}

Iterator it = list.iterator();
while (it.hasNext()) {
    System.out.println(it.next());
}
```

## 3.2 Set集合

Set 实现了 Collection 接口，它主要有两个常用的实现类：HashSet 类和 TreeSet类。

```java
HashSet<String> set = new HashSet();
set.add("A");
set.add("B");
set.add("A");

Iterator<String> it = set.iterator();
while (it.hasNext()) {
    System.out.println(it.next());
}
```

## 3.3 Map集合

Map 接口主要有两个实现类：HashMap 类和 TreeMap 类。其中，HashMap 类按哈希算法来存取键对象，而 TreeMap 类可以对键对象进行排序。

1）在 for 循环中使用 entries 实现 Map 的遍历（最常见和最常用的）。

```java
HashMap<String, Integer> map = new HashMap();
map.put("A", 1);
map.put("B", 2);

for (Map.Entry<String, Integer> entry: map.entrySet()) {
    String k = entry.getKey();
    Integer v = entry.getValue();
    System.out.println("k=" + k + ", v =" + v );
}
```

2）使用 for-each 循环遍历 key 或者 values，一般适用于只需要 Map 中的 key 或者 value 时使用。性能上比 entrySet 较好。

```java
Map<String, Integer> map = new HashMap<String, Integer>();
map.put("A", 1);
map.put("B", 2);

// 打印键集合
for (String key : map.keySet()) {
    System.out.println(key);
}
// 打印值集合
for (Integer value : map.values()) {
    System.out.println(value);
}
```

3）使用迭代器（Iterator）遍历

```java
HashMap map = new HashMap();
map.put("A", 1);
map.put("B", 2);

Iterator it = map.keySet().iterator();
while (it.hasNext()) {
    Object k = it.next();
    Object v = map.get(k);
    System.out.println("k=" + k + ", v =" + v );
}
```



## 3.4 Collections类

Collections 类是 Java提供的一个操作 Set、List 和 Map 等集合的工具类。Collections 类提供了许多操作集合的静态方法，借助这些静态方法可以实现集合元素的排序、查找替换和复制等操作。下面介绍 Collections 类中操作集合的常用方法。

### 排序（正向和逆向）

Collections 提供了如下方法用于对 List 集合元素进行排序。

- void reverse(List list)：对指定 List 集合元素进行逆向排序。
- void shuffle(List list)：对 List 集合元素进行随机排序（shuffle 方法模拟了“洗牌”动作）。
- void sort(List list)：根据元素的自然顺序对指定 List 集合的元素按升序进行排序。
- void sort(List list, Comparator c)：根据指定 Comparator 产生的顺序对 List 集合元素进行排序。
- void swap(List list, int i, int j)：将指定 List 集合中的 i 处元素和 j 处元素进行交换。
- void rotate(List list, int distance)：当 distance 为正数时，将 list 集合的后 distance 个元素“整体”移到前面；当 distance 为负数时，将 list 集合的前 distance 个元素“整体”移到后面。该方法不会改变集合的长度。

### 查找、替换操作

Collections 还提供了如下常用的用于查找、替换集合元素的方法。

- int binarySearch(List list, Object key)：使用二分搜索法搜索指定的 List 集合，以获得指定对象在 List 集合中的索引。如果要使该方法可以正常工作，则必须保证 List 中的元素已经处于有序状态。
- Object max(Collection coll)：根据元素的自然顺序，返回给定集合中的最大元素。
- Object max(Collection coll, Comparator comp)：根据 Comparator 指定的顺序，返回给定集合中的最大元素。
- Object min(Collection coll)：根据元素的自然顺序，返回给定集合中的最小元素。
- Object min(Collection coll, Comparator comp)：根据 Comparator 指定的顺序，返回给定集合中的最小元素。
- void fill(List list, Object obj)：使用指定元素 obj 替换指定 List 集合中的所有元素。
- int frequency(Collection c, Object o)：返回指定集合中指定元素的出现次数。
- int indexOfSubList(List source, List target)：返回子 List 对象在父 List 对象中第一次出现的位置索引；如果父 List 中没有出现这样的子 List，则返回 -1。
- int lastIndexOfSubList(List source, List target)：返回子 List 对象在父 List 对象中最后一次出现的位置索引；如果父 List 中没有岀现这样的子 List，则返回 -1。
- boolean replaceAll(List list, Object oldVal, Object newVal)：使用一个新值 newVal 替换 List 对象的所有旧值 oldVal。

### 复制

Collections 类的 copy() 静态方法用于将指定集合中的所有元素复制到另一个集合中。执行 copy() 方法后，目标集合中每个已复制元素的索引将等同于源集合中该元素的索引

## 3.5 Iterator（迭代器）

Iterator（迭代器）是一个接口，它的作用就是遍历容器的所有元素，也是 Java 集合框架的成员，但它与 Collection 和 Map 系列的集合不一样，Collection 和 Map 系列集合主要用于盛装其他对象，而 Iterator 则主要用于遍历（即迭代访问）Collection 集合中的元素。

```java
// 创建一个集合
Collection objs = new HashSet();
objs.add("C语言中文网Java教程");
objs.add("C语言中文网C语言教程");
objs.add("C语言中文网C++教程");
// 调用forEach()方法遍历集合
// 获取books集合对应的迭代器
Iterator it = objs.iterator();
while (it.hasNext()) {
    // it.next()方法返回的数据类型是Object类型，因此需要强制类型转换
    String obj = (String) it.next();
    System.out.println(obj);
    if (obj.equals("C语言中文网C语言教程")) {
        // 从集合中删除上一次next()方法返回的元素
        it.remove();
    }
    // 对book变量赋值，不会改变集合元素本身
    obj = "C语言中文网Python语言教程";
}
System.out.println(objs);
```

# 四、类

```java
[public][abstract|final]class<class_name>[extends<class_name>][implements<interface_name>] {
    // 定义属性部分
    <property_type><property>;
    [public|protected|private][static][final]<type><variable_name>;
    …
    // 定义方法部分
    function();
    [public|private|protected][static]<void|return_type><method_name>([paramList]) {
        // 方法体
    }
    …
}
```

# 五、异常

```java
try {
    逻辑程序块
} catch(ExceptionType1 e) {
    处理代码块1
} catch (ExceptionType2 e) {
    处理代码块2
    throw(e);    // 再抛出这个"异常"
} finally {
    释放资源代码块
}
```

# 六、注解

`@Override ` 注解是用来指定方法重写的，只能修饰方法并且只能用于方法重写，不能修饰其它的元素。它可以强制一个子类必须重写父类方法或者实现接口的方法。

`@Deprecated` 可以用来注解类、接口、成员方法和成员变量等，用于表示某个元素（类、方法等）已过时。当其他程序使用已过时的元素时，编译器将会给出警告。
