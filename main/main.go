package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
	//"strconv"
)

//import "d7024e"

func main() {
	address := os.Getenv("ADDRESS")
        port := os.Getenv("PORT")
        go mainServer(port) //Gör egen tråd
        //go mainServer(8000, done) 
        <- time.After(1*time.Second)
        server, err := net.ResolveUDPAddr("udp4", address)
        conn, err := net.DialUDP("udp4", nil, server)
        if err != nil {
                fmt.Println(err)
                return
        }
        fmt.Printf("The UDP server is %s\n", conn.RemoteAddr().String())
        defer conn.Close()

        for {
                reader := bufio.NewReader(os.Stdin)
                fmt.Printf("Type message here: ")
                message, _ := reader.ReadString('\n')
                data := []byte(message + "\n")
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
	/*var contact d7024e.Contact
	address := "127.0.0.1:1234"
	id := d7024e.NewRandomKademliaID()
	contact = d7024e.NewContact(id, address)
	fmt.Println(id)
	fmt.Println(contact)
	network := d7024e.CreateNetwork(&contact)
	routingTable := d7024e.NewRoutingTable(contact)

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
			fmt.Print("Enter targetNode id: ")
			targetNode, _ := reader.ReadString('\n')
			//var convertToKademliaId d7024e.KademliaID
			targetNodeId := []byte(targetNode)
			contacts := routingTable.FindClosestContacts(id, 1)
			fmt.Println(contacts)
		}
	}*/
}

func mainServer(port string) {
    //port_input := os.Getenv("PORT")
    port2 := ":" + port
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
            fmt.Print("Type message here: ")
            /*reader := bufio.NewReader(os.Stdin)
            fmt.Print("Type answer here: ")
            text, _ := reader.ReadString('\n')*/
            data := []byte("Hello from 127.0.0.1:" + port + "\n")
            _, err = connection.WriteToUDP(data, addr)
            if err != nil {
                    fmt.Println(err)
                    return
            }
	}
}