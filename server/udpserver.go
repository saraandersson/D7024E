package main

import (
        "fmt"
        "net"
       // "os"
      //  "bufio"
)

func main() {
    numberOfNodes := 1;
    port := 8000;
    for i := 0; i < numberOfNodes; i++ {
        //fmt.Println("Enter for-loop")
        newPort := port + i
        newNode := createNewNode("localhost:" + strconv.Itoa(newNode))
        go newNode.checkNodeIsUp()
    }
}

func createNewNode(address string) *Node{
        s, err := net.ResolveUDPAddr("udp", address)
        if err != nil {
                fmt.Println(err)
                
        }
        connection, err := net.ListenUDP("udp", s)
        if err != nil {
                fmt.Println(err)
                
        }
        defer connection.Close()
        buffer := make([]byte, 1024)

        for {
                n, addr, err := connection.ReadFromUDP(buffer)
                fmt.Print("Message: ", string(buffer[0:n-1]))
               // reader := bufio.NewReader(os.Stdin)
               // fmt.Print("Type answer here: ")
               // text, _ := reader.ReadString('\n')
                data := []byte("WORLD!" + "\n")
                _, err = connection.WriteToUDP(data, addr)
                if err != nil {
                        fmt.Println(err)
                        
                }
        }
        return &Node{address, connection}
}

func (node *Node) checkNodeIsUp() {
    fmt.Println("Hello I am a new node existing on: " + node.address)
    defer node.connection.Close()
}
    