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
	returnedContacts := make([]Contact, 0)
	returnedContactedContacts := make([]Contact, 0)
	sortedContactsToAdd := make([]Contact, 0)
	contactsToAdd := make([]Contact, 0)
	returnMessage := make(chan []Contact)
	contactedContactsRound := make(chan []Contact)
	for i := range shortList { 
		go kademlia.NodeLookUpRound(shortList[i], contactedContactsRound, returnMessage)
		returnedContacts = <- returnMessage
		returnedContactedContacts = <- contactedContactsRound
	}
	fmt.Println("returnedContacts")
	fmt.Println(returnedContacts)
	fmt.Println("ConcatedContacts")
	fmt.Println(contactedContactsRound)
	
	for i := range returnedContacts { 
		contactsToAdd = append(contactsToAdd, returnedContacts[i])
	}
	for i := range returnedContactedContacts { 
		contactedContacts = append(contactedContacts, returnedContactedContacts[i])
	}
	
	fmt.Println("contactsToAdd")
	fmt.Println(contactsToAdd)
	fmt.Println("ConcatedContacts")
	fmt.Println(contactedContacts)

	for i:=0; i<len(contactsToAdd); i++ {
		isInList := false
		if (len(contactedContacts) > 0) {
			for x:=0; x<len(contactedContacts); x++ {
				if (contactsToAdd[i].ID.String() == contactedContacts[x].ID.String()){
					isInList = true
				}
			}
		}
		if (len(shortList) > 0) {
			for x:=0; x<len(shortList); x++ {
				if (contactsToAdd[i].ID.String() == shortList[x].ID.String()){
					isInList = true
				}
			}
		}

		if (isInList == false){
			sortedContactsToAdd = append(sortedContactsToAdd, contactsToAdd[i])
		}
	}
	sortedContactsToAdd = removeDuplicates(sortedContactsToAdd)
	sortedContactsToAdd = kademlia.sortList(contactsToAdd)
	//sortedContactsToAdd = unique(sortedContactsToAdd)

	/*Picks alpha first nodes from the k closest*/
	if (len(sortedContactsToAdd)>kademlia.alpha) {
		sortedContactsToAdd = append(sortedContactsToAdd, sortedContactsToAdd[0:kademlia.alpha]...)
	} else {
		sortedContactsToAdd = append(sortedContactsToAdd, sortedContactsToAdd...)
	}
	/*Perform a node lookup for the current round that returns a updated shortList and contactedContacts,
	 and calls recursively with the returned lists*/
	//nextRoundShortList := make([]Contact, 0)
	//shortList, contactedContacts = kademlia.NodeLookUpRound(shortList, contactedContacts,nextRoundShortList)
	fmt.Println("INNAN RETURN sortedContacts")
	fmt.Println(sortedContactsToAdd)
	return kademlia.NodeLookUp(sortedContactsToAdd, contactedContacts)
}

/*func unique(contactedContacts []Contact) []Contact {
    keys := make(map[Contact.ID]bool)
    list := []Contact{} 
    for _, entry := range contactedContacts {
        if _, value := keys[entry]; !value {
            keys[entry] = true
            list = append(list, entry)
        }
    }    
    return list
}*/

func removeDuplicates(contactedContacts []Contact) []Contact {
	resultList := make([]Contact, 0)
	for i := range contactedContacts {
		isInList := false
		newList := contactedContacts[i+1:]
		for x := range newList {
			if (contactedContacts[i] == newList[x]){
				isInList = true
			} 
		}
		if (isInList == false){
			resultList = append(resultList, contactedContacts[i])
		}
	}
	return resultList
}


/*NodeLookUpRound function*/
/*ShortList = Nodes to contact, contactedContacts = nodes that has been contacted,
 nextRoundShortList = contains the next round of nodes that should be contacted*/
func (kademlia *Kademlia) NodeLookUpRound(contact Contact, contactedContact chan []Contact, kClosestReturned chan []Contact) {
	/*Base case, check if shortList has any nodes to contact*/

	/*Perform FIND-NODE RPC, returns the k-clostest nodes to the target node*/
	returnMessage := make(chan []Contact)
	//kClosestReturnedFromFindNode := make([]Contact, 0)
	go SendFindNodeMessage(kademlia.network.contact, &contact, returnMessage)
	select {
	case kClosestReturnedFromFindNode := <- returnMessage:
		kClosestReturned <- kClosestReturnedFromFindNode
		contactedContacts := make([]Contact, 0)
		contactedContacts = append(contactedContacts, contact)
		contactedContact <- contactedContacts
		/*If no answer from FIND-NODE RPC, perform timeout*/
	case <-time.After(5*time.Second):
		fmt.Println("TIMEOUT")
		returnList := make([]Contact, 0)
		kClosestReturned <- returnList
	}
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
