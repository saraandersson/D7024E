
package main

import (
	//"bufio"
	"fmt"
        //"flag"
	"net"
	"os"
	"time"
        "d7024e"
        //"cli"
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
        contact3 := d7024e.NewContact(d7024e.NewKademliaID("1111111500000000000000000000000000000000"), "localhost:8004")
        contact4 := d7024e.NewContact(d7024e.NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8005")
        contact5 := d7024e.NewContact(d7024e.NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8006")
        
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
        <- time.After(1*time.Second)
        kademliaNetwork := d7024e.NewKademlia(&network, &bootstrapContact, rtBootstrap, 20, 3, done)
        <- time.After(1*time.Second)
        fmt.Println("Här kommer listan:")
        kademliaNetwork.LookupContact(&contact1, &ntContact1, 8002)
        fmt.Println("Här kommer listan:")
        kademliaNetwork.LookupContact(&contact2, &ntContact2, 8003)
        fmt.Println("Här kommer listan:")
        kademliaNetwork.LookupContact(&contact3, &ntContact3, 8004)
        fmt.Println("Här kommer listan:")
        kademliaNetwork.LookupContact(&contact4, &ntContact4, 8005)
        fmt.Println("Här kommer listan:")
        kademliaNetwork.LookupContact(&contact5, &ntContact5, 8002)

       /* cliDone := make(chan bool)
        go cli.CliInput(&kademliaNetwork, cliDone)
        //go cli.ExampleScanner_lines(cliDone)
        <- cliDone*/

        /*Temp store functionality until the other containers work*/


       /* reader := bufio.NewReader(os.Stdin)
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
        //}
*/

}

func GetIPContainer() string{
        containerHostname, _ := os.Hostname()
        addrs, _ := net.LookupHost(containerHostname)
        fmt.Println("Container IP address: " + addrs[0])
        return addrs[0]
}

/*
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
	for {
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
}*/


/*
func (network *Network) SendPingMessage(contact *Contact, port int, donePing chan bool) {
	fmt.Println("Kommer till ping!")
	go network.Listen(contact.Address, port) //Gör egen tråd
	<- time.After(2*time.Second)
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




*/

