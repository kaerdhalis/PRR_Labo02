package mutex
import (
	"../network"
	"bytes"
	"encoding/gob"
	"math"
	"net"
	"time"
)


var n uint = 0
var id uint = 0
var h uint = 0
var pendingReq = false
var cs = false
var hReq uint= 0
var pDiff []uint
var pWait []uint


func run(request chan bool,wait chan bool,end chan bool) {

	var networkMsg chan network.Message
	for{

		select {
		
			case <- request:
			case <- wait:
			case <-end:
		case msg:=<-networkMsg:

			


		}



	}

}
 func request(){

 	pendingReq = true
 	h +=1
 	hReq = h
 	for  i := 0;i< len(pWait);i++{

		sendRequest()
	 }


 }

 func end(){
	h = h+1
	pendingReq= false
	cs = false

	pWait = pDiff
	for i:= 0;i< len(pDiff);i++{
		sendOk()
	}
	pDiff = nil

 }

 func requestTraitement(msg network.Message){


	h = max(msg.Hi,h)+1
	if pendingReq && (hReq< msg.Hi|| (hReq == msg.Hi && id < msg.Id )) {
		pDiff = append(pDiff,msg.Id)
	} else {
		sendOk()
	}

 }

 func okTraitement(ok network.Message){

	 h = max(ok.Hi,h)+1
 	var i =0
 	for pWait[i] != ok.Id{
 		i++
	 }
	 pWait = append(append(pWait[:i], pWait[i+1:]...))
 }

func max(x, y uint) uint {
	if x > y {
		return x
	}
	return y
}

func sendRequest(){

}

func sendOk(){

}