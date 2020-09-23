package main

import (
  "bufio"
  "fmt"
  "os"
  "strings"
  "d7024e"
)

func main() {

  reader := bufio.NewReader(os.Stdin)
  fmt.Println("Type operation below, you can choose between following: store, find, put, get, exit")
  fmt.Println("---------------------")

  for {
    fmt.Print("-> ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimRight(text, "\n")
	
	switch text {
		case "store":
			fmt.Println("enter store")
			fmt.Print("Enter data to store: ")
			data, _ := reader.ReadString('\n')
			data = strings.TrimRight(text, "\n")
			sendData := []byte(data + "\n")
			done := make(chan bool)
			bootstrapContact :=  d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
			rtBootstrap := d7024e.NewRoutingTable(bootstrapContact)
			network := d7024e.NewNetwork(bootstrapContact)
        	kademliaNetwork := d7024e.NewKademlia(&network, &bootstrapContact, rtBootstrap, 20, 3, done)
			d7024e.kademliaNetwork.Store(sendData)
		case "find":
			fmt.Println("enter find")
		case "put":
			fmt.Println("enter put")
		case "get":
			fmt.Println("enter get")
		case "exit":
			fmt.Println("enter exit")
		default:
			fmt.Println("Please type correct operation, you can choose between following: store, find, put, get, exit")
	}

  }

}