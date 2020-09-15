package main

import (
	//"bufio"
	//"fmt"
	//"net"
	//"os"
	"d7024e"
)

func main() {
	address := "127.0.0.1:1234"
	id := d7024e.kademliaid.NewRandomKademliaID()
	node := d7024e.contact.NewContact(id, address)

	d7024e.network.SendPingMessage(node)
}