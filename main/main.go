package main

import (
	//"bufio"
    "fmt"
    //"flag"
     "net"
     "os"
     "time"
     "d7024e"
     "strconv"
)

const defaultPort ="8000"

func main() {
        
        bootstrapID := "FFFFFFFF00000000000000000000000000000000"
        bootstrapIP := "172.19.0.2"

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
                 address = bootstrapIP +":"+ defaultPort
                 contact = d7024e.NewContact(d7024e.NewKademliaID(bootstrapID), address)
                 /*Create network, routingtable for bootstrap node*/
                 routingtable := d7024e.NewRoutingTable(contact)
                 network := d7024e.NewNetwork(contact, *routingtable)
                 returnMessage := make(chan d7024e.Message)
                 go network.Listen(address, currentPort, returnMessage)
                 <-time.After(2*time.Second) 

        } else {
                currentPort = 8080
                /*create a new contact for the node*/
                address = GetIPContainer() + ":" + "8080"
                //address = "0.0.0.0:8080"
                contact = d7024e.NewContact(d7024e.NewRandomKademliaID(), address)
                routingtable := d7024e.NewRoutingTable(contact)
                network := d7024e.NewNetwork(contact, *routingtable)
                /*Create kademlia network for bootstrap node*/
                kademliaNetwork := d7024e.NewKademlia(&network, &contact, 20, 3)
                returnMessage := make(chan d7024e.Message)
                go network.Listen(address, currentPort, returnMessage)
                <-time.After(2*time.Second) 
                /*Perform node lookup and network join if not a bootstrap node*/
                bootstrapAddress := bootstrapIP +":"+ defaultPort
                bootstrapContact := d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), bootstrapAddress)
                /*Network joining*/
                routingtable.AddContact(contact)
                routingtable.AddContact(bootstrapContact)
                donePing := make(chan bool)
                boostrapPortPing, _ := strconv.Atoi(defaultPort)
                go d7024e.SendPingMessage(&contact,&bootstrapContact,boostrapPortPing,donePing)
                <- donePing
                lookUpContactResult :=  kademliaNetwork.LookupContact(*contact.ID)
                fmt.Println("Lookup done! Contacts found: ")
                fmt.Println(lookUpContactResult)
                
        }

        for {

        }
}


func GetIPContainer() string{
        containerHostname, _ := os.Hostname()
        addrs, _ := net.LookupHost(containerHostname)
        fmt.Println("Container IP address: " + addrs[0])
        return addrs[0]
}