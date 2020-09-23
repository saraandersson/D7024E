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

/*Network joining and node lookup*/
func (kademlia *Kademlia) LookupContact(target *Contact, targetRoutingTable *RoutingTable, port int){
	// TODO
	donePing := make(chan bool)
	contacts := kademlia.routingTable.FindClosestContacts(target.ID, kademlia.k)
	addAlphaContacts := make([]Contact, 0)
	

	/*Picks alpha first nodes from the k closest*/
	if (len(contacts)>kademlia.alpha) {
		addAlphaContacts = append(addAlphaContacts, contacts[0:kademlia.alpha]...)
	} else {
		addAlphaContacts = append(addAlphaContacts, contacts...)
	}
	/*Network joining*/
	kademlia.routingTable.AddContact(*target)
	targetRoutingTable.AddContact(*kademlia.contact)
	/*Start for node lookup*/
	fmt.Println(addAlphaContacts)
	for i:=0; i<len(addAlphaContacts); i++ {
		go kademlia.network.SendPingMessage(&addAlphaContacts[i], port, donePing)
		targetRoutingTable.AddContact(addAlphaContacts[i])
		fmt.Println("Ping sent and received response in LookUpContact")
		<- donePing
	}
	//kademlia.done <- true
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	fmt.Println("Enter store")
	fmt.Print(data)
	fileKey := NewRandomKademliaID()
	clostestK := kademlia.routingTable.FindClosestContacts(fileKey, kademlia.k)
	fmt.Println("Närmsta k noderna:")
	fmt.Println(clostestK)
	for i:=0; i<len(clostestK); i++ {
		file := NewFile(fileKey, data, &clostestK[i])
		kademlia.network.StoreDataOnNode(file)
	}


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
