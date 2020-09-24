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

/*Node Lookup:
- Findkclosest
- take first alhpa contacts from k-closest
- send nodelookup request to the k-closest to the alpha nodes
- Check so no duplicates are added to the node, since the alpha nodes can have common
  k clostest nodes it can be duplicates
- End recursion when all alpha k clostest nodes have received response from their nodes
- Use SendFindContactMessage to add contact to the other nodes k buckets
- Use channels to define when a node is contacted

KademliaRandomId fungerar ej som den ska, vissa blir samma id, for fixa. 

*/

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
	fileKey := NewKademliaID("2111111400000000000000000000000000000000")
	clostestK := kademlia.routingTable.FindClosestContacts(fileKey, kademlia.k)
	fmt.Println("Närmsta k noderna:")
	fmt.Println(clostestK)
	for i:=0; i<len(clostestK); i++ {
		file := NewFile(*fileKey, data, clostestK[i])
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
