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

func (network *Network) MessageHandler(funcType string, message *protobuf.Message) {
	switch funcType {
	case "Ping":
		fmt.Println("Ping! kommer till messagehandler")
		fmt.Println("Network id")
		fmt.Println(network.contact.ID.String())
		fmt.Println("Message.SenderID")
		fmt.Println(message.SenderID)
		contact := NewContact(NewKademliaID(message.SenderID), message.SenderAddress)
		network.routingTable.AddContact(contact)
		fmt.Println("RoutingTable")
		fmt.Println(network.routingTable)
	default:
		fmt.Println("Error in MessageHandler: Wrong message type")
	}
}

func (network *Network) Listen(ip string, port int) {
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
	fmt.Println(connection)
    if err != nil {
            fmt.Println(err)
            return
    }
    defer connection.Close()
	buffer := make([]byte, 1024)
	for {
		n, _, _ := connection.ReadFromUDP(buffer)
		fmt.Println("Listen lyssnar")
		newMessage := &protobuf.Message{}
		errorMessage := proto.Unmarshal(buffer[0:n], newMessage)
		if errorMessage!=nil {
			fmt.Println(errorMessage)
		}
		//fmt.Print("\n", string(buffer[0:n-1]))
		fmt.Print(newMessage)
		go network.MessageHandler("Ping", newMessage)
	}
}

func createProtoBufMessage(senderContact *Contact, receiverContact *Contact) *protobuf.Message {
	/*protoBufMessage := []string {
		network.contact.ID.String(), network.contact.Address, contact.ID.String(), contact.Address}*/
	//var text = "hello"
	protoBufMessage := &protobuf.Message {
			SenderID: senderContact.ID.String(),
			SenderAddress: senderContact.Address,
			ReceiverID: receiverContact.ID.String(),
			ReceiverAddress: receiverContact.Address}
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
	message := createProtoBufMessage(senderContact, receiverContact)
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

func (network *Network) SendFindContactMessage(contact *Contact, ) {
	//rt.AddContact(NewContact(NewKademliaID(currentPacket.SourceID), currentPacket.SourceAddress))
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func NewNetwork(contact Contact) Network {
	//sendPing := make(chan bool)
	network := Network{}
	network.contact=&contact
	network.routingTable= NewRoutingTable(contact)
	//go network.SendPingMessage(bootstrapContact, sendPing)
	//<- sendPing
	return network
}
