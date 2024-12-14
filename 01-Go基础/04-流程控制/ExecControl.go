// ExecControl 流程控制语句
package main

import (
	"fmt"
)

func IfSegment() {
	fmt.Println("If Seg")
	var a float32 = 12.125 //变量的定义方式一
	if a > 10 {            //判断条件不用带小括号
		fmt.Println("a>10") //即使语句中只有一条语句，花括号也不能少
	} else if a == 10 { //左半花括号只能和条件在同一行
		fmt.Println("a<=10")
	} else {
		fmt.Println("a<10")
	}
}

func SwitchSegment() {
	fmt.Println("Switch Seg")
	var a int32 = 12
	switch a { //同样，条件不用括号
	case 12: //每个case自带C语言中的break，因此已经省略
		fallthrough //穿透这个case继续执行下一个case
	case 11:
		fallthrough
	case 10:
		fallthrough
	case 9:
		fmt.Println("a==9")
	case 8:
		fmt.Println("a==8")
	default:
		fmt.Println("cannot do switch")
	}
}
func ForSegment() {
	fmt.Println("For Seg")
	var a int = 3
	for { //go中没有while等循环，只有for循环（增强），无条件就是while（true）
		a--
		fmt.Println("while loop")
		if a == 0 {
			break
		}
	}
	//传统的for方式，三个表达式
	for i := 0; i < 5; i++ { //变量定义方式二，使用:=定义并初始化，类型由系统识别
		fmt.Println("trad for loop")
	}
}
func main() {
	IfSegment()
	SwitchSegment()
	ForSegment()
}
