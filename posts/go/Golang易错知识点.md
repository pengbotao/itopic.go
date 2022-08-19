```
{
    "url": "go-tips",
    "time": "2021/02/04 13:00",
    "tag": "Golang",
    "toc": "no"
}
```

### 1. for ... range

```
package main

import "fmt"

func main() {
	m := []int{1, 2, 3}
	n := make(map[int]*int)
	for k, v := range m {
		fmt.Println(k, v)
		n[k] = &v
	}
	fmt.Printf("%+v", n)
}
```

> Output:

```
0, 1, 0xc000014168, 0xc000014180
1, 2, 0xc000014168, 0xc000014180
2, 3, 0xc000014168, 0xc000014180
map[0:0xc000014180 1:0xc000014180 2:0xc000014180]
```

可以看到Range里k,v的变量地址未发生改变，多次循环使用了同一块内存地址接收。

### 2、切片陷阱

```
package main

import "fmt"

func main() {
	a := [...]int{1, 2, 3, 4, 5}
	s1 := a[1:3]
	fmt.Printf("%#v, len:%d, cap: %d\n", s1, len(s1), cap(s1))
	s1[0] = 6
	fmt.Println(a, s1)
}
```

切片的底层结构是数组，上面切片s1指向底层数组a，起始点为a1[1], 长度为2, 容量为4。当s1的值变化时，底层数组也变了，所以会看到a数组的值也变化了。
如果把`s1[0] = 6` 改为 `s1 = append(s1, 6, 7, 8)`，超过了切片的容量，切片开辟一块新的空间扩容，与原数组脱离关系，从而不会被改变。

### 3、随机到相同的值

```
package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println(rand.Intn(100))
}
```

发现每次都随机到相同的值，需要指定下随机种子：`rand.Seed(time.Now().UnixNano())`

