/**
 * Title: 			Labo2 - Mutual exclusion
 * File:			client.go
 * Date:			20.11.12
 * Authors:			Le Guillou Benjamin, Reis de Carvalho Luca
 *
 * Description:		File containing the configurationof the application. It contains the adress of every process and other
 *                  values like the waiting time in the critical section
 */
package config

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

// contains all the values of configuration
var config configuration

//struct representing the adress from the json file
type IpAdress struct {

	Ip string
	Port uint

}

// Configuration struct represents all configurations from the json file
type configuration struct {
	NumberOfProcesses uint
	Address []IpAdress
	ArtificialDelay uint
}


func GetAdressById(id uint) *net.TCPAddr{

	var localAdrr = new(net.TCPAddr)

	localAdrr.IP = net.ParseIP(config.Address[id].Ip)
	localAdrr.Port =int(config.Address[id].Port)

	return localAdrr
}

//read the json file and stock the values
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

	return config.NumberOfProcesses
}

func GetArtificialDelay() uint{

	return config.ArtificialDelay
}

