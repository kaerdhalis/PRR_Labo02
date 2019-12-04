/**
 * Title: 			Labo2 - Mutual exclusion
 * File:			client.go
 * Date:			20.11.12
 * Authors:			Le Guillou Benjamin, Reis de Carvalho Luca
 *
 * Description:		File containing the client side of the process. It can read the inputs of the users and read or
 *                  modify the shared value.
 */
package main

import (
	"../config"
	"../mutex"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

//global value shared by all the processes
var sharedValue  int64 = 0

func main() {

	//read the id of the process passed in argument
	args := os.Args[1:]

	if len(args)!=1{
		log.Fatal("Number of arguments invalid, you need to pass the id of the Process")
	}
	id,_ := strconv.Atoi(args[0])

	//channel to communicate to the mutex
	request:= make(chan bool)
	wait := make(chan bool)
	end := make(chan int64)
	valueChannel := make(chan int64)

	//get the global configuration of the application
	config.SetConfiguration()


	go changeSharedValue(valueChannel)

	//launch the mutex
	go mutex.Run(request,wait,end,valueChannel,uint(id))

	//wait for the server to be launched and the other processes to be up to start reading inputs
	<-wait

	for{
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter (r) to Read or (w <value>) to change sharedValue, (d <time>) to change the artifical delay:")
		fmt.Print(">")

		//read the imput and ignore UperCase
		value, _ := reader.ReadString('\n')
		value = strings.ToLower(value[:len(value)-1])
		tokens := strings.Split(value, " ")

		switch tokens[0] {

		case "r":

			fmt.Printf("sharedValue = %d\n",sharedValue)
			break

		case "d":

			delay, err := strconv.ParseFloat(tokens[1], 64)

			if err != nil || delay <0{

				fmt.Println("Wrong Argument: delay must be a positive float \n")
				break
			}

			config.SetTransmitdelay(delay)

			fmt.Printf("\n new delay is set to %f seconds\n",config.GetTransmitDelay())

			break

		case "w":

			newValue, err := strconv.ParseInt(tokens[1], 10, 64)

			if err != nil{

				fmt.Println("Wrong Argument: new value must be an integer \n")
				break
			}

			fmt.Println("\n***requesting critical section***\n")
			//send a request to the mutex
			request<-true

			//wait for the critical section
			<-wait
			fmt.Println("***entering crtical section***")
			fmt.Printf("\nprevious sharedValue = %d\n",sharedValue)

			sharedValue = int64(newValue)
			fmt.Printf("\nnew value is set to %d\n\n",sharedValue)

			fmt.Println("\n***Doing others artificial stuffs in critical section***\n")
			time.Sleep(time.Duration(config.GetArtificialDelay()) * time.Second)


			fmt.Println("***leaving critical section***\n")
			//send the new sharedValue to the mutex and leave critical section
			end <-sharedValue
		}
	}
}

//Function used to update the shared value if modified by an other process
func changeSharedValue(valueChannel chan int64 ) {

	for {

		value:= <-valueChannel
		sharedValue = value
	}
}

