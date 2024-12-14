// CluserFunction
package main

import (
	"fmt"
)

/*
闭包函数：
返回函数的函数
常见的有next(),取代使用外部index来控制下一个的方式
*/
//返回他的内部函数，由于内部函数具有外部函数的局部变量，
//因此就可以实现，每次调用外部函数，都可以得到一个使用新的外部函数数据的内部函数
//就实现了一个类似全局变量的变量
func GetNext() func() string {
	//下面两个变量是在外部函数内的，但是对于内层函数来说，就和全局变量一样，
	//每个闭包之间相对独立
	var index int
	var heros = [...]string{"houyi", "dianwei", "zhaoyun", "makeboluo", "zhangfei"}
	//因此返回的内部函数，使用了闭包内部的变量，减少了全局变量的使用和内存的持续占用
	return func() string {
		res := heros[index]
		index = (index + 1) % len(heros)
		return res
	}
}

func main() {
	//下面的两次输出结果是完全一样的，因为它们拥有各自的index，在它们的外部函数内
	//并没有使用同一个index
	fmt.Println("Line1-----------")
	line1 := GetNext()
	for i := 0; i < 10; i++ {
		fmt.Println(line1())
	}
	fmt.Println("Line2-----------")
	line2 := GetNext()
	for i := 0; i < 10; i++ {
		fmt.Println(line2())
	}
}
