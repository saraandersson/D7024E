package main

import (
	//"bufio"
	"fmt"
        //"flag"
	"net"
	"os"
	"time"
        "d7024e"
        //"strconv"
        //"strings"
)

//import "d7024e"

const defaultPort ="8000"

func main() {
        /*Add contacts and create routing tables*/
        bootstrapContact :=  d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
        contact1 := d7024e.NewContact(d7024e.NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002")
        contact2 := d7024e.NewContact(d7024e.NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8003")
        contact3 := d7024e.NewContact(d7024e.NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8004")
        contact4 := d7024e.NewContact(d7024e.NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8005")
        contact5 := d7024e.NewContact(d7024e.NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8006")
        
        ntContact1 := d7024e.NewNetwork(contact1)
        ntContact2 := d7024e.NewNetwork(contact2)
        ntContact3 := d7024e.NewNetwork(contact3)
        ntContact4 := d7024e.NewNetwork(contact4)
        ntContact5 := d7024e.NewNetwork(contact5)

        <- time.After(1*time.Second)

        rtBootstrap := d7024e.NewRoutingTable(bootstrapContact)

        /*rtContact1 := d7024e.NewRoutingTable(contact1)
        rtContact2 := d7024e.NewRoutingTable(contact2)
        rtContact3 := d7024e.NewRoutingTable(contact3)
        rtContact4 := d7024e.NewRoutingTable(contact4)
        rtContact5 := d7024e.NewRoutingTable(contact5)
        rtBootstrap := d7024e.NewRoutingTable(bootstrapContact)*/

        /* */ 
        
        //contact = d7024e.NewContact(d7024e.NewRandomKademliaID(), address)
        //routingTableContact := d7024e.NewRoutingTable(contact)

        /*done := make(chan bool)
        var port = flag.String("port", defaultPort,"Test")
        var bootstrapIP = flag.String("bootstrap_ip", "kademliaBootstrapHost", "Test")
        var bootstrapPort = flag.String("bootstrap_port", defaultPort, "Test")
        flag.Parse()*/
        /*Create contact*/
	//var contact d7024e.Contact
       /* address := GetIPContainer() + ":" + *port
        fmt.Println("addressen for noden: ")
        fmt.Println(address)
        contact := d7024e.NewContact(d7024e.NewRandomKademliaID(), address)
        contactNetwork := d7024e.NewNetwork(contact)
        /*create network and kademlianetwork*/        
       /* bootstrapAddress := *bootstrapIP +":"+ *bootstrapPort
        bootstrapContact := d7024e.NewContact(d7024e.NewRandomKademliaID(), bootstrapAddress)
        network := d7024e.NewNetwork(bootstrapContact)
        bootstrapRoutingTable := d7024e.NewRoutingTable(bootstrapContact)
        kademliaNetwork := d7024e.NewKademlia(&network, &bootstrapContact, bootstrapRoutingTable, 20, 3, done)
        /*Call on LookUpContact*/
        /*i, _ := strconv.Atoi(*port)
        fmt.Println("Här kommer listan:")
        go kademliaNetwork.LookupContact(&contact,&contactNetwork, i)
        <- done
        network.Listen(address, i)*/

       /**/
        done := make(chan bool)
        network := d7024e.NewNetwork(bootstrapContact)
        go network.Listen(bootstrapContact.Address, 8000)
        kademliaNetwork := d7024e.NewKademlia(&network, &bootstrapContact, rtBootstrap, 20, 3, done)
        //go network.Listen("localhost:8000", 8000) 
        //go network.Listen("localhost:8002", 8002) 
        <- time.After(1*time.Second)
        fmt.Println("Här kommer listan:")
        kademliaNetwork.LookupContact(&contact1, &ntContact1, 8002)
        fmt.Println("Här kommer listan:")
        kademliaNetwork.LookupContact(&contact2, &ntContact2, 8002)
        fmt.Println("Här kommer listan:")
        kademliaNetwork.LookupContact(&contact3, &ntContact3, 8002)
        fmt.Println("Här kommer listan:")
        kademliaNetwork.LookupContact(&contact4, &ntContact4, 8002)
        fmt.Println("Här kommer listan:")
        kademliaNetwork.LookupContact(&contact5, &ntContact5, 8002)

}

func GetIPContainer() string{
        containerHostname, _ := os.Hostname()
        addrs, _ := net.LookupHost(containerHostname)
        fmt.Println("Container IP address: " + addrs[0])
        return addrs[0]
}

