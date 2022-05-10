package main

import "net"

type User struct {
	Username string
	Addr     string
	C        chan string
	coon     net.Conn
	Server   Server
}

//不断监听客户chan，获取message
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		u.coon.Write([]byte(msg + "\n"))
	}
}

//上线功能
func (u *User) Online(user *User) {
	u.Server.Maplock.Lock()
	u.Server.Useronline[user.Username] = user
	u.Server.Maplock.Unlock()
}

//下线功能
func (u *User) Offline() {
	u.Server.Maplock.Lock()
	delete(u.Server.Useronline, u.Username)
	u.Server.Maplock.Unlock()
}

//发送消息
func (u *User) DoMessage(mes string) {
	if mes == "" {
		return
	}
	mess := "[" + u.Addr + "]" + u.Username + "说:" + mes
	message := Message{Mess: mess, User: u}
	u.Server.Mes <- message
}

//在server端accapt一个请求后，创建一个User
func NewUser(coon net.Conn, server Server) *User {
	useraddr := coon.RemoteAddr().String()
	user := &User{
		Username: useraddr,
		Addr:     useraddr,
		C:        make(chan string),
		coon:     coon,
		Server:   server,
	}
	//不断监听User的chan
	go user.ListenMessage()
	return user
}
