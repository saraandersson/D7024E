package d7024e

import (
	"fmt"
	"strings"
	"strconv"
	"time"

)

type Kademlia struct {
	network 		*Network
	contact 		*Contact
	routingTable 	*RoutingTable
	k 				int
	alpha			int
}

/*TODO: Fixa findValue (likt NodeLookUp) och ändra ping så den skickar fillistan tillbaka som ett meddelande*/


/*Network joining and node lookup*/
func (kademlia *Kademlia) LookupContact(kademliaId KademliaID) []Contact{
	contacts := kademlia.network.routingTable.FindClosestContacts(&kademliaId, kademlia.k)
	addAlphaContacts := make([]Contact, 0)
	

	/*Picks alpha first nodes from the k closest*/
	if (len(contacts)>kademlia.alpha) {
		addAlphaContacts = append(addAlphaContacts, contacts[0:kademlia.alpha]...)
	} else {
		addAlphaContacts = append(addAlphaContacts, contacts...)
	}

	/*Node lookup*/
	contactsToAdd := kademlia.NodeLookUp(addAlphaContacts, addAlphaContacts, *kademlia.network.contact)

	for i:=0; i<len(contactsToAdd); i++ {
		kademlia.network.routingTable.AddContact(contactsToAdd[i])
	}
	return contactsToAdd
}

/*NodeLookUp function*/
func (kademlia *Kademlia) NodeLookUp(alphaContacts []Contact, addedContacts []Contact, currentContact Contact) []Contact {
	/*Base case*/
	if (len(alphaContacts)==0) {
		return addedContacts
	}
	/*Ping alpha node to see if it is alive*/
	pingAlphaNode := make(chan bool)
	stringListAlpha := strings.Split(alphaContacts[0].Address, ":")
	portStringAlpha := stringListAlpha[1]
	portAlpha, _ := strconv.Atoi(portStringAlpha)
	go SendPingMessage(&currentContact,&alphaContacts[0], portAlpha, pingAlphaNode)
	select {
	case <-pingAlphaNode:
		/*If alpha-node is alive, fetch k clostest nodes to the alpha node*/
		returnMessage := make(chan []Contact)
		go SendFindAlphaMessage(&currentContact, &alphaContacts[0], portAlpha, returnMessage)
		select {
		case kContactsFromAlpha := <- returnMessage:
			/*kContactsFromAplha = channelvärdet från message*/
			var addAlphaContacts []Contact
			/*Picks alpha first nodes from the k closest*/
			if (len(kContactsFromAlpha)>kademlia.alpha) {
				addAlphaContacts = append(addAlphaContacts, kContactsFromAlpha[0:kademlia.alpha]...)
			} else {
				addAlphaContacts = append(addAlphaContacts, kContactsFromAlpha...)
			}
			kContactsFromAlpha = addAlphaContacts
			/*Loop through all k nodes*/
			for i:=0; i<len(kContactsFromAlpha); i++ {
				isInList := false
				if (len(addedContacts) > 0) {
					/*Loop through all nodes in addedContacts*/
					for x:=0; x<len(addedContacts); x++ {
						/*Check that a node is not already in addedContacts*/
						if (kContactsFromAlpha[i].ID.String() == addedContacts[x].ID.String()){
							isInList = true
						}
					}
				}
				/*If the node is not in list and the node is not currentNode, ping the node and if it is alive,
					add node to addedContacts
				*/
				if ((isInList==false) && (kContactsFromAlpha[i].ID.String() != currentContact.ID.String())) {
					donePing := make(chan bool)
					stringList := strings.Split(kContactsFromAlpha[i].Address, ":")
					portString := stringList[1]
					port, _ := strconv.Atoi(portString)
					go SendPingMessage(&currentContact,&kContactsFromAlpha[i], port, donePing)
					select {
					case <- donePing:
						addedContacts = append(addedContacts,kContactsFromAlpha[i])
					/*If node does not answer in 10 seconds, the node is dead*/
					case <-time.After(10*time.Second):
						/*Send ping again??*/
						fmt.Println("TIMEOUT")
						continue
					}
				}
			}

	}
	/*If node does not answer in 10 seconds, remove alpha node from list*/
	case <-time.After(10*time.Second):
		addedContacts = RemoveFromList(addedContacts, alphaContacts[0])
	}
	/*Remove first element from list and continue recursion*/
	if (len(alphaContacts) > 1) {
		alphaContacts = alphaContacts[1:]
	} else {
		alphaContacts = make([]Contact, 0)
	}
	return kademlia.NodeLookUp(alphaContacts, addedContacts, currentContact)

}

/*Help functinon, remove a contact from a list of contacts */
func RemoveFromList(contacts []Contact, contactToRemove Contact) []Contact {
	fmt.Println("Node is dead, removing node")
	for i := 0; i < len(contacts); i++ {
		contact := contacts[i]
		if contact == contactToRemove {
			contacts = append(contacts[:i], contacts[i+1:]...)
			i--
		}
	}
	fmt.Println("Remove done")
	return contacts
}


func (kademlia *Kademlia) LookupData(hash string) []byte{
	fmt.Println("Enter LookupData")
	dataKey := NewKademliaID(hash)
	data := kademlia.network.FindData(*dataKey)
	if (data!=nil){
		fmt.Println("Data found on node!")
		fmt.Println(data)
		return data
	} else{
		contacts := kademlia.network.routingTable.FindClosestContacts(&dataKey, kademlia.k)
		addAlphaContacts := make([]Contact, 0)
		/*Picks alpha first nodes from the k closest*/
		if (len(contacts)>kademlia.alpha) {
			addAlphaContacts = append(addAlphaContacts, contacts[0:kademlia.alpha]...)
		} else {
			addAlphaContacts = append(addAlphaContacts, contacts...)
		}
		LookupKClosestData(dataKey, addAlphaContacts)
	}
}

func (kademlia *Kademlia) LookupKClosestData(dataKey KademliaID, addAlphaContacts []Contact) []byte{



}

func (kademlia *Kademlia) Store(data []byte) {
	fmt.Println("Enter store")
	fmt.Print(data)
	fileKey := NewRandomKademliaID()
	fmt.Println("Filekey: " + fileKey.String())
	clostestK := kademlia.LookupContact(*fileKey)
	fmt.Println("Närmsta k noderna:")
	fmt.Println(clostestK)
	for i:=0; i<len(clostestK); i++ {
		//file := NewFile(*clostestK[i].ID, data, clostestK[i])
		//kademlia.network.StoreDataOnNode(file)
		doneStorePing := make(chan bool)
		go SendPingStoreMessage(&clostestK[i], data, 8080, doneStorePing)
		select {
		case <-doneStorePing:
			fmt.Println("File stored on node: " +clostestK[i].ID.String())
		case <-time.After(10*time.Second):
			fmt.Println("TIMEOUT IN STORE")
		}
	}


	// TODO
}

func NewKademlia(network *Network, contact *Contact, k int, alpha int) Kademlia{
	kademlia := Kademlia{}
	kademlia.network=network
	kademlia.contact = contact
	kademlia.routingTable = network.routingTable
	kademlia.k = k
	kademlia.alpha = alpha
	return kademlia
}
