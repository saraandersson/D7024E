package main

import (
	"bufio"
	"fmt"
        //"flag"
	"net"
	"os"
	"time"
        "d7024e"
        //"strconv"
        "strings"
)

//import "d7024e"

const defaultPort ="8000"

func main() {
        /*Add contacts and create routing tables*/
        bootstrapContact :=  d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
        contact1 := d7024e.NewContact(d7024e.NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002")
        contact2 := d7024e.NewContact(d7024e.NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002")
        contact3 := d7024e.NewContact(d7024e.NewKademliaID("1111111500000000000000000000000000000000"), "localhost:8002")
        contact4 := d7024e.NewContact(d7024e.NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002")
        contact5 := d7024e.NewContact(d7024e.NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002")
        
        ntContact1 := d7024e.NewNetwork(contact1)
        ntContact2 := d7024e.NewNetwork(contact2)
        ntContact3 := d7024e.NewNetwork(contact3)
        ntContact4 := d7024e.NewNetwork(contact4)
        ntContact5 := d7024e.NewNetwork(contact5)

        <- time.After(1*time.Second)

        rtBootstrap := d7024e.NewRoutingTable(bootstrapContact)
        done := make(chan bool)
        network := d7024e.NewNetwork(bootstrapContact)
        go network.Listen(bootstrapContact.Address, 8000)
        kademliaNetwork := d7024e.NewKademlia(&network, &bootstrapContact, rtBootstrap, 20, 3, done)
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

        /*Temp store functionality until the other containers work*/


        reader := bufio.NewReader(os.Stdin)
        fmt.Println("Type operation below, you can choose between following: store, find, put, get, exit")
        fmt.Println("---------------------")

  //for {
        fmt.Print("-> ")
                text, _ := reader.ReadString('\n')
                text = strings.TrimRight(text, "\n")
                
                switch text {
                        case "store":
                                fmt.Println("enter store")
                                
                                fmt.Print("Enter data to store: ")
                                data, _ := reader.ReadString('\n')
                                data = strings.TrimRight(text, "\n")
                                sendData := []byte(data + "\n")
                                kademliaNetwork.Store(sendData)
                        case "find":
                                fmt.Println("enter find")
                        case "put":
                                fmt.Println("enter put")
                        case "get":
                                fmt.Println("enter get")
                        case "exit":
                                fmt.Println("enter exit")
                        default:
                                fmt.Println("Please type correct operation, you can choose between following: store, find, put, get, exit")
                }


}

func GetIPContainer() string{
        containerHostname, _ := os.Hostname()
        addrs, _ := net.LookupHost(containerHostname)
        fmt.Println("Container IP address: " + addrs[0])
        return addrs[0]
}

