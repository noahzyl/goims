package main

import (
	"flag"
	"fmt"
	"github.com/noahzyl/goims/goims-client/client"
)

var serverIP string
var serverPort int

// ./client -ip <server ip> -port <server port>
func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "IP of server")
	flag.IntVar(&serverPort, "port", 5090, "Port of server")
}

func main() {
	flag.Parse() // Parse arguments passed by command line

	cli := client.NewClient("127.0.0.1", 5090)
	if cli == nil {
		fmt.Println(">>>>>Connect failed")
		return
	}
	fmt.Println(">>>>>Connect succeed")
	// Listen for response from server
	go cli.DealResponse()
	// Run client
	cli.Run()
}
