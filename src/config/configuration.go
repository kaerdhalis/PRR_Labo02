package config

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

const localAdress = "127.0.0.1"
const defaultPort = 5555

var config configuration

type IpAdress struct {

	Ip string
	Port uint

}

type configuration struct {
	NumberOfProcesses uint
	Address []IpAdress
}

func GetAdressById(id uint) *net.TCPAddr{

	var localAdrr = new(net.TCPAddr)
	localAdrr.IP = net.ParseIP(config.Address[id].Ip)
	localAdrr.Port =int(config.Address[id].Port)

	return localAdrr
}

func SetConfiguration()  {


	file, err:= os.Open("src/config/config.json")

	if err != nil{
		fmt.Println(err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil{
		fmt.Println(err)
	}

}

func GetNumberOfProc() uint{

	//return config.NumberOfProcesses
	return 2
}

