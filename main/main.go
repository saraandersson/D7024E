package main

import (
	//"bufio"
    "fmt"
    "flag"
     "net"
     "os"
     //"time"
     "d7024e"
     "strconv"
)

const defaultPort ="8000"

func main() {
        /*Input flags*/
        var port = flag.String("port", defaultPort,"Test")
        var bootstrapIP = flag.String("bootstrap_ip", "kademliaBootstrapHost", "Test")
        var bootstrapPort = flag.String("bootstrap_port", defaultPort, "Test")
	flag.Parse()

        bootstrapNodeValue := os.Getenv("BOOTSTRAPNODE")
        var contact d7024e.Contact
        var address string
        var currentPort int
        /*If bootstrapNodeValue == "1" => boostrapNode
          Else => kademliaNode
        */
        if (bootstrapNodeValue == "1"){
                currentPort, _ = strconv.Atoi(defaultPort)
                /*create boostrapcontact*/        
                address = *bootstrapIP +":"+ *bootstrapPort
                contact = d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), address)
                 /*Create network, routingtable for bootstrap node*/

        } else {
                currentPort, _ = strconv.Atoi(*port)
                /*create a new contact for the node*/
                address = GetIPContainer() + ":" + *port
                contact = d7024e.NewContact(d7024e.NewRandomKademliaID(), address)
        }

        network := d7024e.NewNetwork(contact)
        routingtable := d7024e.NewRoutingTable(contact)
        /*Create kademlia network for bootstrap node*/
        done := make(chan bool)
        kademliaNetwork := d7024e.NewKademlia(&network, &contact, routingtable, 20, 3, done)

        //go network.Listen(address, currentPort)
        //<-time.After(2*time.Second) 

        /*Perform node lookup and network join if not a bootstrap node*/
        if (bootstrapNodeValue != "1"){
                fmt.Println("Här kommer listan:")
                bootstrapAddress := *bootstrapIP +":"+ *bootstrapPort
                bootstrapContact := d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), bootstrapAddress)
                routingtable.AddContact(bootstrapContact)
                donePing := make(chan bool)
                boostrapPortPing, _ := strconv.Atoi(defaultPort)
                go network.SendPingMessage(&bootstrapContact,boostrapPortPing,donePing)
                kademliaNetwork.LookupContact(&contact, &network, currentPort)
                <- donePing
        }
}

func GetIPContainer() string{
        containerHostname, _ := os.Hostname()
        addrs, _ := net.LookupHost(containerHostname)
        fmt.Println("Container IP address: " + addrs[0])
        return addrs[0]
}