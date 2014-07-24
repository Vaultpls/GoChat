// gochat project main.go
package main

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"strings"
)

type Newuser struct {
	Chatroom map[net.Conn]string
	Name     map[net.Conn]string
	isMod    map[net.Conn]int
}

func newuser() *Newuser {
	return &Newuser{
		Chatroom: make(map[net.Conn]string),
		Name:     make(map[net.Conn]string),
		isMod:    make(map[net.Conn]int),
	}
}

//START BROADCAST

func (user *Newuser) ChatroomBroadcast(conn net.Conn, msg string) {
	for conns, chatroom := range user.Chatroom {
		if user.Chatroom[conn] == chatroom {
			fmt.Fprintf(conns, user.Name[conn]+": "+msg+"\n")
		}
	}
	fmt.Println(user.Chatroom[conn] + ": " + user.Name[conn] + ": " + msg)
}

func (user *Newuser) Broadcast(msg string) {
	for conns := range user.Chatroom {
		fmt.Fprintf(conns, msg)
	}
}

//END BROADCAST

//START COMMANDS
func (user *Newuser) CommandInterpretter(conn net.Conn, command string) {
	if user.isMod[conn] == 1 {
		if strings.HasPrefix(command, "/broadcast ") {
			newstring := strings.Replace(command, "/broadcast ", "", 1)
			user.Broadcast(newstring)
		}
	}
	if strings.HasPrefix(command, "/modme") {
		user.isMod[conn] = 1
	}
}

//END COMMANDS

func (user *Newuser) handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)
	for {
		line, err := tp.ReadLine()
		if err != nil {
			delete(user.Chatroom, conn)
			delete(user.Name, conn)
			break // break loop on errors
		}
		if strings.HasPrefix(line, "<>CHATROOM ") {
			newstring := strings.Replace(line, "<>CHATROOM ", "", 1)
			if len(newstring) < 40 {
				user.Chatroom[conn] = newstring
				fmt.Println(conn.RemoteAddr().String() + " Joined: " + newstring)
			}
		} else if strings.HasPrefix(line, "<>NAME ") {
			newstring := strings.Replace(line, "<>NAME ", "", 1)
			if len(newstring) < 25 {
				user.Name[conn] = newstring
				fmt.Println(conn.RemoteAddr().String() + " Name Change: " + newstring)
			}
		} else if strings.HasPrefix(line, "/") {
			user.CommandInterpretter(conn, line)
		} else {
			if user.Chatroom[conn] != "" || user.Name[conn] != "" || len(line) < 255 {
				user.ChatroomBroadcast(conn, line)
			}
		}
	}
}
func main() {
	fmt.Println("Hello World!")
	user := newuser()
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Trouble binding socket!")
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		go user.handleConnection(conn)
	}
}
