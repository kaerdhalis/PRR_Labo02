/**
 * Title: 			Labo2 - Mutual exclusion
 * File:			client.go
 * Date:			20.11.12
 * Authors:			Le Guillou Benjamin, Reis de Carvalho Luca
 *
 * Description:		File containing the mutex side of the process. It implements the Carvalho & Roucairol algorithm and
 * 					manages the interaction with the network side of the process
 */
package mutex

import (
	"../config"
	"../network"
	"bytes"
	"encoding/gob"
	"fmt"
)

const(
	REQ = false
	OK = true
)
var(
	 id uint = 0 //Id of the process
	 h uint =  0 //timestamp
	 pendingReq = false //used to check if process is requesting critical section
	 cs = false //used to know if the process is in the critical section
	 hReq uint= 0 //timestamp of the request to access the critical section
	 pDiff []uint //array containing the differed processes
	 pWait []uint //array containg the processes on which it wait an OK
)

func Run(request chan bool,wait chan bool,end chan int64,valchannel chan int64,processId uint) {

	//get the id of the process from the client
	id = processId

	//get the adress of the process from the configuration
	var localAdrr = config.GetAdressById(id)

	//channel to communicate with the network
	networkMsg := make(chan network.Message)
	sharedVal := make(chan network.SharedValueMessage)

	//put all the processes in the wait list
	for i:=0 ; i<int(config.GetNumberOfProc()) ; i++ {

		if uint(i) != id {

			pWait = append(pWait, uint(i))
		}
	}


	//launch the server to listen to connections
	go network.ClientReader(localAdrr,networkMsg,sharedVal)

	//check that all others processes are ready
	checkAllProcessAreReady()

	//signal the client it can begin read inputs
	wait<-true

	//main loop of the algorithm
	for{

		select {
		
			case <- request:

				requestHandle()

			case newValue:=<-end:

				transmitSharedValue(newValue)
				endHandle()

			case msg:=<-networkMsg:

				//check the type of message
				if msg.MsgType ==REQ {
					requestTraitement(msg)
				}else  {
					okTraitement(msg,wait)
				}

			case val :=<- sharedVal:

				//transmit the new value from other process to the client
				valchannel <- val.SharedValue
		}
	}
}

//handle the request of CS from the client
 func requestHandle(){

 	h += 1
 	pendingReq = true
 	hReq = h
 	for  i := 0;i< len(pWait);i++{

 		if pWait[i] != id {

			sendMessage(hReq,pWait[i],REQ)
		}
	 }
 }

 //handle the exit of the critical section by the client
 func endHandle(){

	h = h+1
	cs = false
	pendingReq= false
	for i:= 0;i< len(pDiff);i++ {

		if pDiff[i] != id {

			sendMessage(h, pDiff[i], OK)
		}
	}

	pWait = pDiff
	pDiff = nil

 }

 //handle the incomings requests from other processes
 func requestTraitement(rqst network.Message){

	h = max(rqst.Hi,h)+1

	if pendingReq==false{

		sendMessage(rqst.Hi,rqst.Id,OK)
		pWait = append(pWait,rqst.Id)

	}else if cs || (hReq< rqst.Hi)|| (hReq == rqst.Hi && id < rqst.Id ) {

		pDiff = append(pDiff,rqst.Id)

	} else {

		sendMessage(h,rqst.Id,OK)
		pWait = append(pWait,rqst.Id)
		sendMessage(hReq,rqst.Id,REQ)
	}
 }

 //handle the incomings OK from the other processes
 func okTraitement(ok network.Message, wait chan bool){

 	h = max(ok.Hi,h)+1

 	for i :=0;i< len(pWait);i++{

 		if pWait[i] == ok.Id{

			pWait = append(append(pWait[:i], pWait[i+1:]...))
		}
	}

	//if the wait list is empty signal the client that it can enter to critical section
	 if len(pWait)==0 {
	 	cs = true
	 	wait<-true
	 }
 }


 func sendMessage(hi uint,procesId uint,messageType bool){

	var buf bytes.Buffer
	var adress = config.GetAdressById(procesId)

	if err := gob.NewEncoder(&buf).Encode(network.Message{MsgType:messageType,Id:id,Hi:hi}); err != nil {
		fmt.Println(err)
	}

	network.ClientWriter(adress,buf)
}


//util function to get max of two int values
 func max(x, y uint) uint {
	if x > y {
		return x
	}
	return y
}

 //function which check if all the other processes are ready
 func checkAllProcessAreReady(){

	for i:=0 ;i<int(config.GetNumberOfProc()) ;i++ {

		if uint(i) != id {
		network.PingAdress(config.GetAdressById(uint(i)), uint(i))
		}
	}

	fmt.Println("All Process are Ready")
}

//transmit the new value to all the other processes
func transmitSharedValue(value int64){


	var buf bytes.Buffer

	for i:=0 ;i<int(config.GetNumberOfProc()) ;i++  {
		buf.Reset()

		var adress = config.GetAdressById(uint(i))

		if err := gob.NewEncoder(&buf).Encode(network.SharedValueMessage{value}); err != nil {

			fmt.Println(err)
		}

		network.ClientWriter(adress,buf)
	}

}
