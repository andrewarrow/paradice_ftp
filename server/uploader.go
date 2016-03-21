package server

import (
	"fmt"
	"io"
)

func (p *Paradise) HandleStore() {
	fmt.Println(p.ip, p.command, p.param)

	p.writeMessage(150, "Data transfer starting")

	_, err := p.storeOrAppend()
	if err == io.EOF {
		p.writeMessage(226, "OK, received some bytes") // TODO send total in message
	} else {
		p.writeMessage(550, "Error with upload: "+err.Error())
	}
}

func (p *Paradise) storeOrAppend() (int64, error) {
	var err error
	err = p.readFirst512Bytes()
	if err != nil {
		return 0, err
	}

	// TODO run p.buffer thru mime type checker
	// if mime type bad, reject upload

	// TODO send p.buffer to where u want bits stored

	var total int64
	var n int
	total = int64(len(p.buffer))
	for {
		temp_buffer := make([]byte, 20971520) // reads 20MB at a time
		n, err = p.passiveConn.Read(temp_buffer)
		total += int64(n)

		if err != nil {
			break
		}
		// TODO send temp_buffer to where u want bits stored
		if err != nil {
			break
		}
	}

	return total, err
}

func (p *Paradise) readFirst512Bytes() error {
	p.buffer = make([]byte, 0)
	var err error
	p.waiter.Wait()
	for {
		temp_buffer := make([]byte, 512)
		n, err := p.passiveConn.Read(temp_buffer)

		if err != nil {
			break
		}
		p.buffer = append(p.buffer, temp_buffer[0:n]...)

		if len(p.buffer) >= 512 {
			break
		}
	}

	if err != nil && err != io.EOF {
		return err
	}

	// you have a buffer filled to 512, or less if file is less than 512
	return nil
}
