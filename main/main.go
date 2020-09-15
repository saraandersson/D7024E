package main

import (
	//"bufio"
	"fmt"
	//"net"
	//"os"
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
	network.SendPingMessage(&contact)
}