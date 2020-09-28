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
		fmt.Printf("Ping! kommer till messagehandler")
		contact := NewContact(NewKademliaID(message.SenderID), message.SenderAddress)
		network.routingTable.AddContact(contact)
	default:
		fmt.Println("Error in MessageHandler: Wrong message type")
	}
}

func (network *Network) Listen(ip string, port int) {
	fmt.Println("kommer till listen")
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
    buffer := make([]byte, 1024)
	n, addr, err := connection.ReadFromUDP(buffer)
	newMessage := &protobuf.Message{}
	errorMessage := proto.Unmarshal(buffer[0:n], newMessage)
	if errorMessage!=nil {
		fmt.Println(errorMessage)
	}
	//fmt.Print("\n", string(buffer[0:n-1]))
	fmt.Print(newMessage)
	go network.MessageHandler("Ping", newMessage)
	data := []byte("Hello from " + ip + "\n")
	_, err = connection.WriteToUDP(data, addr)
	if err != nil {
			fmt.Println(err)
			return
	}
}

func (network *Network) createProtoBufMessage(contact *Contact) *protobuf.Message {
	/*protoBufMessage := []string {
		network.contact.ID.String(), network.contact.Address, contact.ID.String(), contact.Address}*/
	//var text = "hello"
	protoBufMessage := &protobuf.Message {
			SenderID: network.contact.ID.String(),
			SenderAddress: network.contact.Address,
			ReceiverID: contact.ID.String(),
			ReceiverAddress: contact.Address}
	return protoBufMessage
}

func (network *Network) SendPingMessage(contact *Contact, port int, donePing chan bool) {
	fmt.Println("Kommer till ping!")
	go network.Listen(contact.Address, port) //Gör egen tråd
	<- time.After(5*time.Second)
	server, err := net.ResolveUDPAddr("udp4", contact.Address)
	conn, err := net.DialUDP("udp4", nil, server)
	if err != nil {
			fmt.Println(err)
			return
	}
	fmt.Printf("The UDP server is %s\n", conn.RemoteAddr().String())
	defer conn.Close()
	message := network.createProtoBufMessage(contact)
	data,_ := proto.Marshal(message)
	//data := []byte("Ping " + "\n")
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
	fmt.Printf("Answer: %s", string(buffer[0:n]))
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
