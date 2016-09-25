package server

import "fmt"

import "crypto/tls"
import "bufio"

//import "net"

func (p *Paradise) HandleUser() {
	p.user = p.param
	p.writeMessage(331, "User name ok, password required")
}

func (p *Paradise) HandlePass() {
	// think about using https://developer.bitium.com
	if AuthManager.CheckUser(p.user, p.param, &p.userInfo) {
		p.writeMessage(230, "Password ok, continue")
	} else {
		p.writeMessage(530, "Incorrect password, not logged in")
		p.theConnection.Close()
		delete(ConnectionMap, p.cid)
	}
}

func (p *Paradise) HandleAuth() {
	fmt.Println(p.param)
	p.writeMessage(234, "AUTH command ok. Expecting TLS Negotiation.")

	// openssl req -new -nodes -x509 -out server.pem -keyout server.key -days 3650 -subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=www.random.com/emailAddress=foo@foo.com"

	cert, cerr := tls.LoadX509KeyPair("server.pem", "server.key")
	if cerr != nil {
		fmt.Println(cerr)
		return
	}
	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		//ClientAuth:         tls.VerifyClientCertIfGiven,
		ClientAuth:         tls.NoClientCert,
		InsecureSkipVerify: true,
		ServerName:         "localhost"}
	tlsConn := tls.Server(p.theConnection, &config)
	fmt.Println("handshake1 ", tlsConn.ConnectionState())
	err := tlsConn.Handshake()
	if err == nil {
		p.theConnection = tlsConn
		p.writer = bufio.NewWriter(tlsConn)
		p.reader = bufio.NewReader(tlsConn)
		p.tls = true
	}
	fmt.Println("handshake2 ", tlsConn.ConnectionState())

}

func (p *Paradise) HandlePbsz() {
	fmt.Println(p.param)
	p.writeMessage(200, "PBSZ OK.")
}

func (p *Paradise) HandleProt() {
	fmt.Println(p.param)
	p.writeMessage(200, "PROT OK.")
}
