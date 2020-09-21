package main

import (
	//"bufio"
	"fmt"
        "flag"
	"net"
	"os"
	"time"
        "d7024e"
	"strconv"
)

//import "d7024e"

const defaultPort ="8000"

func main() {
        done := make(chan bool)
        var port = flag.String("port", defaultPort,"specify port for the connections.")
        var bootstrapIP = flag.String("bootstrap_ip", "kademliaBootstrapHost","The bootstrap node IP address to join")
        var bootstrapPort = flag.String("bootstrap_port", defaultPort, "The port of bootstrap node")
        flag.Parse()

        fmt.Println("Bootstrap Address: "+ *bootstrapIP + ":" + *bootstrapPort)
        fmt.Println("Node is listeneing to Port " + *port)

        
        
	var contact d7024e.Contact
        address := GetIPContainer() + ":" + *port
        fmt.Println("addressen for noden: ")
        fmt.Println(address)
        contact = d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
        routingTableContact := d7024e.NewRoutingTable(contact)

        bootstrapAddress := *bootstrapIP +":"+ *bootstrapPort
        bootstrapContact := d7024e.NewContact(d7024e.NewRandomKademliaID(), bootstrapAddress)
        network := d7024e.NewNetwork(bootstrapContact, &bootstrapContact)
        bootstrapRoutingTable := d7024e.NewRoutingTable(bootstrapContact)
        kademliaNetwork := d7024e.NewKademlia(&network, &bootstrapContact, bootstrapRoutingTable, 20, 3, done)
        go d7024e.Listen(address, 8080)
        <- time.After(1*time.Second)
        //lookupContact := d7024e.NewContact(d7024e.NewRandomKademliaID(), "0.0.0.0:"+ *port)
        //kademliaNetwork.routingTable.AddContact(bootstrapContact)
        //closestTargets := network.routingTable.FindClosestContacts(contact.ID, 3)
        i, _ := strconv.Atoi(*port)
        go kademliaNetwork.LookupContact(&contact, routingTableContact, i)
        fmt.Println("Här kommer listan:" )
        <- done

        //fmt.Println(closestTargets)

        /*go network.JoinNetwork(contact, *routingTable, test)
        <- test*/
        //kademliaNetwork.LookupContact(&lookupContact)
        //network.Listen(*port)
        
/*
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

func GetIPContainer() string{
        containerHostname, _ := os.Hostname()
        addrs, _ := net.LookupHost(containerHostname)
        fmt.Println("Container IP address: " + addrs[0])
        return addrs[0]
}

