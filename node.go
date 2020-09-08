package d7024e

import (
	"fmt"
	"net"
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
        newNode := createNewNode("localhost:" + strconv.Atoi(port))
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
