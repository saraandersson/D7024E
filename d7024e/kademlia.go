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
	contacts := kademlia.routingTable.FindClosestContacts(target.ID, kademlia.k)
	kademlia.routingTable.AddContact(*target)
	targetRoutingTable.AddContact(*kademlia.contact)
	fmt.Println(contacts)
	for i:=0; i<len(contacts); i++ {
		go kademlia.network.SendPingMessage(&contacts[i], port, donePing)
		targetRoutingTable.AddContact(contacts[i])
		fmt.Println("Ping sent and received response in LookUpContact")
		<- donePing
	}
	//kademlia.done <- true

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
