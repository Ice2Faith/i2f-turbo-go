// UnifyTools
package main

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe" //测量内存占用大小的系统包
)

func main() {
	SizeOfTest()
	SienceNumber()
	StringJoin()
	BigToSmallLost()
	DefineConstImmi()
	FmtFunctions()
	BaseTypeAndStringConvert()
	TypeDef2OtherName()
	SimpleOperator()
}

//内存占用测定
func SizeOfTest() {
	fmt.Println(unsafe.Sizeof(1.26))
	fmt.Println(unsafe.Sizeof(12))
	fmt.Println(unsafe.Sizeof('A'))
	fmt.Println(unsafe.Sizeof(true))
	fmt.Println(unsafe.Sizeof("hello"))
}

//科学计数法
func SienceNumber() {
	var test int32 = 1e8
	fmt.Println(test)
	fmt.Println(3.14e-3)
	//e的大小写仅仅影响输出的时候的e的大小写
	fmt.Printf("%e,%E\n", 3.157, 3.157)
}

//字符串拼接
func StringJoin() {
	var str = "hello"
	var str1 = "go"
	str2 := str + " " + str1
	fmt.Println(str2)
	//字符串拼接
	var heros = []string{"ahri", "alis", "darwen"}
	res := strings.Join(heros, "--")
	fmt.Println(res)
}

//强制类型转换精度丢失
func BigToSmallLost() {
	var a uint8 = 125
	fmt.Println(int16(a))
	var b int32 = 12345
	fmt.Println(int8(b)) //数据精度丢失
}

//常量自动增长
func DefineConstImmi() {
	const (
		Zero = iota //类似枚举的默认值，自动增长，iota处为0，实现十进制0 1 2 3 4 5
		One
		Two
		Three
		Four
		Five
	)
	fmt.Println(Zero, One, Two, Three, Four, Five)
	const (
		D = 1 << iota //利用左移运算（1左移0..位），实现1 10 100 1000 10000的二进制实现1 2 4 8 46
		C
		B
		A
		S
	)
	fmt.Println(D, C, B, A, S)
}

//fmt的sprint系列应用
func FmtFunctions() {
	//使用sprint系列函数，产生格式化字符串
	str := fmt.Sprint(12, "nihao", 12.125)
	fmt.Println(str)
}

//基本数据类型和字符串的相互转换
func BaseTypeAndStringConvert() {
	//使用strconv实现互转
	//使用format系列，将数据转换成字符串
	stri := strconv.FormatInt(1285, 2) //参数：数据，目标进制
	fmt.Println(stri)
	strb := strconv.FormatBool(true)
	fmt.Println(strb)
	//精度为5的浮点类型（f/e/E/b/G）输出64位浮点数
	strf := strconv.FormatFloat(12.125, 'f', 5, 64)
	fmt.Println(strf)

	//使用parse系列，将字符串转换成数据
	bvalue, berr := strconv.ParseBool("true")
	fmt.Println(bvalue, berr)
	//参数：字符串，源进制，数据大小（32位够用）
	ivalue, ierr := strconv.ParseInt("100100", 2, 32)
	fmt.Println(ivalue, ierr)
	fvalue, ferr := strconv.ParseFloat("3.14", 32)
	fmt.Println(fvalue, ferr)
	//转换失败案例
	evalue, eerr := strconv.ParseBool("turely")
	if eerr == nil {
		fmt.Println(evalue)
	} else { //转换失败的执行线
		fmt.Println(eerr)
	}

	//更加便捷的整型和字符串转换方式(itoa==int to array)
	cival, cierr := strconv.Atoi("1257")
	fmt.Println(cival, cierr)
	cstr := strconv.Itoa(cival)
	fmt.Println(cstr)
}

//类型的别名
func TypeDef2OtherName() {
	type Enum = int8
	var Orange Enum = 2
	fmt.Printf("%T,%v", Orange, Orange)
	type _person struct {
		name string
		age  int32
	}
	//通过type关键字实现类型的别名定义
	type People = _person
	var p = People{
		name: "karwen",
		age:  20,
	}
	fmt.Println(p)
}

//简化的运算符号
func SimpleOperator() {
	var a int = 10
	a += 2
	a -= 3
	a *= 2
	a /= 2
	a %= 4
	a &= 11
	a |= 32
	a ^= 19
	a <<= 2
	a >>= 3
	fmt.Println(a)
}
