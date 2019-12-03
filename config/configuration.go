package config

import (
	"encoding/json"
	"net"
	"os"
)

const localAdress = "127.0.0.1"
const defaultPort = 5555

var adresses[] string
var ports[] uint

var numberOfProc uint
var config Configuration

type ipAdress struct {

	adress string
	port uint

}

type Configuration struct {
	NumberOfProcesses int
	Adresses []ipAdress
}

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

func SetConfiguration() Configuration {

	file, _ := os.Open("src/config/config.json")

	decoder := json.NewDecoder(file)
	_ = decoder.Decode(&config)
	return  config
}

func GetNumberOfProc() uint{

	//return numberOfProc
	return 4
}

