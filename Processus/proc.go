package main



import (
"../network"
"bytes"
"encoding/gob"
"net"
"time"
)

func main() {

	var buf bytes.Buffer
	for {

		buf.Reset()

		adrress := new(net.TCPAddr)
		adrress.IP = net.ParseIP("127.0.0.1")
		adrress.Port = 8000

		err := gob.NewEncoder(&buf).Encode(network.Message{true,15,28		})
		if err != nil {
			// handle error
		}
		time.Sleep(2*time.Second) //simulates network delay of artificialNetworkDelay milliseconds
		network.ClientWriter(adrress, buf)
	}
}
