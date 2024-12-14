// LoopControl
package main

import (
	"fmt"
)

func embedLoop() {
	//以99乘法表输出为例
	for i := 1; i < 10; i++ {
		for j := 1; j <= i; j++ {
			fmt.Printf("%d*%d=%02d ", i, j, i*j)
		}
		fmt.Println()
	}
}
func breakPointLoop() {
	//跳出指定循环
outer:
	for i := 1; i < 10; i++ {
		for j := 1; j < 10; j++ {
			if i == 5 && j == 5 {
				break outer
			}
			fmt.Print("* ")
		}
		fmt.Println()
	}
}
func main() {
	embedLoop()
	breakPointLoop()
}
