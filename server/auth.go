package server

import "crypto/tls"
import "bufio"

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
	config := Load509Config()
	if config == nil {
		p.writeMessage(550, "Auth system cannot find X509KeyPair files")
		return
	}
	p.writeMessage(234, "AUTH command ok. Expecting TLS Negotiation.")

	p.theConnection = tls.Server(p.theConnection, config)
	p.reader = bufio.NewReader(p.theConnection)
	p.writer = bufio.NewWriter(p.theConnection)
	p.tls = true
}

func (p *Paradise) HandleProt() {
	p.writeMessage(200, "OK")
}

func (p *Paradise) HandlePbsz() {
	p.writeMessage(200, "OK")
}
