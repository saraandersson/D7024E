package main

import (
	"fmt"
    "net"
    //"time"
    "os"
    "strconv"
)

type Node struct {
    address string
    connection *net.UDPConn
}


func main() {
    numberOfNodes := 4;
    startingPort := 8000;
    newPort := 0;
    port := os.Getenv("PORT")
    for i := 0; i < numberOfNodes; i++ {
        fmt.Println("Enter for-loop")
        newPort = startingPort + i;
        newNode := createNewNode("localhost:" + strconv.Itoa(newPort))
        go newNode.checkNodeIsUp() 
    }
    go mainServer(strconv.Itoa(port))
    conn, err := net.Dial("udp", "127.0.0.1:" + port)
    if err != nil {
        fmt.Printf("ERROR: %v", err)
        return
	}
	fmt.Printf("Send request")
    fmt.Fprintf(conn, "Hello")
    defer conn.Close()
}

func mainServer(port int) {
	p := make([]byte, 2048)
    addr := net.UDPAddr{
        Port: port,
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
        go sendResponse(ser, remoteaddr)
    }
}
func sendResponse(conn *net.UDPConn, addr *net.UDPAddr) {
    _,err := conn.WriteToUDP([]byte("World"), addr)
    if err != nil {
        fmt.Println("ERROR SERVER: %v", err)
    }
    fmt.Println("Request received")
}

/*

func main() {
    numberOfNodes := 4;
    startingPort := 8000;
    newPort := 0;
    port := os.Getenv("PORT")
    for i := 0; i < numberOfNodes; i++ {
        fmt.Println("Enter for-loop")
        newPort = startingPort + i;
        newNode := createNewNode("localhost:" + strconv.Itoa(newPort))
        go newNode.checkNodeIsUp() 
    }
    conn, err := net.Dial("udp", "localhost:" + strconv.Itoa(port))
    if err != nil {
        fmt.Printf("ERROR: %v", err)
        return
    }
    fmt.Fprintf(conn, "Hello")


}
*/

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
    //defer node.connection.Close()
}
