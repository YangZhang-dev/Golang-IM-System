package main

import "fmt"

func main() {
	fmt.Println("创建服务器中...")
	server := NewServer("127.0.0.1", 8888)
	server.Start()

}
