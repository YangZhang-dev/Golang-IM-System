package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type Message struct {
	Mess string
	User *User
}

type Server struct {
	Ip         string
	Port       int
	Mes        chan Message
	Maplock    sync.RWMutex
	Useronline map[string]*User
}

//创建一个不断监听Mes的goroutine
func (s *Server) ListenMessager() {
	for {
		messtruct := <-s.Mes
		mes := messtruct.Mess
		mesuser := messtruct.User
		s.Maplock.Lock()
		for _, cil := range s.Useronline {
			if cil != mesuser {
				cil.C <- mes
			}
		}
		s.Maplock.Unlock()
	}
}

//广播消息
func (s *Server) Broatcast(user *User, mes string) {
	if mes == "" {
		return
	}
	mess := "[" + user.Addr + "]" + user.Username + "说:" + mes
	message := Message{Mess: mess, User: user}
	s.Mes <- message
}

//服务端accept到一个客户端进行的操作
func (s *Server) Handler(conn net.Conn) {
	//创建一个用户
	user := NewUser(conn, *s)
	fmt.Printf("用户:%s上线\n", user.Username)
	//将用户加入在线队列中
	user.Online(user)
	//广播上线消息
	s.Broatcast(user, "我上线了")
	//接收用户的读入并进行广播
	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		//ctrl+c返回长度为0
		if n == 0 {
			s.Broatcast(user, "我下线了")
			user.Offline()
			return
		}
		if err != nil && err != io.EOF {
			fmt.Println("用户输入异常")
			return
		}
		//去除"\n"
		mes := string(buf[:n-1])
		//广播用户消息
		user.DoMessage(mes)
	}
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
		Mes:        make(chan Message),
		Useronline: make(map[string]*User),
	}
	return server
}
