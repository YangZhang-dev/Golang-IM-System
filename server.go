package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

//创建一个Server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

func (s *Server) handler(conn net.Conn) {
	//进行操作
	fmt.Println("连接建立成功")
}

//启动一个server
func (s *Server) Start() {
	//new socket
	listen, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("new server error:")
		log.Fatal(err)
		return
	}
	//关闭socket
	defer listen.Close()
	for {
		//accept
		accept, err := listen.Accept()
		if err != nil {
			fmt.Println("accept error:")
			log.Fatal(err)
			continue
		}
		go s.handler(accept)
	}
}
