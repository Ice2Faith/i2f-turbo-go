// PackagePath
/*
实现代码的复用，可以将复用代码放到GOPATH环境变量中，否则不可使用（绝对路径）
一个包就是一个文件夹，包内放置包的相关代码文件，文件可以有多个
工具，GOPATH管理，添加源代码所在的绝对路径
源代码需要方法放在：GOPATH路径/src/包文件夹/源代码文件
因此GOPATH之下需要有一个src文件夹存放包（文件夹），包内存放源代码（go源代码文件）
*/
package main

import (
	"fmt"
	"packutil"
)

func main() {
	fmt.Println(packutil.IsLeapYear(1999))
	fmt.Println(getNumbers(52219))
}

func getNumbers(number int) (result [5]int) {
	temp := number
	i := 0
	for {
		if temp == 0 {
			break
		}
		result[4-i] = temp % 10
		temp /= 10
		i++

	}
	return
}
