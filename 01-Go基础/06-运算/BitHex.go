// BitHex 位运算的简易使用之常见进制转换
package main

import (
	"fmt"
)

func Pow2N(n int32) (result int32) {
	if n < 0 {
		fmt.Println("无法计算的次方，", n)
		return
	}
	result = 1 << n
	return
}
func Int32ToOther(num int32, hex int32) (result string) {
	move := 0
	switch hex {
	case 16:
		move = 4
	case 8:
		move = 3
	case 2:
		move = 1
	default:
		fmt.Println("无法处理的目标进制，", hex)
		return
	}
	//切片
	mapping := [...]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}
	result = ""
	for {
		if num == 0 {
			break
		}
		result += mapping[num&(hex-1)]
		num = num >> move
	}
	return
}
func main() {
	fmt.Println(Int32ToOther(12764, 16))
	fmt.Println(Pow2N(3))
}
