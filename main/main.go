package main

import (
	//"bufio"
	//"fmt"
	//"net"
	//"os"
	"kademlia"
)

func main() {
	address := "127.0.0.1:1234"
	id := kademlia.kademliaid.NewRandomKademliaID()
	node := kademlia.contact.NewContact(id, address)

	kademlia.network.SendPingMessage(node)
}