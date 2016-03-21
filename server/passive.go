package server

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func (p *Paradise) HandlePassive() {
	fmt.Println(p.ip, p.command, p.param)

	laddr, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
	passiveListen, _ := net.ListenTCP("tcp", laddr)

	add := passiveListen.Addr()
	parts := strings.Split(add.String(), ":")
	port, _ := strconv.Atoi(parts[len(parts)-1])

	p.waiter.Add(1)

	go func() {
		p.passiveConn, _ = passiveListen.AcceptTCP()
		passiveListen.Close()
		p.waiter.Done()
	}()

	if p.command == "PASV" {
		p1 := port / 256
		p2 := port - (p1 * 256)
		addr := p.theConnection.LocalAddr()
		tokens := strings.Split(addr.String(), ":")
		host := tokens[0]
		quads := strings.Split(host, ".")
		target := fmt.Sprintf("(%s,%s,%s,%s,%d,%d)", quads[0], quads[1], quads[2], quads[3], p1, p2)
		msg := "Entering Passive Mode " + target
		p.writeMessage(227, msg)
	} else {
		msg := fmt.Sprintf("Entering Extended Passive Mode (|||%d|)", port)
		p.writeMessage(229, msg)
	}
}
