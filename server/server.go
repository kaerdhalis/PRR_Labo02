package main


type ipAdress struct {

	adress string
	port uint

}

type configuration struct {
	NumberOfProcesses int
	adresses []ipAdress
}
