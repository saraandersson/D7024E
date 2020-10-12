package d7024e

import (
	//"github.com/golang/protobuf/proto"
	"protobuf"
)

func createProtoBufDataReturnMessage(data []byte, contactId []string, contactAddress []string) *protobuf.ContactsMessage {
	protoBufMessage := &protobuf.ContactsMessage {
		Data: data,
		ContactsID: contactId,
		ContactsAddress: contactAddress}
	return protoBufMessage
}

func createProtoBufMessageForContacts(contactId []string, contactAddress []string) *protobuf.ContactsMessage {
	protoBufMessage := &protobuf.ContactsMessage{
			ContactsID: contactId,
			ContactsAddress: contactAddress}
	return protoBufMessage
}

func createProtoBufMessage(senderContact *Contact, receiverContact *Contact, messageType string) *protobuf.Message {
	protoBufMessage := &protobuf.Message {
			SenderID: senderContact.ID.String(),
			SenderAddress: senderContact.Address,
			ReceiverID: receiverContact.ID.String(),
			ReceiverAddress: receiverContact.Address,
			MessageType: messageType}
	return protoBufMessage
}

func createProtoBufPingMessage(receiverContact *Contact, messageType string) *protobuf.Message {
	protoBufMessage := &protobuf.Message {
			ReceiverID: receiverContact.ID.String(),
			ReceiverAddress: receiverContact.Address,
			MessageType: messageType}
	return protoBufMessage
}

func createProtoBufDataMessage(senderContact *Contact, data []byte, key KademliaID, messageType string) *protobuf.Message {
	protoBufMessage := &protobuf.Message {
		SenderID: senderContact.ID.String(),
		SenderAddress: senderContact.Address,
		Data: data,
		Key: key.String(),
		MessageType: messageType}

	return protoBufMessage
}

func createProtoBufFindMessage(senderContact *Contact, key KademliaID, messageType string) *protobuf.Message {
	protoBufMessage := &protobuf.Message {
		SenderID: senderContact.ID.String(),
		SenderAddress: senderContact.Address,	
		Key: key.String(),
		MessageType: messageType}

	return protoBufMessage
}