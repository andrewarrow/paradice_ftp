package server

import (
	//"fmt"
	"strings"
)

func (p *Paradise) HandleSyst() {
	p.writeMessage(215, "UNIX Type: L8")
}

func (p *Paradise) HandlePwd() {
	p.writeMessage(257, "\""+p.path+"\" is the current directory")
}

func (p *Paradise) HandleType() {
	p.writeMessage(200, "Type set to binary")
}

func (p *Paradise) HandleQuit() {
	p.writeMessage(221, "Goodbye")
	p.theConnection.Close()
}

func (p *Paradise) HandleCwd() {
	if p.param == ".." {
		p.path = "/"
	} else {
		p.path = p.param
	}
	if !strings.HasPrefix(p.path, "/") {
		p.path = "/" + p.path
	}
	p.writeMessage(250, "CD worked")
}

func (p *Paradise) HandleSize() {
	p.writeMessage(450, "downloads not allowed")
}

func (p *Paradise) HandleRetr() {
	p.writeMessage(551, "downloads not allowed")
}

func (p *Paradise) HandleStat() {
	p.writeMessage(551, "downloads not allowed")
}
