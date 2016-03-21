package main

import (
	"fmt"
	"net"

	"github.com/andrewarrow/paradise_ftp/server"
)

func main() {

	url := fmt.Sprintf("localhost:%d", 2121) // change to 21 in production
	var listener net.Listener
	listener, err := net.Listen("tcp", url)

	if err != nil {
		fmt.Println("cannot listen on: ", url)
		return
	}
	fmt.Println("listening on: ", url)

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("listening error ", err)
			break
		}
		p := server.NewParadise(connection)

		go p.HandleCommands()
	}
}
