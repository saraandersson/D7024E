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
	address := "127.0.0.1:1234"
	id := d7024e.NewRandomKademliaID()
	node := d7024e.NewContact(id, address)
	fmt.Println(id)
	fmt.Println(node)

	//test := d7024e.SendPingMessage(node)
}