package mutex

import (
	"../network"
	"../configuration"
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

func run(request chan bool,wait chan bool,end chan bool,processId uint) {

	id = processId
	var localAdrr = configuration.GetAdressById(id)
	var networkMsg chan network.Message

	for i:=0 ;i<int(configuration.GetNumberOfProc()) ;i++  {

		pWait= append(pWait,uint(i))
	}

	//faire un ping avec dialTimeout

	go network.ClientReader(localAdrr,networkMsg)

	checkAllProcessAreReady()

	wait<-true

	for{

		select {
		
			case <- request:
				fmt.Println("demande acces a la section crittique")
				requestHandle()

			case <-end:
				fmt.Println("sortie de la section critique")
				endHandle()

			case msg:=<-networkMsg:
				if msg.MsgType ==REQ {
					requestTraitement(msg)
				}else  {
					okTraitement(msg,wait)
				}

			


		}



	}

}
 func requestHandle(){

 	pendingReq = true
 	h +=1
 	hReq = h
 	for  i := 0;i< len(pWait);i++{

		sendMessage(hReq,pWait[i],REQ)
	 }

 }

 func endHandle(){

	h = h+1
	pendingReq= false
	cs = false

	pWait = pDiff
	for i:= 0;i< len(pDiff);i++{
		sendMessage(h,pDiff[i],OK)
	}
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
		sendMessage(hReq,id,REQ)
	}

 }

 func okTraitement(ok network.Message, wait chan bool){

 	h = max(ok.Hi,h)+1
 	var i =0
 	for pWait[i] != ok.Id{
 		i++
	 }
	 pWait = append(append(pWait[:i], pWait[i+1:]...))

	 if len(pWait)==0 {
	 	cs = true
	 	wait<-true
	 }
 }


func sendMessage(hi uint,procesId uint,messageType bool){

	var buf bytes.Buffer
	var adress = configuration.GetAdressById(procesId)

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

	for i:=0 ;i<int(configuration.GetNumberOfProc()) ;i++  {

		network.PingAdress(configuration.GetAdressById(uint(i)),uint(i))
	}

	fmt.Println("All Process are Ready")
}
