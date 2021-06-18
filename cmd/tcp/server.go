package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

type User struct {
	ID             int
	Addr           string
	EnterAt        time.Time
	MessageChannel chan string
}

func main() {
	listener, err := net.Listen("tcp", ":2021")
	if err != nil {
		panic(err)
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConn(conn)
	}
}

// 1.新用户进来 2.用户普通消息 3, 用户离开
func broadcaster() {
	users := make(map[*User]struct{})

	for {
		select {
		case user := <-enteringChannel:
			// 新用户进入
			users[user] = struct{}{}
		case user := <-leavingChannel:
			// 用户离开
			delete(users, user)
			// 避免goroutine泄漏
			close(user.MessageChannel)
		case msg := <-messageChannel:
			for user := range users {
				user.MessageChannel <- msg
			}
		}
	}
}

var (
	enteringChannel = make(chan *User)
	leavingChannel  = make(chan *User)
	messageChannel  = make(chan string, 8)
)

func handleConn(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	// 1.新用户进来, 构建该用户的实例
	user := &User{
		ID:             rand.Int(),
		Addr:           conn.RemoteAddr().String(),
		EnterAt:        time.Now(),
		MessageChannel: make(chan string, 8),
	}

	// 2. 开一个 goroutine 用于写操作
	// 读写 goroutine 之间可以通过 channel 进行通信
	go sendMessage(conn, user.MessageChannel)

	// 3.给当前用户发送欢迎信息;给所有用户告知新用户的到来
	user.MessageChannel <- "Welcome, " + user.Addr
	messageChannel <- "user:`" + strconv.Itoa(user.ID) + "` has enter"

	// 4.将该记录到全局的用户列表中,避免用锁
	enteringChannel <- user

	// 5.循环读取用户的输入
	input := bufio.NewScanner(conn)
	for input.Scan() {
		messageChannel <- strconv.Itoa(user.ID) + ":" + input.Text()
	}

	if err := input.Err(); err != nil {
		log.Println("读取错误: ", err)
	}

	// 6.用户离开
	leavingChannel <- user
	messageChannel <- "user:`" + strconv.Itoa(user.ID) + "` has left"
}

func sendMessage(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		_, err := fmt.Fprintln(conn, msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func genUserID() int {
	id := genAutoIncreaseInt()
	return id()
}

func genAutoIncreaseInt() func() int {
	id := 1
	return func() int {
		id = id + 1
		return id
	}
}
