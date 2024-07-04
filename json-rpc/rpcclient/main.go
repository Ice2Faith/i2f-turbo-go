package main

import (
	"fmt"
	"rpcclient/gorpc"
)

func main() {
	client := gorpc.GetDefaultRpcClient("127.0.0.1")
	client.Run()

	var req string
	var resp string
	for {
		fmt.Println("$> 请输入发送的内容，输入#退出：")
		fmt.Print(">/ ")
		fmt.Scanln(&req)
		if req == "#" {
			break
		}

		err := client.Call("hello.EchoBack", req, &resp)
		if err != nil {
			fmt.Printf("$> rpc invoke err: %v\n", err)
		}
		fmt.Println("$> rpc resp:", resp)
	}

	client.Stop()
}
