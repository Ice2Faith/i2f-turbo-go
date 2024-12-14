// CalculUse 运算，算术运算，逻辑运算，位运算
package main

import (
	"fmt"
)

//算术运算，必须是同种类型的数据才能进行运算，否则需要强转
func MathCalcul() {
	var a int = 12
	var b float32 = 12.125
	res := float32(a) + b
	fmt.Println("result=", res)
}

//和其他语言基本一致，不赘述
func LogicalCalcul() {
	var b bool = false
	b = !b
	a := true
	a = a && b
	a = a || b
	fmt.Println(a, b)
}

//和C语言类似，不赘述
func BitCalcul() {
	a := 12
	a = a & 3
	a = a | 5
	a = a ^ 0xa2
	a = a >> 2
	a = a << 1
	fmt.Println(a)
}
func main() {
	MathCalcul()
	LogicalCalcul()
	BitCalcul()
}
