package config

import (
	"encoding/json"
	"log"
	"net"
	"os"
)

const localAdress = "127.0.0.1"
const defaultPort = 5555

var config configuration

type ipAdress struct {

	ip string
	port uint

}

type configuration struct {
	numberOfProcesses uint
	address []ipAdress
}

func GetAdressById(id uint) *net.TCPAddr{

	log.Println(config.numberOfProcesses)
	log.Println("totot")
	var localAdrr = new(net.TCPAddr)
	localAdrr.IP = net.ParseIP(config.address[id].ip)
	localAdrr.Port =int(config.address[id].port)

	return localAdrr
}

func SetConfiguration()  {

	file, err:= os.Open("src/config/config.json")

	if err != nil{
		log.Println(err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil{
		log.Println(err)
	}
}

func GetNumberOfProc() uint{

	return config.numberOfProcesses
}

