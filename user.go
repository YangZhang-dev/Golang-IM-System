package main

import (
	"net"
)

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

//用户端返回给自己系统提示
func (u *User) SendMessage(message Message) {
	u.coon.Write([]byte(message.Mess))
}

//处理消息
func (u *User) DoMessage(mes string) {
	var message Message

	if mes == "" {
		return
	} else if mes == "who" {
		//如果是who则应该返回所有在UserOline的用户
		u.Server.Maplock.Lock()
		for _, user := range u.Server.Useronline {
			if user != u {
				message = GetSysMes("["+user.Username+"]"+"在线中\n", u)
				u.SendMessage(message)
			}
		}
		u.Server.Maplock.Unlock()
	} else if len(mes) > 7 && mes[:7] == "rename|" {
		//更改用户名
		newusername := mes[7:]
		_, ok := u.Server.Useronline[newusername]
		if ok {
			message = GetSysMes("用户名已经存在\n", u)

		} else {
			u.Server.Broatcast(u, "我改名为"+newusername)
			u.Server.Maplock.Lock()
			delete(u.Server.Useronline, u.Username)
			u.Server.Useronline[newusername] = u
			u.Server.Maplock.Unlock()
			u.Username = newusername
			message = GetSysMes("成功改名为:"+u.Username+"\n", u)
		}
		u.SendMessage(message)
	} else {
		//广播消息
		u.Server.Broatcast(u, mes)
	}

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

func GetSysMes(mes string, user *User) Message {
	return Message{Mess: mes, User: user}
}

func GetMessage(mes string, user *User) Message {
	mess := "[" + user.Addr + "]" + user.Username + "说:" + mes
	return Message{Mess: mess, User: user}
}
