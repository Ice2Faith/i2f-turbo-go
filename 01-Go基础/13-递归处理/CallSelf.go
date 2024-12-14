// CallSelf
package main

import (
	"fmt"
)

//递归求阶乘
func Multiplys(number int) (result int) {
	//终止条件
	if number <= 1 {
		result = 1
		return
	}
	//递归调用
	return number * Multiplys(number-1)
}

//求斐波那契数列第N项，这里仅仅是介绍递归，实际上这种效率是比较低的
func FibonaciiSequeue(index int) (result int) {
	if index <= 2 {
		result = 1
		return
	}
	return FibonaciiSequeue(index-1) + FibonaciiSequeue(index-2)
}
func main() {
	fmt.Println("Hello World!")
	fmt.Println(Multiplys(5))
	fmt.Println(FibonaciiSequeue(5))
}
