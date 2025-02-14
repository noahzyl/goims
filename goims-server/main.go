package main

import (
	"github.com/noahzyl/goims/goims-server/server"
)

func main() {
	s := server.NewServer("127.0.0.1", 5090)
	s.Start()
}
