package main

import (
	//"bufio"
	"fmt"
        "flag"
	//"net"
	//"os"
	//"time"
        //"d7024e"
	//"strconv"
)

//import "d7024e"

const defaultPort ="8000"
func main() {

        var port = flag.String("port", defaultPort,"specify port for the connections.")
        var bootstrapIP = flag.String("bootstrap_ip", "kademlia_bootstrap_host","The bootstrap node IP address to join")
        var bootstrapPort = flag.String("bootstrap_port", defaultPort, "The port of bootstrap node")
        flag.Parse()

        fmt.Println("Bootstrap Address: "+ *bootstrapIP + ":" + *bootstrapPort)
        fmt.Println("Node is listeneing to Port" + *port)

        
        /*
	var contact d7024e.Contact
	address := "127.0.0.1:1234"
	id := d7024e.NewRandomKademliaID()
	contact = d7024e.NewContact(id, address)
	fmt.Println(id)
	fmt.Println(contact)
	network := d7024e.CreateNetwork(&contact)
	//routingTable := d7024e.NewRoutingTable(contact)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Type operation here: ")
	for {
		text, _ := reader.ReadString('\n')
		if text == "ping" {
			network.SendPingMessage(&contact)
		}
		if text == "join network" {
			//Questions: 
			// How to use and create a bootstrap node?
			// How to implement? 
			// Why do we have the same kademliaid of all containers?

		}
		if text == "node lookup" {
			//fmt.Print("Enter targetNode id: ")
			//targetNode, _ := reader.ReadString('\n')
			//var convertToKademliaId d7024e.KademliaID
			//targetNodeId := []byte(targetNode)
			//contacts := routingTable.FindClosestContacts(id, 1)
			//fmt.Println(contacts)
		}
	}
        */
}

