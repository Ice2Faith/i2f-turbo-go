// VarType
package main

/*
布尔类型：
	bool true false ,default = false
数值类型：
	int 整型 12 ,default =0
	float32 浮点 12.125 ,default =0.0
	complex64 复数 3+4i ,default =0
字符串类型：
	string "hello" ,default =""
复合类型：
	[10]int 数组 []float32 切片（变长数组）
	map[string]float 映射字典（键:值对） map[string]float32{"Size":12.125}
	struct 结构
	struct{
			name string
			age int
	}
	*int 整型指针 var pVal *int=&a ,default = nil
函数类型：
	func 函数类型
	var addFunc func()=func(a int,b int){
		return a+b
	}
接口类型：
	空接口类型，可以表示任意类型
	var Animal interface{}
	变量名首字母大写，包外部可见public，小写时不可见private
*/
import (
	"fmt"
)

func main() {
	// 通过var关键字声明变量，通过:=声明并赋值
	var isSuccess bool = true
	size := 12.155
	//%T 数据类型占位符 %v 值占位符
	fmt.Println("Var Type!", isSuccess, size)
}
