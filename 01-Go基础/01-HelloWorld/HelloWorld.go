// HelloWorld
// 包名为main才会执行，否则不会执行
package main

// 标准（格式化）输入输出 导入包 format
import (
	"fmt"
)

// 主函数，行没有分号结束，格式化格式类似C
func main() {
	fmt.Println("Hello,What's your name?developer.")
	var name string
	fmt.Scanf("%s", &name)
	fmt.Println("Hello " + name + ",Welcome to Golang.")
}
