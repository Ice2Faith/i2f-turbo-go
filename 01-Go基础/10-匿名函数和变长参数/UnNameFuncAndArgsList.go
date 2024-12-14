// UnNameFuncAndArgsList
package main

//GO中也支持使用goto，但是依旧不建议使用
import (
	"fmt"
	"time"
)

func main() {
	UnNameFunc()
	ArgsListFunc()
}

//匿名函数
func UnNameFunc() {
	//匿名函数，带参数的传值，通过函数之后加上函数参数列表的形式
	go func(a, b int) {
		fmt.Println(a + b)
	}(12, 32)
	//将匿名函数保存在变量中
	var f1 = func(a, b int) (ret int) {
		ret = (a + b)
		return
	}

	ret := f1(12, 24)
	fmt.Println(ret)
	fmt.Printf("%T,%v\n", f1, f1)
	time.Sleep(5 * time.Second) //等待协程结束
}

//变长参数
func ArgsListFunc() {
	test(1, 2, 3, 4, 5)
	test2(1, 2.2, "abc")
}

func test(nums ...int) {
	fmt.Println(nums)
}
func test2(args ...interface{}) {
	fmt.Println(args)
	//变长参数的遍历，利用range，以便获得下标
	for index, value := range args {
		fmt.Println(index, value)
	}
}
