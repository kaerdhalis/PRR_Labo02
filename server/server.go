package main

import (
	"../config"
	"fmt"
)


type ipAdress struct {

	adress string
	port uint

}

type configuration struct {
	NumberOfProcesses int
	adresses []ipAdress
}


func main() {


	var conf config.Configuration
	conf = config.SetConfiguration()

	fmt.Println(conf.NumberOfProcesses)

}
