package d7024e

import (
	"fmt"
	//"bufio"
	"net"
	//"os"
	"time"
	"strconv"
)

type Network struct {
	contact *Contact
	routingTable *RoutingTable 
	kademlia *Kademlia
}

func Listen(ip string, port int) {
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

    for {
			n, addr, err := connection.ReadFromUDP(buffer)
			fmt.Print("\n" + "Message: ", string(buffer[0:n-1]))
            data := []byte("Hello from " + ip + "\n")
            _, err = connection.WriteToUDP(data, addr)
            if err != nil {
                    fmt.Println(err)
                    return
            }
    }
}

func (network *Network) SendPingMessage(contact *Contact) {
	fmt.Println("Kommer till ping!")
	go Listen(contact.Address, 1234) //Gör egen tråd
	<- time.After(1*time.Second)
	server, err := net.ResolveUDPAddr("udp4", contact.Address)
	conn, err := net.DialUDP("udp4", nil, server)
	if err != nil {
			fmt.Println(err)
			return
	}
	fmt.Printf("The UDP server is %s\n", conn.RemoteAddr().String())
	defer conn.Close()

	for {
			data := []byte("Ping " + "\n")
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
	}
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func NewNetwork(contact Contact, bootstrapContact *Contact) Network {
	network := Network{}
	network.contact=&contact
	network.routingTable= NewRoutingTable(contact)
	go network.SendPingMessage(bootstrapContact)
	return network
}

func (network *Network) NodeLookup(k int, targetNodeId *KademliaID) {
	/*for i:=0; i<contact.length; i++ {
		go FindNode(contacts[i].id)
	}*/

}

func (network *Network) JoinNetwork(target Contact){
	closestTargets := network.routingTable.FindClosestContacts(target.ID, 3)
	fmt.Println(closestTargets)
	network.routingTable.AddContact(target)
	targetNetwork := NewNetwork(target,&target)
	targetNetwork.routingTable.AddContact(network.contact)
	for i:=0; i < len(closestTargets);i++{
		targetNetwork.SendPingMessage(&closestTargets[i])
	}

}

/*func FindNode(id *KademliaID) Network {
	
}*/
