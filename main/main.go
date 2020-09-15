package main

import (
	"bufio"
	"fmt"
	//"net"
	"os"
	//"d7024e"
)

import "d7024e"

func main() {
	var contact d7024e.Contact
	address := "127.0.0.1:1234"
	id := d7024e.NewRandomKademliaID()
	contact = d7024e.NewContact(id, address)
	fmt.Println(id)
	fmt.Println(contact)
	network := d7024e.CreateNetwork(&contact)
	reader := bufio.NewReader(os.Stdin)
    fmt.Print("Type operation here: ")
	text, _ := reader.ReadString('\n')
	if text == "ping" {
		network.SendPingMessage(&contact)
	}
	if text == "join network" {

	}
	if text == "node lookup" {
		fmt.Print("Enter targetNode id: ")
		targetNode, _ := reader.ReadString('\n')
		var targetNodeId d7024e.KademliaID
		targetNodeId = targetNode
		contacts := d7024e.FindClosestContacts(targetNodeId, 1)
		fmt.Println(contacts)
	}
}