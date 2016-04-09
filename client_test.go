package main

import "testing"
import "os"
import "time"
import "paradise/server"
import "paradise/client"

var file *os.File
var fileBytes []byte

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestSimple(t *testing.T) {
	go server.Start()
	time.Sleep(1 * (time.Second * 1))
	testConnect(t)
	testLots(t)
}

func testConnect(t *testing.T) {
	c := client.NewClient()
	c.Connect()
	c.List()
	c.Stor(1024)
}

func testLots(t *testing.T) {
	for {
		c := client.NewClient()
		c.Connect()
		c.List()
		c.Stor(1024)
		time.Sleep(1 * (time.Second * 1))
	}
}