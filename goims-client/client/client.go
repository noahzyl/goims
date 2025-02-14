package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

type Client struct {
	ServerIP   string
	ServerPort int
	conn       net.Conn
	option     int
	scanner    *bufio.Scanner // Read from standard input
}

// Create a client
func NewClient(serverIP string, serverPort int) *Client {
	server := &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
		option:     -1,
		scanner:    bufio.NewScanner(os.Stdin),
	}
	// Connect to server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	server.conn = conn
	return server
}

// Deal with response from server
func (this *Client) DealResponse() {
	// Listen for messages from the server, and copy the message to standard output
	io.Copy(os.Stdout, this.conn)
}

// Read a line from standard input
func (this *Client) ReadLine() string {
	if this.scanner.Scan() {
		line := this.scanner.Text()
		// Remove '\n'
		line = strings.TrimRight(line, "\n")
		return line
	}
	return ""
}

// Send a message to server
func (this *Client) SendMessage(msg string) bool {
	_, err := this.conn.Write([]byte(msg + "\n"))
	if err != nil {
		fmt.Println(">>>>>conn.Write error:", err)
		return false
	}
	return true
}

// Show menu
func (this *Client) Menu() bool {
	var op int
	fmt.Println(">>>>>Menu<<<<<")
	fmt.Println("1. Public chat")
	fmt.Println("2. Private chat")
	fmt.Println("3. View all online users")
	fmt.Println("4. Update username")
	fmt.Println("0. Quit")
	fmt.Println(">>>>>Input an option:")

	fmt.Scanln(&op) // Read option

	if op >= 0 && op <= 4 {
		this.option = op
		return true
	} else {
		fmt.Println(">>>>>Invalid option, please try again")
		return false
	}
}

// Public chat mode
func (this *Client) PublicChat() {
	var msg string
	fmt.Println(">>>>>Public chat mode, use \"exit\" go back to the previous level")
	fmt.Println(">>>>>Message:")
	// fmt.Scanln(&msg) // Error: fmt.Scanln() will ignore spaces in user's input
	msg = this.ReadLine()
	// Receive user's input and send to the server
	for msg != "exit" {
		// Send the message when the message is not empty
		if len(msg) > 0 {
			ok := this.SendMessage(msg)
			if !ok {
				return
			}
		}
		fmt.Println(">>>>>Message:")
		// fmt.Scanln(&msg)
		msg = this.ReadLine()
	}
}

// Private chat mode
func (this *Client) PrivateChat() {
	var targetUserName string
	var msg string
	fmt.Println(">>>>>Private chat mode, use \"exit\" go back to the previous level")
	this.ViewOnline()
	fmt.Println(">>>>>Input a username:")
	// fmt.Scanln(&targetUserName)
	targetUserName = this.ReadLine()
	for targetUserName != "exit" {
		fmt.Println(">>>>>Message:")
		// fmt.Scanln(&msg)
		msg = this.ReadLine()
		for msg != "exit" {
			if len(msg) > 0 {
				sendMsg := "to|" + targetUserName + "|" + msg
				ok := this.SendMessage(sendMsg)
				if !ok {
					return
				}
			}
			fmt.Println(">>>>>Message:")
			// fmt.Scanln(&msg)
			msg = this.ReadLine()
		}
		this.ViewOnline()
		fmt.Println(">>>>>Input a username:")
		// fmt.Scanln(&targetUserName)
		targetUserName = this.ReadLine()
	}
}

// View all online users
func (this *Client) ViewOnline() {
	msg := "online users"
	this.SendMessage(msg)
}

// Update username
func (this *Client) UpdateName() {
	var newName string
	fmt.Println(">>>>>>Please input your new username:")
	// fmt.Scanln(&newName)
	newName = this.ReadLine()
	if newName == "" {
		fmt.Println("Username can not be empty")
		return
	}
	msg := "rename|" + newName
	this.SendMessage(msg)
}

// Run the client and execute transactions
func (this *Client) Run() {
	for this.option != 0 {
		for !this.Menu() {
			// Continue looping until user input a valid option
		}
		switch this.option {
		case 1:
			this.PublicChat()
			break
		case 2:
			this.PrivateChat()
			break
		case 3:
			this.ViewOnline()
			break
		case 4:
			this.UpdateName()
			break
		}
	}
}
