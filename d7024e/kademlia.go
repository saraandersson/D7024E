package d7024e

import (
	"fmt"
	//"strings"
	//"strconv"
	"time"
	"math/rand"

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
	fmt.Println("kommer till lookUpContact")
	contacts := kademlia.network.routingTable.FindClosestContacts(&kademliaId, kademlia.k)
	addAlphaContacts := make([]Contact, 0)
	

	/*Picks alpha first nodes from the k closest*/
	if (len(contacts)>kademlia.alpha) {
		addAlphaContacts = append(addAlphaContacts, contacts[0:kademlia.alpha]...)
	} else {
		addAlphaContacts = append(addAlphaContacts, contacts...)
	}

	/*Node lookup*/
	contactedContacts := make([]Contact, 0)
	contactsToAdd := kademlia.NodeLookUp(addAlphaContacts, contactedContacts)
	//contactsToAdd = sortList(contactsToAdd)

	for i:=len(contactsToAdd)-1; i>0; i-- {
		kademlia.network.UpdateKBucket(contact);
		//kademlia.network.routingTable.AddContact(contactsToAdd[i])
	}
	return contactsToAdd
}
/*ShortList = Nodes to contact*/
func (kademlia *Kademlia) NodeLookUp(shortList []Contact, contactedContacts []Contact) []Contact {
	if (len(shortList)==0){
		contactedContacts = kademlia.sortList(contactedContacts)
		return contactedContacts
	}
	nextRound := make([]Contact, 0)
	shortList, contactedContacts = kademlia.NodeLookUpRound(shortList, contactedContacts,nextRound)
	return kademlia.NodeLookUp(shortList, contactedContacts)
}

/*NodeLookUp function*/
func (kademlia *Kademlia) NodeLookUpRound(shortList []Contact, contactedContacts []Contact, nextRound []Contact) ([]Contact, []Contact) {
	/*Base case*/
	if (len(shortList)==0) {
		nextRound = kademlia.sortList(nextRound)
		if (len(nextRound)>kademlia.alpha) {
			nextRound = append(nextRound, nextRound[0:kademlia.alpha]...)
		}
		return nextRound, contactedContacts
	}
	returnMessage := make(chan []Contact)
	go SendFindNodeMessage(kademlia.network.contact, &shortList[0], returnMessage)
	select {
	case kContactsReturned := <- returnMessage:
		contactedContacts = append(contactedContacts, shortList[0])
		/*kContactsFromAplha = channelvärdet från message*/
		var alphaContacts []Contact
		/*Picks alpha first nodes from the k closest*/
		if (len(kContactsReturned)>kademlia.alpha) {
			alphaContacts = append(alphaContacts, kContactsReturned[0:kademlia.alpha]...)
		} else {
			alphaContacts = append(alphaContacts, kContactsReturned...)
		}
		/*Loop through all k nodes*/
		for i:=0; i<len(alphaContacts); i++ {
			isInList := false
			if (len(contactedContacts) > 0) {
				/*Loop through all nodes in shortList*/
				for x:=0; x<len(contactedContacts); x++ {
					/*Check that a node is not already in shortList*/

					if (alphaContacts[i].ID.String() == contactedContacts[x].ID.String()){
						isInList = true
					}
				}
			}
			if (len(nextRound) > 0) {
				/*Loop through all nodes in shortList*/
				for x:=0; x<len(nextRound); x++ {
					/*Check that a node is not already in shortList*/
					if (alphaContacts[i].ID.String() == nextRound[x].ID.String()){
						isInList = true
					}
				}
			}
			if (len(shortList) > 0) {
				for x:=0; x<len(shortList); x++ {
					if (alphaContacts[i].ID.String() == shortList[x].ID.String()){
						isInList = true
					}
				}
			}
			/*If the node is not in list and the node is not currentNode, ping the node and if it is alive,
				add node to shortList
			*/
			if (isInList==false) {
				nextRound = append(nextRound,alphaContacts[i])
				nextRound = kademlia.sortList(nextRound)
			}
		}
	case <-time.After(10*time.Second):
		fmt.Println("TIMEOUT")

	}
	/*Remove first element from list and continue recursion*/
	if (len(shortList) > 1) {
		shortList = shortList[1:]
	} else {
		shortList = make([]Contact, 0)
	}
	return kademlia.NodeLookUpRound(shortList, contactedContacts, nextRound)

}



func (kademlia *Kademlia) sortList(shortList []Contact) []Contact{
	distanceList := make([]KademliaID, len(shortList))
	for i:=0; i<len(shortList); i++{
		distance := kademlia.network.contact.ID.CalcDistance(shortList[i].ID)
		distanceList[i] = *distance
	}
	shortList = quicksort(distanceList, shortList)
	return shortList
}

func quicksort(distanceList []KademliaID, shortList []Contact) []Contact {
    if len(shortList) < 2 {
        return shortList
    }
      
    left, right := 0, len(distanceList)-1
      
    pivot := rand.Int() % len(distanceList)
      
	distanceList[pivot], distanceList[right] = distanceList[right], distanceList[pivot]
	shortList[pivot], shortList[right] = shortList[right], shortList[pivot]
      
    for i, _ := range distanceList {
        if distanceList[i].Less(&distanceList[right]) {
			distanceList[left], distanceList[i] = distanceList[i], distanceList[left]
			shortList[left], shortList[i] = shortList[i], shortList[left]
            left++
        }
    }
	distanceList[left], distanceList[right] = distanceList[right], distanceList[left]
	shortList[left], shortList[right] = shortList[right], shortList[left]
      
    quicksort(distanceList[:left],shortList[:left] )
    quicksort(distanceList[left+1:], shortList[left+1:])
      
    return shortList
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
		contacts := kademlia.network.routingTable.FindClosestContacts(dataKey, kademlia.k)
		addAlphaContacts := make([]Contact, 0)
		/*Picks alpha first nodes from the k closest*/
		if (len(contacts)>kademlia.alpha) {
			addAlphaContacts = append(addAlphaContacts, contacts[0:kademlia.alpha]...)
		} else {
			addAlphaContacts = append(addAlphaContacts, contacts...)
		}
		contactedContacts := make([]Contact, 0)
		result := kademlia.NodeLookUpData(addAlphaContacts, contactedContacts, *dataKey)
		return result
	}
	return nil
}

/*ShortList = Nodes to contact*/
func (kademlia *Kademlia) NodeLookUpData(shortList []Contact, contactedContacts []Contact, dataKey KademliaID) []byte {
	if (len(shortList)==0){
		return nil
	}
	nextRound := make([]Contact, 0)
	var data []byte
	shortList, contactedContacts, data = kademlia.NodeLookUpRoundData(shortList, contactedContacts,nextRound, dataKey)
	if (data!=nil){
		return data
	}
	return kademlia.NodeLookUpData(shortList, contactedContacts, dataKey)
}

/*NodeLookUp function*/
func (kademlia *Kademlia) NodeLookUpRoundData(shortList []Contact, contactedContacts []Contact, nextRound []Contact, dataKey KademliaID) ([]Contact, []Contact, []byte) {
	/*Base case*/
	if (len(shortList)==0) {
		nextRound = kademlia.sortList(nextRound)
		if (len(nextRound)>kademlia.alpha) {
			nextRound = append(nextRound, nextRound[0:kademlia.alpha]...)
		}
		return nextRound, contactedContacts, nil
	}
	returnData:= make(chan []byte)
	go SendFindDataMessage(kademlia.network.contact, &shortList[0], dataKey, returnData)
	select {
	case returnedData := <- returnData:
		if (returnedData!=nil){
			return nextRound, contactedContacts, returnedData
		}
	case <-time.After(10*time.Second):
		fmt.Println("TIMEOUT")
	}
	returnMessage := make(chan []Contact)
	go SendFindNodeMessage(kademlia.network.contact, &shortList[0], returnMessage)
	select {
	case kContactsReturned := <- returnMessage:
		contactedContacts = append(contactedContacts, shortList[0])
		/*kContactsFromAplha = channelvärdet från message*/
		var alphaContacts []Contact
		/*Picks alpha first nodes from the k closest*/
		if (len(kContactsReturned)>kademlia.alpha) {
			alphaContacts = append(alphaContacts, kContactsReturned[0:kademlia.alpha]...)
		} else {
			alphaContacts = append(alphaContacts, kContactsReturned...)
		}
		/*Loop through all k nodes*/
		for i:=0; i<len(alphaContacts); i++ {
			isInList := false
			if (len(contactedContacts) > 0) {
				/*Loop through all nodes in shortList*/
				for x:=0; x<len(contactedContacts); x++ {
					/*Check that a node is not already in shortList*/

					if (alphaContacts[i].ID.String() == contactedContacts[x].ID.String()){
						isInList = true
					}
				}
			}
			if (len(nextRound) > 0) {
				/*Loop through all nodes in shortList*/
				for x:=0; x<len(nextRound); x++ {
					/*Check that a node is not already in shortList*/
					if (alphaContacts[i].ID.String() == nextRound[x].ID.String()){
						isInList = true
					}
				}
			}
			if (len(shortList) > 0) {
				for x:=0; x<len(shortList); x++ {
					if (alphaContacts[i].ID.String() == shortList[x].ID.String()){
						isInList = true
					}
				}
			}
			/*If the node is not in list and the node is not currentNode, ping the node and if it is alive,
				add node to shortList
			*/
			if (isInList==false) {
				nextRound = append(nextRound,alphaContacts[i])
				nextRound = kademlia.sortList(nextRound)
			}
		}
	case <-time.After(10*time.Second):
		fmt.Println("TIMEOUT")

	}
	/*Remove first element from list and continue recursion*/
	if (len(shortList) > 1) {
		shortList = shortList[1:]
	} else {
		shortList = make([]Contact, 0)
	}
	return kademlia.NodeLookUpRoundData(shortList, contactedContacts, nextRound, dataKey)

}

func (kademlia *Kademlia) Store(data []byte, donePut chan bool) {
	fmt.Println("Enter store")
	fmt.Print(data)
	fileKey := NewRandomKademliaID()
	fmt.Println("Filekey: " + fileKey.String())
	clostest := kademlia.LookupContact(*fileKey)
	for i:=0; i<len(clostest); i++ {
		//file := NewFile(*clostestK[i].ID, data, clostestK[i])
		//kademlia.network.StoreDataOnNode(file)
		doneStorePing := make(chan bool)
		go SendStoreMessage(kademlia.network.contact, &clostest[i], data, *fileKey, doneStorePing)
		select {
		case <-doneStorePing:
			fmt.Println("File stored on node: " +clostest[i].ID.String())
		case <-time.After(10*time.Second):
			fmt.Println("TIMEOUT IN STORE")
		}
	}
	donePut <- true


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
