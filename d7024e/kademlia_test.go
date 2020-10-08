package d7024e

import (
	"fmt"
	"testing"
	"time"
	"bytes"

)

func TestLookUpContact(t *testing.T) {
	
	contact := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	rt := NewRoutingTable(contact)
	network := NewNetwork(contact, *rt)
	returnMessage := make(chan Message)
	go network.Listen("localhost:8000", 8000, returnMessage)
	<-time.After(2*time.Second) 
	/*Create kademlia network for bootstrap node*/
	kademliaNetwork := NewKademlia(&network, &contact, 20, 3)
	rt.AddContact(contact)

	contact1 := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8001")
	rt1 := NewRoutingTable(contact1)
	network1 := NewNetwork(contact1, *rt1)
	returnMessage1 := make(chan Message)
	go network1.Listen("localhost:8001", 8001, returnMessage1)
	<-time.After(2*time.Second) 
	rt.AddContact(contact1)

	contact2 := NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002")
	rt2 := NewRoutingTable(contact2)
	network2 := NewNetwork(contact2, *rt2)
	returnMessage2 := make(chan Message)
	go network2.Listen("localhost:8002", 8002, returnMessage2)
	<-time.After(2*time.Second)
	rt.AddContact(contact2)


	contact3 := NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8003")
	rt3 := NewRoutingTable(contact3)
	network3 := NewNetwork(contact3, *rt3)
	returnMessage3 := make(chan Message)
	go network3.Listen("localhost:8003", 8003, returnMessage3)
	<-time.After(2*time.Second)
	rt.AddContact(contact3)

	returnContacts := kademliaNetwork.LookupContact(*contact.ID)
	fmt.Println(returnContacts)

}

func TestStore(t *testing.T) {
	contact := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8004")
	rt := NewRoutingTable(contact)
	network := NewNetwork(contact, *rt)
	returnMessage := make(chan Message)
	go network.Listen("localhost:8004", 8004, returnMessage)
	<-time.After(2*time.Second) 
	kademliaNetwork := NewKademlia(&network, &contact, 20, 3)
	rt.AddContact(contact)

	contact1 := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8005")
	rt1 := NewRoutingTable(contact1)
	network1 := NewNetwork(contact1, *rt1)
	returnMessage1 := make(chan Message)
	go network1.Listen("localhost:8005", 8005, returnMessage1)
	<-time.After(2*time.Second) 
	rt.AddContact(contact1)

	data := "test"
	sendData := []byte(data)
	donePut := make(chan bool)
	go kademliaNetwork.Store(sendData, donePut)
	select {
	case <-donePut:
		fmt.Println("Put is done!")
		break
	}

}

func TestFind(t *testing.T){
	contact := NewContact(NewRandomKademliaID(), "localhost:8006")
	rt := NewRoutingTable(contact)
	network := NewNetwork(contact, *rt)
	returnMessage := make(chan Message)
	go network.Listen("localhost:8006", 8006, returnMessage)
	<-time.After(2*time.Second) 
	kademliaNetwork := NewKademlia(&network, &contact, 20, 3)
	rt.AddContact(contact)
	data := "test"
	sendData := []byte(data)
	fileKey := NewRandomKademliaID()
	file := NewFile(*fileKey, sendData)
	network.StoreDataOnNode(file)
	fileReturned := kademliaNetwork.LookupData(fileKey.String())
	equal := bytes.Equal(sendData, fileReturned)
	if (equal){
		fmt.Println(fileReturned)
	} else{
		t.Errorf("Correct data was not returned in Find Data!")
	}
}

func TestSortList(t *testing.T){
	contact := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	rt := NewRoutingTable(contact)
	network := NewNetwork(contact, *rt)
	/*Create kademlia network for bootstrap node*/
	kademliaNetwork := NewKademlia(&network, &contact, 20, 3)
	contact1 := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8001")
	contact2 := NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002")
	contacts := []Contact{contact, contact1, contact2}
	returnedContacts := kademliaNetwork.sortList(contacts)
	fmt.Println(returnedContacts)
}


