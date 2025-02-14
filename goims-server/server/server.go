package server

import (
	"fmt"
	"io"
	"net"
	"runtime"
	"sync"
	"time"
)

type Server struct {
	IP        string
	Port      int
	OnlineMap map[string]*User // Map of online users (critical resource)
	mapLock   sync.Mutex
	Message   chan string
}

// Create a Server
func NewServer(ip string, port int) *Server {
	return &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

// Broadcast a message
func (this *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.UserName + ":" + msg
	this.Message <- sendMsg
}

// Listen message and send the message to all users
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		// Send the message to all users
		this.mapLock.Lock()
		for _, user := range this.OnlineMap {
			user.UserChannel <- msg
		}
		this.mapLock.Unlock()
	}
}

// Listen a message from a user, and broadcast the message
func (this *Server) ListenUserMessage(user *User, isAlive chan bool) {
	buf := make([]byte, 4096)
	for {
		n, err := user.ReadMessage(buf)
		if err != nil && err != io.EOF {
			fmt.Println("user.ReadMessage error:", err)
			return
		}
		if n == 0 {
			user.Offline()
			return
		}
		// Get the message, and remove the '\n'
		msg := string(buf[:n-1])
		// User handles the message
		user.DoMessage(msg)
		isAlive <- true
	}
}

// Handler
func (this *Server) Handler(conn net.Conn) {
	user := NewUser(conn, this)
	go user.ListenMessage() // Start to listening user's messages
	// Online
	user.Online()
	// Listen user's messages
	isAlive := make(chan bool) // Listen whether current user is alive
	go this.ListenUserMessage(user, isAlive)
	// Block current handler
	for {
		select {
		case <-isAlive:
		// Current user is alive
		// Do not do anything...
		case <-time.After(time.Second * 300): // Timer (5 minutes)
			user.SendMessage("Your session has timed out")
			user.Offline()
			runtime.Goexit()
		}
	}
}

// Start a server
func (this *Server) Start() {
	// Set a listening socket
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.IP, this.Port))
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}
	defer listener.Close() // Close listening socket after server exited
	go this.ListenMessage()
	// Continue to listening and handling requests
	for {
		// Accept a user's request and set a connection socket
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept error:", err)
			continue
		}
		// Use a goroutine to do a handler
		go this.Handler(conn)
	}
}
