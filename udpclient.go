package main

import (
	"fmt"
    "net"
    //"time"
    //"os"
    "strconv"
)

type Node struct {
    address string
    connection *net.UDPConn
}


func main() {
    numberOfNodes := 4;
    port := 8000;
    for i := 0; i < numberOfNodes; i++ {
        fmt.Println("Enter for-loop")
        newPort = port + i;
        newNode := createNewNode("localhost:" + strconv.Itoa(newPort))
        go newNode.checkNodeIsUp()
        
    }
}

func createNewNode(address string) *Node{
    fmt.Println("Enter createNewNode")
	newAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println("ERROR: %v", err)
	}
	conn, err := net.ListenUDP("udp", newAddr)
	if err != nil {
		fmt.Println("ERROR: %v", err)
    }
    return &Node{address, conn}
    //defer conn.Close()
}

func (node *Node) checkNodeIsUp() {
    fmt.Println("Hello I am a new node existing on: " + node.address)
    defer node.connection.Close()
}

/*func main() {
    done := make(chan bool)
    //fmt.Println("port: ", os.Getenv("PORT"))
    //portOwn := os.Getenv("PORTOWN")
    portSending := os.Getenv("PORTSENDING")
    i2, err1 := strconv.Atoi(portSending)
    if err1 != nil {
        go mainServer(done, i2) //Gör egen tråd
    }
    go mainServer(done)
	<- time.After(1*time.Second)
    conn, err := net.Dial("udp", "127.0.0.1:" + portSending)
    if err != nil {
        fmt.Printf("ERROR: %v", err)
        return
	}
	fmt.Printf("Send request")
    fmt.Fprintf(conn, "Hello")
    defer conn.Close()
    <-done
}*/

func mainServer(done chan bool) {
	p := make([]byte, 2048)
    addr := net.UDPAddr{
        Port: 8000,
        IP: net.ParseIP("localhost"),
    }
    ser, err := net.ListenUDP("udp", &addr)
    if err != nil {
        fmt.Println("ERROR %v\n", err)
        return
    }
    for {
        _,remoteaddr,err := ser.ReadFromUDP(p)
        if err != nil {
            fmt.Println("ERROR %v\n", err)
            return
        }
        go sendResponse(ser, remoteaddr, done)
    }
}
func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, done chan bool) {
    _,err := conn.WriteToUDP([]byte("World"), addr)
    if err != nil {
        fmt.Println("ERROR SERVER: %v", err)
    }
    fmt.Println("Request received")
    done <- true
}