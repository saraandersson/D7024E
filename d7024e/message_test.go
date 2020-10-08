package d7024e

import (
	"fmt"
	"testing"
)

func TestDataReturnMessage(t *testing.T) {
	data := "testingMessage"
	sendData := []byte(data)
	messageType:= "Test"
	messageReturned := createProtoBufDataReturnMessage(sendData, messageType) 
	fmt.Println(messageReturned)
}

func TestMessageForContacts(t *testing.T) {
	contactId := []string{"FFFFFFFF00000000000000000000000000000000", "1111111100000000000000000000000000000000"}
	contactAddress := []string{"localhost:8002", "localhost:8003"}
	messageReturned := createProtoBufMessageForContacts(contactId, contactAddress) 
	fmt.Println(messageReturned)
}

func TestMessage(t *testing.T) {
	senderContact := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8002")
	receiverContact := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8003")
	messageType := "Test"
	messageReturned := createProtoBufMessage(&senderContact, &receiverContact, messageType)
	fmt.Println(messageReturned)
}

func TestMessagePing(t *testing.T) {
	receiverContact := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8003")
	messageType := "Test"
	messageReturned := createProtoBufPingMessage(&receiverContact, messageType)
	fmt.Println(messageReturned)
}

func TestDataMessage(t *testing.T) {
	senderContact := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8003")
	data := "testingMessage"
	sendData := []byte(data)
	key := NewRandomKademliaID()
	messageType := "Test"
	messageReturned := createProtoBufDataMessage(&senderContact, sendData, *key, messageType)
	fmt.Println(messageReturned)
}

func TestFindMessage(t *testing.T) {
	senderContact := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8003")
	key := NewRandomKademliaID()
	messageType := "Test"
	messageReturned := createProtoBufFindMessage(&senderContact, *key, messageType)
	fmt.Println(messageReturned)
}



