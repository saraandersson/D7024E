package d7024e

import (
	"fmt"
	//"bufio"
	"net"
	//"os"
	//"time"
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
	contact 	Contact
}

type Message struct {
	contacts []Contacts
}

func NewFile(key KademliaID, data []byte, contact Contact) File{
	file := File{}
	file.key = key
	file.data = data
	file.contact = contact
	fmt.Println("FILE:")
	fmt.Println(file)
	return file
}

func (network *Network) StoreDataOnNode(file File) {
	network.fileList = append(network.fileList, file)
	fmt.Println("file-listan")
	fmt.Println(network.fileList)
}

func (network *Network) MessageHandlerPing(message *protobuf.Message) {
		contact := NewContact(NewKademliaID(message.SenderID), message.SenderAddress)
		network.routingTable.AddContact(contact)
}

func (network *Network) MessageHandlerAlpha(message *protobuf.Message, returnedMessage chan(Message)){
	returnedMessage <- network.routingTable.FindClosestContacts(network.contact.ID, 20)
}

func (network *Network) Listen(ip string, port int, returnedMessage chan(Message)) {
	fmt.Println("kommer till listen")
	fmt.Println("ip i listen:")
	fmt.Println(ip)
	fmt.Println("Port i listen")
	fmt.Println(port)
	port2 := ":" + strconv.Itoa(port)
    s, err := net.ResolveUDPAddr("udp4", port2)
    if err != nil {
            fmt.Println(err)
            return
    }
	connection, err := net.ListenUDP("udp4", s)
	fmt.Println("connection i listen")
	//fmt.Println(connection.LocalAddr().String())
    if err != nil {
            fmt.Println(err)
            return
    }
    defer connection.Close()
	for {
		fmt.Println("Listen lyssnar")
		buffer := make([]byte, 1024)
		n, addr, err := connection.ReadFrom(buffer)
		newMessage := &protobuf.Message{}
		errorMessage := proto.Unmarshal(buffer[0:n], newMessage)
		if addr != nil {
			fmt.Println("addressen")
			fmt.Println(addr)
		}
		if err != nil {
			fmt.Println(err)
		}
		if errorMessage!=nil {
			fmt.Println(errorMessage)
		}
		//fmt.Print("\n", string(buffer[0:n-1]))
		fmt.Print(newMessage)
		if (newMessage.MessageType == "Ping"){
			go network.MessageHandlerPing(newMessage)
		}
		if (newMessage.MessageType == "Alpha"){
			returnedMessage.contacts <- network.routingTable.FindClosestContacts(network.contact.ID, 20)
			data := []byte(returnedMessage.contacts)
            _, err = connection.WriteToUDP(data, addr)
            if err != nil {
                    fmt.Println(err)
                    return
            }
		}

	}
	fmt.Println("exit for loop in listen")
}

func createProtoBufMessage(senderContact *Contact, receiverContact *Contact, messageType string) *protobuf.Message {
	/*protoBufMessage := []string {
		network.contact.ID.String(), network.contact.Address, contact.ID.String(), contact.Address}*/
	//var text = "hello"
	protoBufMessage := &protobuf.Message {
			SenderID: senderContact.ID.String(),
			SenderAddress: senderContact.Address,
			ReceiverID: receiverContact.ID.String(),
			ReceiverAddress: receiverContact.Address
			MessageType: messageType}
	return protoBufMessage
}

func SendPingMessage(senderContact *Contact, receiverContact *Contact, port int, donePing chan bool) {
	fmt.Println("Kommer till ping!")
	/*go network.Listen(contact.Address, port) //Gör egen tråd
	<- time.After(2*time.Second)*/
	fmt.Println("Kontaktaddress")
	fmt.Println(receiverContact.Address)
	server, err := net.ResolveUDPAddr("udp4", receiverContact.Address)
	conn, err := net.DialUDP("udp4", nil, server)

	if err != nil {
			fmt.Println(err)
			return
	}
	fmt.Printf("The UDP server is %s\n", conn.RemoteAddr().String())
	defer conn.Close()
	message := createProtoBufMessage(senderContact, receiverContact, "Ping")
	data,_ := proto.Marshal(message)
	fmt.Println("meddelande som skickas")
	fmt.Println(message)
	//data := []byte("Ping " + "\n")
	_, err = conn.Write(data)
	if err != nil {
			fmt.Println(err)
			return
	}
	donePing <- true
}

func SendFindAlphaMessage(senderContact *Contact, receiverContact *Contact, port int, donePing chan bool, returnMessage chan(Message)){
	fmt.Println("Kommer till sendFindAlphaMessage!")
	/*go network.Listen(contact.Address, port) //Gör egen tråd
	<- time.After(2*time.Second)*/
	fmt.Println("Kontaktaddress")
	fmt.Println(receiverContact.Address)
	server, err := net.ResolveUDPAddr("udp4", receiverContact.Address)
	conn, err := net.DialUDP("udp4", nil, server)

	if err != nil {
			fmt.Println(err)
			return
	}
	fmt.Printf("The UDP server is %s\n", conn.RemoteAddr().String())
	defer conn.Close()
	message := createProtoBufMessage(senderContact, receiverContact, "Alpha")
	data,_ := proto.Marshal(message)
	fmt.Println("meddelande som skickas")
	fmt.Println(message)
	//data := []byte("Ping " + "\n")
	_, err = conn.Write(data)
	if err != nil {
			fmt.Println(err)
			return
	}


	donePing <- true
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
