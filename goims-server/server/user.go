package server

import (
	"net"
	"strings"
)

type User struct {
	UserName    string
	Addr        string
	UserChannel chan string
	conn        net.Conn
	server      *Server
}

// Create a user
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	return &User{
		UserName:    userAddr,
		Addr:        userAddr,
		UserChannel: make(chan string),
		conn:        conn,
		server:      server,
	}
}

func (this *User) Online() {
	// Add current user to server's OnlineMap
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.UserName] = this
	this.server.mapLock.Unlock()
	this.server.Broadcast(this, "online")
}

func (this *User) Offline() {
	// Remove current user from server's OnlineMap
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.UserName)
	this.server.mapLock.Unlock()
	// Release resources
	defer close(this.UserChannel)
	defer this.conn.Close()
	this.server.Broadcast(this, "offline")
}

// Get user's message
func (this *User) ReadMessage(buffer []byte) (int, error) {
	n, err := this.conn.Read(buffer)
	return n, err
}

// Send a message to relative client
func (this *User) SendMessage(msg string) {
	this.conn.Write([]byte(msg + "\n"))
}

// Handle user's message
func (this *User) DoMessage(msg string) {
	if msg == "online users" { // Get all online users
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.UserName + ": online"
			this.SendMessage(onlineMsg)
		}
		this.server.mapLock.Unlock()
		return
	} else if len(msg) > 7 && msg[:7] == "rename|" { // Update username
		// Format: rename <new username>
		newName := strings.Split(msg, "|")[1]
		if newName == "" {
			this.SendMessage("User name can not be empty")
			return
		}
		// Check whether new username has been used
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMessage("New username has been used")
			return
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.UserName)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()
			this.UserName = newName
			this.SendMessage("Username has been updated as " + this.UserName)
			return
		}
	} else if len(msg) > 3 && msg[:3] == "to|" { // Private chat
		// Format: to <target username> <content>
		targetUserName := strings.Split(msg, "|")[1]
		if targetUserName == "" {
			this.SendMessage("Username can not be empty")
			return
		}
		targetUser, ok := this.server.OnlineMap[targetUserName]
		if !ok { // Target user does not exist
			this.SendMessage("Target user does not exist")
			return
		}
		words := strings.Split(msg, "|")[2:]
		content := strings.Join(words, "|")
		targetUser.SendMessage("[" + this.Addr + "]" + this.UserName + " to you: " + content)
		return
	} else {
		this.server.Broadcast(this, msg)
		return
	}
}

// Listening user messages and send messages to relative client
func (this *User) ListenMessage() {
	for {
		msg := <-this.UserChannel // Get a message from server
		this.SendMessage(msg)     // Send the message to client
	}
}
