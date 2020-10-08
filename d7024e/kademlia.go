package d7024e

import (
	"fmt"
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


/*Node lookup*/
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

	/*Performs node lookup with alpha contacts as input, and an empty list as contactedContacts*/
	contactedContacts := make([]Contact, 0)
	contactsToAdd := kademlia.NodeLookUp(addAlphaContacts, contactedContacts)
	/*Updates k-buckets for current node with the list contactsToAdd*/
	for i:=len(contactsToAdd)-1; i>0; i-- {
		kademlia.network.UpdateKBucket(contactsToAdd[i]);
	}
	return contactsToAdd
}
/*ShortList = Nodes to contact, contactedContaces = node that has been contacted*/
func (kademlia *Kademlia) NodeLookUp(shortList []Contact, contactedContacts []Contact) []Contact {
	/*If no nodes left to contact, sort the contactedContact list and return the sorted list*/
	if (len(shortList)==0){
		contactedContacts = kademlia.sortList(contactedContacts)
		return contactedContacts
	}
	/*Perform a node lookup for the current round that returns a updated shortList and contactedContacts,
	 and calls recursively with the returned lists*/
	nextRoundShortList := make([]Contact, 0)
	shortList, contactedContacts = kademlia.NodeLookUpRound(shortList, contactedContacts,nextRoundShortList)
	return kademlia.NodeLookUp(shortList, contactedContacts)
}

/*NodeLookUpRound function*/
/*ShortList = Nodes to contact, contactedContacts = nodes that has been contacted,
 nextRoundShortList = contains the next round of nodes that should be contacted*/
func (kademlia *Kademlia) NodeLookUpRound(shortList []Contact, contactedContacts []Contact, nextRoundShortList []Contact) ([]Contact, []Contact) {
	/*Base case, check if shortList has any nodes to contact*/
	if (len(shortList)==0) {
		/*Sorts the list and returns alpha nodes to contact and contactedContacts*/
		nextRoundShortList = kademlia.sortList(nextRoundShortList)
		if (len(nextRoundShortList)>kademlia.alpha) {
			nextRoundShortList = append(nextRoundShortList, nextRoundShortList[0:kademlia.alpha]...)
		}
		return nextRoundShortList, contactedContacts
	}

	/*Perform FIND-NODE RPC, returns the k-clostest nodes to the target node*/
	returnMessage := make(chan []Contact)
	go SendFindNodeMessage(kademlia.network.contact, &shortList[0], returnMessage)
	select {
	case kContactsReturned := <- returnMessage:
		/*Since first contact in shortList has been contacted, it is added to contactedContacts*/
		contactedContacts = append(contactedContacts, shortList[0])
		var alphaContacts []Contact
		/*Picks alpha first nodes from the k closest to the target node*/
		if (len(kContactsReturned)>kademlia.alpha) {
			alphaContacts = append(alphaContacts, kContactsReturned[0:kademlia.alpha]...)
		} else {
			alphaContacts = append(alphaContacts, kContactsReturned...)
		}
		/*Loop through all alpha nodes from the target node*/
		for i:=0; i<len(alphaContacts); i++ {
			isInList := false
			if (len(contactedContacts) > 0) {
				/*Loop through all nodes in contactedContacts*/
				for x:=0; x<len(contactedContacts); x++ {
					/*Check that a node is not already in contactedContacts,
					 to make sure that we do not contact a node twice*/
					if (alphaContacts[i].ID.String() == contactedContacts[x].ID.String()){
						isInList = true
					}
				}
			}
			if (len(nextRoundShortList) > 0) {
				/*Loop through all nodes in nextRoundShortList*/
				for x:=0; x<len(nextRoundShortList); x++ {
					/*Check that a node is not already in nextRoundShortList*/
					if (alphaContacts[i].ID.String() == nextRoundShortList[x].ID.String()){
						isInList = true
					}
				}
			}
			if (len(shortList) > 0) {
				/*Loop through all nodes in shortList*/
				for x:=0; x<len(shortList); x++ {
					/*Check that a node is not already in shortList*/
					if (alphaContacts[i].ID.String() == shortList[x].ID.String()){
						isInList = true
					}
				}
			}
			/*If the node is not in any of the above lists, add node to nextRoundShortList
			*/
			if (isInList==false) {
				nextRoundShortList = append(nextRoundShortList,alphaContacts[i])
				nextRoundShortList = kademlia.sortList(nextRoundShortList)
			}
		}
	/*If no answer from FIND-NODE RPC, perform timeout*/
	case <-time.After(5*time.Second):
		fmt.Println("TIMEOUT")

	}
	/*Remove first element from list and continue recursion*/
	if (len(shortList) > 1) {
		shortList = shortList[1:]
	} else {
		shortList = make([]Contact, 0)
	}
	return kademlia.NodeLookUpRound(shortList, contactedContacts, nextRoundShortList)

}


/*Sorts a list of contacts based on distance*/
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

/*Find object with help of node lookup for the given key*/
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
		/*Performs node lookup with alpha contacts as input, and an empty list as contactedContacts*/
		contactedContacts := make([]Contact, 0)
		/*NodeLookUpData returns the data if found, else it returns nil*/
		result := kademlia.NodeLookUpData(addAlphaContacts, contactedContacts, *dataKey)
		return result
	}
	return nil
}

/*ShortList = Nodes to contact, contactedContaces = node that has been contacted*/
func (kademlia *Kademlia) NodeLookUpData(shortList []Contact, contactedContacts []Contact, dataKey KademliaID) []byte {
	/*If no nodes left to contact, return nil*/
	if (len(shortList)==0){
		return nil
	}
	nextRoundShortList := make([]Contact, 0)
	var data []byte
	shortList, contactedContacts, data = kademlia.NodeLookUpRoundData(shortList, contactedContacts,nextRoundShortList, dataKey)
	if (data!=nil){
		return data
	}
	/*Perform a node lookup for the current round that returns a updated shortList and contactedContacts,
	 and calls recursively with the returned lists*/
	return kademlia.NodeLookUpData(shortList, contactedContacts, dataKey)
}

/*NodeLookUpRound function*/
/*ShortList = Nodes to contact, contactedContacts = nodes that has been contacted, dataKey = key to the requested data
 nextRoundShortList = contains the next round of nodes that should be contacted*/
func (kademlia *Kademlia) NodeLookUpRoundData(shortList []Contact, contactedContacts []Contact, nextRoundShortList []Contact, dataKey KademliaID) ([]Contact, []Contact, []byte) {
	/*Base case, check if shortList has any nodes to contact*/
	if (len(shortList)==0) {
		/*Sorts the list and returns alpha nodes to contact and contactedContacts*/
		nextRoundShortList = kademlia.sortList(nextRoundShortList)
		if (len(nextRoundShortList)>kademlia.alpha) {
			nextRoundShortList = append(nextRoundShortList, nextRoundShortList[0:kademlia.alpha]...)
		}
		return nextRoundShortList, contactedContacts, nil
	}
	/*Perform FIND-DATA RPC, returns data if found*/
	returnData:= make(chan []byte)
	go SendFindDataMessage(kademlia.network.contact, &shortList[0], dataKey, returnData)
	select {
	case returnedData := <- returnData:
		if (returnedData!=nil){
			return nextRoundShortList, contactedContacts, returnedData
		}
	case <-time.After(10*time.Second):
		fmt.Println("TIMEOUT")
	}
	returnMessage := make(chan []Contact)
	go SendFindNodeMessage(kademlia.network.contact, &shortList[0], returnMessage)
	select {
	case kContactsReturned := <- returnMessage:
		contactedContacts = append(contactedContacts, shortList[0])
		var alphaContacts []Contact
		if (len(kContactsReturned)>kademlia.alpha) {
			alphaContacts = append(alphaContacts, kContactsReturned[0:kademlia.alpha]...)
		} else {
			alphaContacts = append(alphaContacts, kContactsReturned...)
		}
		for i:=0; i<len(alphaContacts); i++ {
			isInList := false
			if (len(contactedContacts) > 0) {
				for x:=0; x<len(contactedContacts); x++ {
					if (alphaContacts[i].ID.String() == contactedContacts[x].ID.String()){
						isInList = true
					}
				}
			}
			if (len(nextRoundShortList) > 0) {
				for x:=0; x<len(nextRoundShortList); x++ {
					if (alphaContacts[i].ID.String() == nextRoundShortList[x].ID.String()){
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
			if (isInList==false) {
				nextRoundShortList = append(nextRoundShortList,alphaContacts[i])
				nextRoundShortList = kademlia.sortList(nextRoundShortList)
			}
		}
	case <-time.After(10*time.Second):
		fmt.Println("TIMEOUT")

	}
	if (len(shortList) > 1) {
		shortList = shortList[1:]
	} else {
		shortList = make([]Contact, 0)
	}
	return kademlia.NodeLookUpRoundData(shortList, contactedContacts, nextRoundShortList, dataKey)

}

/*Store functino, sends a STORE-RPC to the closest nodes found*/
func (kademlia *Kademlia) Store(data []byte, donePut chan bool) {
	fmt.Println("Enter store")
	fmt.Print(data)
	fileKey := NewRandomKademliaID()
	fmt.Println("Filekey: " + fileKey.String())
	clostest := kademlia.LookupContact(*fileKey)
	for i:=0; i<len(clostest); i++ {
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
}

/*Creates a new kademlia*/
func NewKademlia(network *Network, contact *Contact, k int, alpha int) Kademlia{
	kademlia := Kademlia{}
	kademlia.network=network
	kademlia.contact = contact
	kademlia.routingTable = network.routingTable
	kademlia.k = k
	kademlia.alpha = alpha
	return kademlia
}
