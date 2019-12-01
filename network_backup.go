package network

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

type client chan<- string // an outgoing message channel
var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

type Message struct{

	Msg string
}

func ClientWriter(address *net.TCPAddr,buf bytes.Buffer) {

	conn, err := net.DialTCP("tcp",nil,address)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		_, err  = buf.WriteTo(conn) // NOTE: ignoring errors
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()

	conn.Close()
	<-done // wait for background goroutine to finish
}

func ClientReader(address *net.TCPAddr) {
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
		go handleConn(conn)
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

func handleConn(conn *net.TCPConn) {
ch := make(chan string) // channel 'client' mais utilisé ici dans les 2 sens
go func() {             // clientwriter
for msg := range ch { // clientwriter <- broadcaster, handleConn
fmt.Println( msg) // netcat Client <- clientwriter
}
}()

who := conn.RemoteAddr().String()
ch <- "You are " + who           // clientwriter <- handleConn
//messages <- who + " has arrived" // broadcaster <- handleConn
entering <- ch


messages <- who + ": " + decrypt(conn).Msg // broadcaster <- handleConn

leaving <- ch
messages <- who + " has left" // broadcaster <- handleConn
conn.Close()
}

func decrypt(conn *net.TCPConn) Message {

	buf := make([]byte, 1024)


	input := bufio.NewScanner(conn)

	for input.Scan() { // handleConn <- netcat client
		buf = input.Bytes()
	}
	buf = input.Bytes()
	//n,err := conn.Read(buf) // n,addr, err := p.ReadFrom(buf)


	//if err != nil {
	//	log.Fatal(err)
	//}

	var msg Message
	if err := gob.NewDecoder(bytes.NewReader(buf)).Decode(&msg); err != nil {
		// handle error
	}

	return msg

}
