package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type client chan<- Message // an outgoing message channel
var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan Message) // all incoming client messages
)

type Message struct{

	MsgType bool
	Id uint
	Hi uint
}

type SharedValueMessage struct {
	SharedValue int64
}


func ClientWriter(address *net.TCPAddr,buf bytes.Buffer) {

	conn, err := net.DialTCP("tcp",nil,address)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {

		_, err  = conn.Write(buf.Bytes())// NOTE: ignoring errors
		//log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	if _, err := conn.Write(buf.Bytes()); err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	<-done // wait for background goroutine to finish
}

func ClientReader(address *net.TCPAddr,message chan Message,sharedValue chan SharedValueMessage) {
	// error testing suppressed to compact listing on slides

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn,message,sharedValue)
	}
}
func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages: // broadcaster <- handleConn
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg // clientwriter (handleConn) <- broadcaster
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func handleConn(conn *net.TCPConn,message chan Message,value chan SharedValueMessage) {
	ch := make(chan Message) // channel 'client' mais utilisé ici dans les 2 sens
	go func() {             // clientwriter
		for msg := range ch { // clientwriter <- broadcaster, handleConn
			fmt.Fprintln(os.Stdout, msg.Id)
			message <-msg// netcat Client <- clientwriter

		}
	}()

	//who := conn.RemoteAddr().String()
	//ch <- "You are " + who           // clientwriter <- handleConn
	//messages <- who + " has arrived" // broadcaster <- handleConn
	entering <- ch
	buf := make([]byte, 1024)


	n,_ := conn.Read(buf) // n,addr, err := p.ReadFrom(buf)


	var msg Message
	var val SharedValueMessage
	if err := gob.NewDecoder(bytes.NewReader(buf[:n])).Decode(&msg); err == nil {
		message<-msg
	}else if err := gob.NewDecoder(bytes.NewReader(buf[:n])).Decode(&val); err == nil{

		value <- val
	}
	leaving <- ch
	//messages <- who + " has left" // broadcaster <- handleConn
	conn.Close()
}


func decrypt(conn *net.TCPConn,message chan Message,value chan SharedValueMessage) {

	buf := make([]byte, 1024)


	n,_ := conn.Read(buf) // n,addr, err := p.ReadFrom(buf)


	var msg Message
	var val SharedValueMessage
	if err := gob.NewDecoder(bytes.NewReader(buf[:n])).Decode(&msg); err == nil {
		message<-msg
	}else if err := gob.NewDecoder(bytes.NewReader(buf[:n])).Decode(&val); err == nil{
		
		value <- val
	}
	

}

func PingAdress(address *net.TCPAddr,id uint) {

	timeout := time.Duration(1 * time.Second)
	for {

		conn, err := net.DialTimeout("tcp", address.String(), timeout)
		if err != nil {
			//log.Println("Site unreachable, error: ", err)

		} else {

			fmt.Printf("Processus %d is Up and Ready\n",id)
			conn.Close()
			break


		}
	}
}
