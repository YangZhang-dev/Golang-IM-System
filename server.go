package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	Mes        chan string
	Maplock    sync.RWMutex
	Useronline map[string]*User
}

//创建一个不断监听Mes的goroutine
func (s *Server) ListenMessager() {
	for {
		mes := <-s.Mes
		s.Maplock.Lock()
		for _, cil := range s.Useronline {
			cil.C <- mes
		}
		s.Maplock.Unlock()
	}
}

//广播消息
func (s *Server) Broatcast(user *User, mes string) {
	mess := "[" + user.Addr + "]" + user.Username + mes
	s.Mes <- mess
}

//服务端accept到一个客户端进行的操作
func (s *Server) Handler(conn net.Conn) {
	//创建一个用户
	user := NewUser(conn)
	fmt.Printf("用户:%s上线\n", user.Username)
	//将用户加入在线队列中
	s.Maplock.Lock()
	s.Useronline[user.Username] = user
	s.Maplock.Unlock()
	//广播上线消息
	s.Broatcast(user, "上线了")
	select {}
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
	fmt.Println("服务器创建成功!!!")
	//关闭socket
	defer listen.Close()
	//在创建一个服务器的同时就应该使其启动监听MES队列的方法
	go s.ListenMessager()

	for {
		//accept等待连接
		accept, err := listen.Accept()
		if err != nil {
			fmt.Println("accept error:")
			log.Fatal(err)
			continue
		}
		go s.Handler(accept)
	}
}

//创建一个Server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:         ip,
		Port:       port,
		Mes:        make(chan string),
		Useronline: make(map[string]*User),
	}
	return server
}
