package d7024e

import (
	"fmt"
)

type Network struct {
}

func Listen(ip string, port int) {
	// TODO
}

func (network *Network) SendPingMessage(contact *Contact) {
	fmt.Println("Kommer till ping!")
	// TODO
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func (network *Network) createNetwork() {
	network := Network{}
	return network
}
