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

func (kademlia *Kademlia) LookupContact(target *Contact, targetRoutingTable *RoutingTable, port int){
	// TODO
	donePing := make(chan bool)
	doneAdd1 := make(chan bool)
	doneAdd2 := make(chan bool)
	contacts := kademlia.routingTable.FindClosestContacts(target.ID, kademlia.k)
	kademlia.routingTable.AddContact(*target, doneAdd1)
	<- doneAdd1
	targetRoutingTable.AddContact(*kademlia.contact, doneAdd2)
	fmt.Println(contacts)
	<- doneAdd2
	for i:=0; i<len(contacts); i++ {
		go kademlia.network.SendPingMessage(&contacts[i], port, donePing)
		//targetRoutingTable.AddContact(contacts[i],)
		fmt.Println("Sara")
		<- donePing
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
