package d7024e

import (
	"fmt"
	//"bufio"
	"net"
	//"os"
	"time"
	"strconv"
	"github.com/golang/protobuf/proto"
	"protobuf"
)

type Network struct {
	contact 		*Contact
	routingTable 	*RoutingTable 
	kademlia 		*Kademlia
	fileList 		[]File
}

type File struct {
	key 		KademliaID
	data 		[]byte
}

type Message struct {
	contactsID 		[]string
	contactsAddress []string
}

func NewFile(key KademliaID, data []byte) File{
	file := File{}
	file.key = key
	file.data = data
	return file
}

func (network *Network) FindData(key KademliaID) []byte {
	fmt.Println("Enter FindData")
	if (len(network.fileList)>0){
		fmt.Println("Filelist")
		fmt.Println(network.fileList)
		for i:=0; i<len(network.fileList); i++ {
			if (network.fileList[i].key == key){
				fmt.Println("Hittat data:")
				fmt.Println(network.fileList[i].data)
				return network.fileList[i].data
			}
		}
	}
	return nil
}

func (network *Network) StoreDataOnNode(file File) {
	network.fileList = append(network.fileList, file)
}

func (network *Network) UpdateKBucket(message *protobuf.Message) {
		contact := NewContact(NewKademliaID(message.SenderID), message.SenderAddress)
		/*See if bucket is full -> Ping last node in bucket, if no response, add contact, if response, SKIT I DET!*/
		bucketIndex := network.routingTable.getBucketIndex(contact.ID)
		bucket := network.routingTable.buckets[bucketIndex]
		bucketLen := bucket.Len();
		contactList := bucket.list
		if (bucketLen < 20){
			network.routingTable.AddContact(contact)
		} else {
			fmt.Println("BUCKETLIST")
			fmt.Println(contactList.Back().Value.(Contact).ID)
			fmt.Println(contactList.Back().Value.(Contact).Address)
			bucketLast := network.routingTable.buckets[20]
			contactLast := bucketLast.list.Back().Value.(Contact)
			alive := make(chan bool)
			go SendPingMessage(&contactLast, alive)
			select{
			case <- alive:
			case <- time.After(5*time.Second):
				contactList.Remove(bucketLast.list.Back())
				network.routingTable.AddContact(contact)
			}

		}
}


func (network *Network) Listen(ip string, port int, returnedMessage chan(Message)) {
	port2 := ":" + strconv.Itoa(port)
    s, err := net.ResolveUDPAddr("udp4", port2)
    if err != nil {
            fmt.Println(err)
            return
    }
	connection, err := net.ListenUDP("udp4", s)
    if err != nil {
            fmt.Println(err)
            return
    }
    defer connection.Close()
	for {
		buffer := make([]byte, 1024)
		n, addr, _ := connection.ReadFromUDP(buffer)
		newMessage := &protobuf.Message{}
		errorMessage := proto.Unmarshal(buffer[0:n], newMessage)
		if errorMessage!=nil {
			fmt.Println(errorMessage)
		}

		
		if (newMessage.MessageType == "Ping"){
			
		} else{
			go network.UpdateKBucket(newMessage)
		}
		if (newMessage.MessageType == "FindNode"){
			contacts := network.routingTable.FindClosestContacts(network.contact.ID, 20)
			contactId :=  make([]string, len(contacts))
			contactAddress := make([]string, len(contacts))
			count := 0
			for i:=0; i<len(contacts); i++{
				count = count +1
				contactId[i] = contacts[i].ID.String()
				contactAddress[i] = contacts[i].Address
			}
			message := createProtoBufMessageForContacts(contactId, contactAddress)
			data,_ := proto.Marshal(message)
			_, err = connection.WriteToUDP(data, addr)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if (newMessage.MessageType == "Store") {
			newFile := NewFile(*NewKademliaID(newMessage.Key), newMessage.Data)
			network.fileList = append(network.fileList, newFile)
			//fmt.Println("New file added to node with kademliaID: " + network.contact.ID.String())
		}

		if (newMessage.MessageType == "Find"){
			newKey := NewKademliaID(newMessage.Key)
			fileData := network.FindData(*newKey)
			message := createProtoBufDataReturnMessage(fileData, "Find")
			data,_ := proto.Marshal(message)
			_, err = connection.WriteToUDP(data, addr)
			if err != nil {
					fmt.Println(err)
					return
			}
		}

	}
	
}

func createProtoBufDataReturnMessage(data []byte, messageType string) *protobuf.Message {
	protoBufMessage := &protobuf.Message {
		Data: data,
		MessageType: messageType}

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

func SendPingMessage(receiverContact *Contact, alive chan bool) {
	server, err := net.ResolveUDPAddr("udp4", receiverContact.Address)
	conn, err := net.DialUDP("udp4", nil, server)

	if err != nil {
			fmt.Println(err)
			return
	}
	fmt.Printf("The UDP server is %s\n", conn.RemoteAddr().String())
	defer conn.Close()
	message := createProtoBufPingMessage(receiverContact, "Ping")
	data,_ := proto.Marshal(message)
	_, err = conn.Write(data)
	if err != nil {
			fmt.Println(err)
			return
	}
	alive <- true
}

func SendFindNodeMessage(senderContact *Contact, receiverContact *Contact, returnMessage chan []Contact){
	server, err := net.ResolveUDPAddr("udp4", receiverContact.Address)
	conn, err := net.DialUDP("udp4", nil, server)

	if err != nil {
			fmt.Println(err)
			return
	}
	fmt.Printf("The UDP server is %s\n", conn.RemoteAddr().String())
	defer conn.Close()
	message := createProtoBufMessage(senderContact, receiverContact, "FindNode")
	data,_ := proto.Marshal(message)
	_, err = conn.Write(data)
	if err != nil {
			fmt.Println(err)
			return
	}
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
			fmt.Println(err)
			return
	}
	newMessage := &protobuf.ContactsMessage{}
	errorMessage := proto.Unmarshal(buffer[0:n], newMessage)
	if (errorMessage !=nil){
		fmt.Println(errorMessage)
	}
	contacts := make([]Contact, len(newMessage.ContactsID))
	fmt.Println("LENGTH CONTACTSID")
	fmt.Println(len(newMessage.ContactsID))
	fmt.Println("LENGTH CONTACTSADDRESS")
	fmt.Println(len(newMessage.ContactsAddress))
	for i:=0; i<len(newMessage.ContactsID); i++ {
		contacts[i] = NewContact(NewKademliaID(newMessage.ContactsID[i]), newMessage.ContactsAddress[i])
	}

	returnMessage <- contacts

}

func SendStoreMessage(senderContact *Contact, receiverContact *Contact, data []byte, key KademliaID, donePing chan bool) {
	server, err := net.ResolveUDPAddr("udp4", receiverContact.Address)
	conn, err := net.DialUDP("udp4", nil, server)

	if err != nil {
			fmt.Println(err)
			return
	}
	fmt.Printf("The UDP server is %s\n", conn.RemoteAddr().String())
	defer conn.Close()
	message := createProtoBufDataMessage(senderContact, data, key, "Store")
	dataMarshal,_ := proto.Marshal(message)
	_, err = conn.Write(dataMarshal)
	if err != nil {
			fmt.Println(err)
			return
	}
	donePing <- true
}

func SendFindDataMessage(senderContact *Contact, receiverContact *Contact, key KademliaID, findValue chan []byte) {
	server, err := net.ResolveUDPAddr("udp4", receiverContact.Address)
	conn, err := net.DialUDP("udp4", nil, server)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("The UDP server is %s\n", conn.RemoteAddr().String())
	defer conn.Close()
	message := createProtoBufFindMessage(senderContact, key, "Find")
	dataMarshal,_ := proto.Marshal(message)
	_, err = conn.Write(dataMarshal)
	if err != nil {
			fmt.Println(err)
			return
	}
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
			fmt.Println(err)
			return
	}
	newMessage := &protobuf.Message{}
	errorMessage := proto.Unmarshal(buffer[0:n], newMessage)
	if (errorMessage !=nil){
		fmt.Println(errorMessage)
	}

	findValue <- newMessage.Data
}



func NewMessage(contactsID []string, contactsAddress []string) Message{
	newMessage := Message{}
	newMessage.contactsID = contactsID
	newMessage.contactsAddress = contactsAddress
	return newMessage
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	//rt.AddContact(NewContact(NewKademliaID(currentPacket.SourceID), currentPacket.SourceAddress))
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func NewNetwork(contact Contact, routingTable RoutingTable) Network {
	//sendPing := make(chan bool)
	network := Network{}
	network.contact=&contact
	network.routingTable= &routingTable
	//go network.SendPingMessage(bootstrapContact, sendPing)
	//<- sendPing
	return network
}
