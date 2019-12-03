package mutex

import (
	"../network"
	"../config"
	"bytes"
	"encoding/gob"
	"fmt"
)

const(
	REQ = false
	OK = true
)
var(
	 n uint = 0
	 id uint = 0
	 h uint = 0
	 pendingReq = false
	 cs = false
	 hReq uint= 0
	 pDiff []uint
	 pWait []uint
)

func Run(request chan bool,wait chan bool,end chan int64,valchannel chan int64,processId uint) {

	id = processId
	var localAdrr = config.GetAdressById(id)

	networkMsg := make(chan network.Message)
	sharedVal := make(chan network.SharedValueMessage)

	for i:=0 ;i<int(config.GetNumberOfProc()) ;i++ {

		if uint(i) != id {

		pWait = append(pWait, uint(i))
		}
	}


	go network.ClientReader(localAdrr,networkMsg,sharedVal)

	checkAllProcessAreReady()

	wait<-true

	for{

		select {
		
			case <- request:
				requestHandle()

			case newValue:=<-end:

				transmitSharedValue(newValue)
				endHandle()

			case msg:=<-networkMsg:
				if msg.MsgType ==REQ {
					requestTraitement(msg)
				}else  {
					okTraitement(msg,wait)
				}
			case val :=<- sharedVal:
				valchannel <- val.SharedValue


		}



	}

}
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

 func okTraitement(ok network.Message, wait chan bool){



	 h = max(ok.Hi,h)+1

 	for i :=0;i< len(pWait);i++{

 		if pWait[i] == ok.Id{

			pWait = append(append(pWait[:i], pWait[i+1:]...))
		}
	}

	 if len(pWait)==0 {
	 	cs = true
	 	wait<-true
	 }
 }


func sendMessage(hi uint,procesId uint,messageType bool){

	var buf bytes.Buffer
	var adress = config.GetAdressById(procesId)

	if err := gob.NewEncoder(&buf).Encode(network.Message{MsgType:messageType,Id:id,Hi:hi}); err != nil {
		// handle error
	}

	network.ClientWriter(adress,buf)
}



func max(x, y uint) uint {
	if x > y {
		return x
	}
	return y
}

func checkAllProcessAreReady(){

	for i:=0 ;i<int(config.GetNumberOfProc()) ;i++ {

		if uint(i) != id {
		network.PingAdress(config.GetAdressById(uint(i)), uint(i))
	}
	}

	fmt.Println("All Process are Ready")
}

func transmitSharedValue(value int64){


	var buf bytes.Buffer

	for i:=0 ;i<int(config.GetNumberOfProc()) ;i++  {
		buf.Reset()

		var adress = config.GetAdressById(uint(i))

		if err := gob.NewEncoder(&buf).Encode(network.SharedValueMessage{value}); err != nil {
			// handle error
		}
		network.ClientWriter(adress,buf)

	}

}
