package main

import (
	"../network"
	"net"
)

func main() {

	adrress := new(net.TCPAddr)
	adrress.IP = net.ParseIP("127.0.0.1")
	adrress.Port = 8000

	network.ClientReader(adrress)
	
}
