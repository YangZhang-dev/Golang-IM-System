package main

import "net"

type User struct {
	Username string
	Addr     string
	C        chan string
	coon     net.Conn
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		u.coon.Write([]byte(msg + "\n"))
	}
}

func NewUser(coon net.Conn) *User {
	useraddr := coon.RemoteAddr().String()
	user := &User{
		Username: useraddr,
		Addr:     useraddr,
		C:        make(chan string),
		coon:     coon,
	}

	go user.ListenMessage()
	return user
}
