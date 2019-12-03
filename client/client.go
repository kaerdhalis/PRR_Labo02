package main

import (
	"../mutex"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var sharedValue  int64 = 0

func main() {

		request:= make(chan bool)
		wait := make(chan bool)
		end := make(chan int64)
		valcnannel := make(chan int64)

	go changeSharedValue(valcnannel)
	go mutex.Run(request,wait,end,valcnannel,0)

	<-wait

	for{

		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter (r) to Read or (w <value>) to change sharedValue:")
		fmt.Print(">")

		value, _ := reader.ReadString('\n')
		value = strings.ToLower(value[:len(value)-1])

		tokens := strings.Split(value, " ")

		switch tokens[0] {

		case "r":
			fmt.Printf("sharedValue = %d\n",sharedValue)
			break

		case "w":
			newValue, err := strconv.ParseInt(tokens[1], 10, 64)

			if err != nil{

				fmt.Println("new value must be an integer \n")
				break
			}

			fmt.Println("requesting critical section")
			request<-true

			<-wait
			fmt.Println("entering crtical section")
			fmt.Printf("previous sharedValue = %d\n",sharedValue)

			sharedValue = int64(newValue)
			fmt.Printf("new value is set to %d\n",sharedValue)

			end <-sharedValue
		}


	}



}

func changeSharedValue(valchannel chan int64 ) {
	for {

		value:= <-valchannel
		sharedValue = value
	}
}

