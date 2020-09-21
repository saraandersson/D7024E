package d7024e

import (
	"fmt"

)

type Kademlia struct {
	network 		*Network
	contact 		*Contact
	routingTable 	*RoutingTable
	k 				int
	alpha			int
	done 			chan bool
}

func (kademlia *Kademlia) LookupContact(target *Contact, targetRoutingTable *RoutingTable){
	// TODO
	donePing := make(chan bool)
	kademlia.routingTable.AddContact(*target)
	contacts := kademlia.routingTable.FindClosestContacts(target.ID, 2)
	fmt.Println(contacts)
	for i:=0; i<len(contacts); i++ {
		fmt.Println(contacts[i])
		go kademlia.network.SendPingMessage(&contacts[i], donePing)
		<- donePing
		targetRoutingTable.AddContact(contacts[i])
	}
	kademlia.done <- true

}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

func NewKademlia(network *Network, contact *Contact, routingTable *RoutingTable, k int, alpha int, done chan bool) Kademlia{
	kademlia := Kademlia{}
	kademlia.network=network
	kademlia.contact = contact
	kademlia.routingTable = routingTable
	kademlia.k = k
	kademlia.alpha = alpha
	kademlia.done = done
	return kademlia
}
