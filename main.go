package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/andrewarrow/paradise_ftp/server"
)

var port int

func init() {
	flag.IntVar(&port, "port", 2121, "port to listen on")
}

func main() {
	flag.Parse()

	url := fmt.Sprintf("localhost:%d", port) // change to 21 in production
	var listener net.Listener
	listener, err := net.Listen("tcp", url)

	if err != nil {
		fmt.Println("cannot listen on:", url)
		return
	}
	fmt.Println("listening on:", url)

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
