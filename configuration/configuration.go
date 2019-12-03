package configuration

import (
	"net"
)

const localAdress = "localhost"
const defaultPort = 5555

var adresses[] string
var ports[] uint

var numberOfProc uint

func GetAdressById(id uint) *net.TCPAddr{

	var adress string
	var port uint
	if adresses!=nil && ports != nil {

		adress = adresses[id]
		port = ports[id]
	} else {

		adress = "127.0.0.1"
		port = defaultPort + id

	}

	var localAdrr = new(net.TCPAddr)
	localAdrr.IP = net.ParseIP(adress)
	localAdrr.Port =int(port)

	return localAdrr
}

func GetNumberOfProc() uint{

	//return numberOfProc
	return 2
}

