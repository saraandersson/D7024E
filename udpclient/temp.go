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

const defaultPort ="8000"

func main() {
        /*Add contacts and create routing tables*/
        var port = flag.String("port", defaultPort,"Test")
        var bootstrapIP = flag.String("bootstrap_ip", "kademliaBootstrapHost", "Test")
        var bootstrapPort = flag.String("bootstrap_port", defaultPort, "Test")
        flag.Parse()
        address := GetIPContainer() + ":" + *port
        fmt.Println("addressen for noden: ")
        fmt.Println(address)
        contact := d7024e.NewContact(d7024e.NewRandomKademliaID(), address)
        contactNetwork := d7024e.NewNetwork(contact)
        /*create network and kademlianetwork*/        
        bootstrapAddress := *bootstrapIP +":"+ *bootstrapPort
        bootstrapContact := d7024e.NewContact(d7024e.NewRandomKademliaID(), bootstrapAddress)
        network := d7024e.NewNetwork(bootstrapContact)
		bootstrapRoutingTable := d7024e.NewRoutingTable(bootstrapContact)
		done := make(chan bool)
        kademliaNetwork := d7024e.NewKademlia(&network, &bootstrapContact, bootstrapRoutingTable, 20, 3, done)
        /*Call on LookUpContact*/
        i, _ := strconv.Atoi(*port)
        go network.Listen(bootstrapAddress, 8000)
        <- time.After(1*time.Second)
        fmt.Println("Här kommer listan:")
        kademliaNetwork.LookupContact(&contact, &contactNetwork, i)
}

func GetIPContainer() string{
        containerHostname, _ := os.Hostname()
        addrs, _ := net.LookupHost(containerHostname)
        fmt.Println("Container IP address: " + addrs[0])
        return addrs[0]
}

