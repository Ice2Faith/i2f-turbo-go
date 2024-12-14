// BaseGrammerGame
package main

import (
	"fmt"
)

//结构
type Person struct {
	name  string
	money int
}

//排序-冒泡排序
func sortPerson(peos []Person) {
	for i := 0; i < len(peos); i++ {
		swap := false
		for j := 1; j < len(peos)-i; j++ {
			if peos[j].money < peos[j-1].money {
				peos[j], peos[j-1] = peos[j-1], peos[j]
				swap = true
			}
		}
		if swap == false {
			break
		}
	}
}
func main() {

	//切片，可以理解为集合
	var Peoples []Person
	for {
		var name string
		var money int
		//获取输入
		fmt.Println("Please input name end of #:")
		fmt.Scanf("%s\n", &name) //这里使用\n强制读取一行

		if name == "#" {
			break
		}

		fmt.Println("Please input money end of #:")
		fmt.Scanf("%d\n", &money)
		//加入切片
		Peoples = append(Peoples, Person{name, money})
	}
	//输出排序前
	fmt.Println(Peoples)
	//排序
	sortPerson(Peoples)
	//输出排序后
	fmt.Println(Peoples)
}
