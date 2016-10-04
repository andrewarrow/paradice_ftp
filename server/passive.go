package server

import "fmt"
import "sync"
import "net"
import "strconv"
import "strings"
import "time"
import "crypto/tls"
import "math/rand"

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false
	case <-time.After(timeout):
		return true
	}
}

type Passive struct {
	listenSuccessAt int64
	listenFailedAt  int64
	closeSuccessAt  int64
	closeFailedAt   int64
	listenAt        int64
	connectedAt     int64
	connection      net.Conn
	command         string
	param           string
	cid             string
	port            int
	waiter          sync.WaitGroup
}

func (p *Paradise) closePassive(passive *Passive) {
	err := passive.connection.Close()
	if err != nil {
		passive.closeFailedAt = time.Now().Unix()
	} else {
		passive.closeSuccessAt = time.Now().Unix()
	}
}

func getThatPassiveConnection(passiveListen net.Listener, p *Passive) {
	var perr error
	p.connection, perr = passiveListen.Accept()
	if perr != nil {
		p.listenFailedAt = time.Now().Unix()
		p.waiter.Done()
		return
	}
	passiveListen.Close()
	// start reading from p.passive, it will block, wait for err. Err means client killed connection.
	p.listenSuccessAt = time.Now().Unix()
	p.waiter.Done()
}

func NewPassive(passiveListen net.Listener, cid string, now int64) *Passive {
	p := Passive{}
	p.cid = cid
	p.listenAt = now

	add := passiveListen.Addr()
	parts := strings.Split(add.String(), ":")
	p.port, _ = strconv.Atoi(parts[len(parts)-1])

	p.waiter.Add(1)
	p.listenFailedAt = 0
	p.listenSuccessAt = 0
	p.listenAt = time.Now().Unix()
	go getThatPassiveConnection(passiveListen, &p)

	return &p
}

func (p *Paradise) HandlePassive() {
	port := 50000 + rand.Intn(51009-50000)
	laddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+fmt.Sprintf("%d", port))
	var passiveListen net.Listener
	passiveListen, err = net.ListenTCP("tcp", laddr)
	if err != nil {
		p.writeMessage(550, "Error with passive: "+err.Error())
		return
	}

	if p.tls {
		config := Load509Config()
		if config == nil {
			p.writeMessage(550, "Passive connection could not load 509 Key Pairs")
			return
		}
		passiveListen = tls.NewListener(passiveListen, config)
	}

	cid := genClientID()
	passive := NewPassive(passiveListen, cid, time.Now().Unix())
	passive.command = p.command
	passive.param = p.param
	p.passive = passive

	if p.command == "PASV" {
		p1 := passive.port / 256
		p2 := passive.port - (p1 * 256)
		addr := p.theConnection.LocalAddr()
		tokens := strings.Split(addr.String(), ":")
		host := tokens[0]
		quads := strings.Split(host, ".")
		target := fmt.Sprintf("(%s,%s,%s,%s,%d,%d)", quads[0], quads[1], quads[2], quads[3], p1, p2)
		msg := "Entering Passive Mode " + target
		p.writeMessage(227, msg)
	} else {
		msg := fmt.Sprintf("Entering Extended Passive Mode (|||%d|)", passive.port)
		p.writeMessage(229, msg)
	}
}
