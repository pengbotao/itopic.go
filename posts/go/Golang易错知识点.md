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

### 2、Slice陷阱


### 3、随机到相同的值
