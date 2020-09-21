package d7024e

import (
	"fmt"
	"time"

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
	kademlia.routingTable.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	go Listen("localhost:8001", 8001)
	<- time.After(1*time.Second)
	kademlia.routingTable.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	go Listen("localhost:8002", 8002)
	<- time.After(1*time.Second)
	/*kademlia.routingTable.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	kademlia.routingTable.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	kademlia.routingTable.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	kademlia.routingTable.AddContact(NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002"))
	*/
	contacts := kademlia.routingTable.FindClosestContacts(target.ID, kademlia.k)
	//kademlia.routingTable.AddContact(*target)
	//targetRoutingTable.AddContact(*kademlia.contact)
	fmt.Println(contacts)
	for i:=0; i<len(contacts); i++ {
		go kademlia.network.SendPingMessage(&contacts[i], port, donePing)
		targetRoutingTable.AddContact(contacts[i])
		fmt.Println("Ping message sent in LookUpContact")
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
