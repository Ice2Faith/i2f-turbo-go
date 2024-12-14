// ConRun
package main

// 在Go中，把线程换了一个概念为协程，是一个微线程，创建和删除线程更快
// 包含线程睡眠包time
import (
	"fmt"
	"time"
)

//协程函数，协程内串行执行，协程之间并发执行
func doThing(a int) {
	fmt.Println("thread output:", a)
}
func main() {
	//并发100条协程（微线程）直接通过go 关键字
	for i := 0; i < 100; i++ {
		go doThing(i)
	}
	//主协程
	fmt.Println("main output")
	//协程睡眠,如果主协程结束，子协程也会立即结束（可能还没结束）
	time.Sleep(100 * time.Millisecond)
}
