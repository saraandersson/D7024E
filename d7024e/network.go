package d7024e

import (
	"fmt"
	"net"
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

func (network *Network) UpdateKBucket(contact Contact) {
		/*See if bucket is full -> Ping last node in bucket, if no response, add contact, if response, skip*/
		bucketIndex := network.routingTable.getBucketIndex(contact.ID)
		bucket := network.routingTable.buckets[bucketIndex]
		bucketLen := bucket.Len();
		contactList := bucket.list
		if (bucketLen < 20){
			network.routingTable.AddContact(contact)
		} else {
			contactLast := bucket.list.Back().Value.(Contact)
			alive := make(chan bool)
			go SendPingMessage(&contactLast, alive)
			select{
			case <- alive:
			case <- time.After(5*time.Second):
				contactList.Remove(bucket.list.Back())
				network.routingTable.AddContact(contact)
			}

		}
}

/*Listen function, handles messageInput*/

func (network *Network) Listen(ip string, port int, returnedMessage chan(Message)) {
	fmt.Println("Kommer till listen")
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
	/*Read udp message and unMarshal*/
	for {
		buffer := make([]byte, 1024)
		n, addr, _ := connection.ReadFromUDP(buffer)
		newMessage := &protobuf.Message{}
		errorMessage := proto.Unmarshal(buffer[0:n], newMessage)
		if errorMessage!=nil {
			fmt.Println(errorMessage)
		}
		if (newMessage.MessageType == "Ping"){
			/*Check if node is alive*/
		} else{
			/*Update K-bucket for node with the sender contact */
			contact := NewContact(NewKademliaID(newMessage.SenderID), newMessage.SenderAddress)
			go network.UpdateKBucket(contact)
		}
		/*Finds k-closest nodes and sends the answer back via udp*/
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
		/*Create new file with the given key and data and add it to the network*/
		if (newMessage.MessageType == "Store") {
			newFile := NewFile(*NewKademliaID(newMessage.Key), newMessage.Data)
			network.fileList = append(network.fileList, newFile)
		}
		/*Finds the node with the given key and sends back nil if data is not found on current node,
		 if data exists, data is sent back to the node via UDP*/
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

/*Classic ping, check if node is alive*/
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

/*Find closest k-nodes to the receivercontact (NodeLookUp) via udp request*/
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
	buffer := make([]byte, 4096)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
			fmt.Println(err)
			return
	}
	newMessage := &protobuf.ContactsMessage{}
	/*Message is returned with a list of contacts id and addresses*/
	errorMessage := proto.Unmarshal(buffer[0:n], newMessage)
	if (errorMessage !=nil){
		fmt.Println(errorMessage)
	}
	contacts := make([]Contact, len(newMessage.ContactsID))
	/*Creates new contacts with the id and address that was returned, adds it to a contact list*/
	for i:=0; i<len(newMessage.ContactsID); i++ {
		contacts[i] = NewContact(NewKademliaID(newMessage.ContactsID[i]), newMessage.ContactsAddress[i])
	}

	returnMessage <- contacts
}

/*Store message, sends a store request with given data and key*/
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

/*Finds data connected to the key and if data is found on node, it is returned, otherwise nil is returned*/
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
	buffer := make([]byte, 4096)
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

/*Creates a new message that contains contacts id and addresses*/
func NewMessage(contactsID []string, contactsAddress []string) Message{
	newMessage := Message{}
	newMessage.contactsID = contactsID
	newMessage.contactsAddress = contactsAddress
	return newMessage
}

/*Create a new network for given contact*/
func NewNetwork(contact Contact, routingTable RoutingTable) Network {
	network := Network{}
	network.contact=&contact
	network.routingTable= &routingTable
	return network
}
