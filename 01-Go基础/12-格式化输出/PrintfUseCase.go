// PrintfUseCase
package main

import (
	"fmt"
)

/*
printf占位符：
	通用：
	%v:值的默认格式输出，简单认为是输出值
	&+v:类似%v,但输出结构时会带上字段名
	%#v:值的GO语法表示
	%T:值的类型的GO表示，简单认为是输出类型
	%%：转义，输出一个%号
	布尔值：
	%t:输出布尔值true,false
	整数值：
	%b:整数的二进制输出
	%c:值对应的UNICODE编码
	%d:十进制输出
	%o:八进制输出
	%x:小写形式的十六进制输出
	%X:大写形式的十六进制输出
	%q:值对应的单引号包裹的字符值，会采用安全的转义表示
	%U:表示为UNICODE格式，比如：U+1234==U+%04X
	浮点|复数值：
	%b:无小数部分，二进制指数的科学计数法
	%e:科学计数法，小e表示
	%E:科学计数法，大E表示
	%f/%F:单纯的浮点数，有小数部分，无指数部分
	%g/%G:根据数字的情况，选择合适的输出方式（%e/%f方式输出）
	字符串和[]byte:
	%s:直接输出字符串或者[]byte
	%q:值对应的字符串字面值，会自动转义
	%x/%X:每个字节对应的十六进制数输出
	指针：
	%p:标识为16进制，并有前导符号0x
	格式控制：
	%-m.nf:输出总共宽为m，小数部分占n位，左对齐方式（注意，小数点占一位，整数部分过长，则域宽会自动跳帧==调整）

*/
func main() {
	a := 12
	b := 12.1254
	c := "25559"
	d := 0.0000000054178
	e := false
	f := &a
	fmt.Printf("%v\n", a)
	fmt.Printf("%v\n", c)
	fmt.Printf("%t\n", e)
	fmt.Printf("%t\n", b)
	fmt.Printf("%d\n", a)
	fmt.Printf("%#v\n", e)
	fmt.Printf("%e\n", d)
	fmt.Printf("%g\n", d)
	fmt.Printf("%s\n", c)
	fmt.Printf("%p\n", f)
	fmt.Printf("%b\n", a)
}
