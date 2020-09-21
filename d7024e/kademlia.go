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

func (kademlia *Kademlia) LookupContact(target *Contact){
	// TODO
	kademlia.routingTable.AddContact(*target)
	contacts := kademlia.routingTable.FindClosestContacts(target.ID, 2)
	kademlia.done <- true
	fmt.Println(contacts)
	/*for i=0;i<len(contacts);i++{

	}*/

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
