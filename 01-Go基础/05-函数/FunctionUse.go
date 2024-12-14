// FunctionUse 函数的定义和使用
package main

import (
	"fmt"
)

//func 函数名(参数列表)(返回值列表){函数体}
func DoubleVoidFunc() {
	fmt.Println("无参数，无返回值的函数")
}

func ParamFunc(a int32, b int32) {
	fmt.Printf("带参数的函数，sum=%d", a+b)
}

func ReturnFunc() int { //声明返回值类型为int
	fmt.Println("带返回值的函数")
	return 0
}

//interface{}类型可以表示所有的类型，类似其他语言的Object
func ParamInterfaceFunc(a interface{}, b interface{}) {
	fmt.Println("使用空接口类型参数的函数,", a, b)
}

//两个返回值的函数
func ReturnMultiFunc() (int, float32) {
	fmt.Println("返回多个返回值的函数")
	return 10, 12.125
}

//返回参数预先定义，需要返回时，可以只return即可
func ReturnParamPreDefine(a int, b int) (sum int, avg float32) {
	fmt.Println("返回值参数预定义的函数")
	sum = a + b
	avg = float32(float32(a+b) / 2.0)
	return
}
func main() {
	DoubleVoidFunc()
	ParamFunc(3, 5)
	ret := ReturnFunc()
	fmt.Println(ret)
	ParamInterfaceFunc(10, "hello")
	a, b := ReturnMultiFunc()
	fmt.Println(a, b)
	sum, avg := ReturnParamPreDefine(10, 12)
	fmt.Println(sum, avg)
}
