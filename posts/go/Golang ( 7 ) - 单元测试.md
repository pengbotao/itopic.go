```
{
    "url": "go-testing",
    "time": "2021/02/01 21:01",
    "tag": "Golang",
    "toc": "yes",
    "public": "no"
}
```

Go语言在工具链和标准库中提供了对测试的原生支持，`go test`会执行`_test.go`中复合TestXxx命名规则的测试函数。比如创建一个demo.go

```
package demo

var times = 2

func Mul(a, b int) int {
	return a * b * times
}
```

然后再同目录创建一个`demo_test.go`：

```
package demo

import (
	"testing"
)

func TestMul(t *testing.T) {
	tests := map[string]struct {
		a, b   int
		expect int
	}{
		"case1": {a: 1, b: 2, expect: 4},
		"case2": {a: 1, b: 0, expect: 0},
		"case3": {a: 2, b: 3, expect: 12},
	}
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			if m := Mul(v.a, v.b); m != v.expect {
				t.Errorf("Mul(%d, %d) = %d Error, Expect %d", v.a, v.b, m, v.expect)
			}
		})
	}
}
```

测试函数的签名需要是`func TestXXX( t *testing.T )`，通常编辑器上都有直接运行测试用例，或者直接手动执行：

```
$ go test -run ^TestMul$ . -v
=== RUN   TestMul
=== RUN   TestMul/case2
=== RUN   TestMul/case3
=== RUN   TestMul/case1
--- PASS: TestMul (0.00s)
    --- PASS: TestMul/case2 (0.00s)
    --- PASS: TestMul/case3 (0.00s)
    --- PASS: TestMul/case1 (0.00s)
PASS
ok      demo    0.199s
```

`-run`：指定测试用例的名称，支持正则表达式，可不指定该参数

`-v`：显示详细输出

**一、包内/外测试**

根据测试代码与被测试包的位置关系可以分为包内测试和包外测试，包内测试是指将测试代码放在与被测包同名的包中，而包外测试则是将代码放在名为被测包名+“_test”的包中。他们之间的区别与联系是：

- 包内测试可以访问该包下的所有符号，无论是导出的还是未导出；包外测试只能访问被导出的符号。
- 包内测试本质上是一种面相实现的白盒测试，覆盖面更广，维护成本更高；包外测试则是一种面向接口的黑盒测试，可能存在测试盲区，但接口稳定性更高，
- 包内可能存在循环引用的硬伤，包外测试可以解决。

使用包外测试时如果想访问非导出的符号，可以通过在包内的测试文件来中转，该文件位于被测包名下，仅仅用来将被测包的内部符号在测试阶段暴露给包外测试代码，如：

```
//demo_test.go

package demo

var Times = times
```

**二、代码组织方式**

第一种方式：平铺

```
func TextXx1(t *testing.T) {

}

func TextXx2(t *testing.T) {

}
```

第二种方式：层级，如文章开头的示例。

```
func testXx1(t *testing.T) {

}

func testXx2(t *testing.T) {

}

func TextXxx(t *testing.T) {
	t.Run("TestXx1", testXx1)
	t.Run("TestXx2", testXx2)
}
```

**三、测试固件**

测试固件是指一个人造的、确定性的环境，一个测试用例或一个测试集。比如给一个测试函数创建和销毁。

```
func setUp(name string) func() {
	fmt.Println("setUp for ", name)
	return func() {
		fmt.Println("tearDown for ", name)
	}
}

func TestXx1(t *testing.T) {
	defer setUp("TestXx1")()
	fmt.Println("Run Test For TestXx1")
}

---------------------------------
=== RUN   TestXx1
setUp for  TestXx1
Run Test For TestXx1
tearDown for  TestXx1
--- PASS: TestXx1 (0.00s)
PASS
ok  	demo	0.413s
```

1.14版本增加了testing.Cleanup方法上面可以改写为：

```
func TestXx1(t *testing.T) {
	//defer setUp("TestXx1")()
	t.Cleanup(setUp("TestXx1"))
	fmt.Println("Run Test For TestXx1")
}
```

有时候我们需要将所有测试函数放在一个更大范围的测试固件环境中执行，这就是包级别测试固件。在1.4中增加了TestMain函数，使得可以创建和销毁包级别测试固件。

```
func setUp(name string) func() {
	fmt.Println("setUp for ", name)
	return func() {
		fmt.Println("tearDown for ", name)
	}
}

func TestXx1(t *testing.T) {
	defer setUp("TestXx1")()
	fmt.Println("Run Test For TestXx1")
}

func pkgSetUp(name string) func() {
	fmt.Println("package setUp for ", name)
	return func() {
		fmt.Println("package tearDown for", name)
	}
}

func TestMain(m *testing.M) {
	defer pkgSetUp("demo_test")()
	m.Run()
}

---------------------------------
package setUp for  demo_test
=== RUN   TestXx1
setUp for  TestXx1
Run Test For TestXx1
tearDown for  TestXx1
--- PASS: TestXx1 (0.00s)
PASS
package tearDown for demo_test
ok  	demo	0.423s
```

**基准测试**

可以像对普通单元测试那样再"_test.go"文件中创建被测对象的性能基准测试，以Benchmark前缀开头，如：

```
func BenchmarkMul(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Mul(1, 2)
	}
}
```

上面是基于顺序执行的性能基准测试，如果想执行并行执行的基准测试，可以按如下操作：

```
func BenchmarkXxx(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			Mul(1, 2)
		}
	})
}
```
