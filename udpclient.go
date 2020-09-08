package main

import (
        "bufio"
        "fmt"
        "net"
        "os"
)

func main() {
        address := os.Getenv("ADDRESS")
        server, err := net.ResolveUDPAddr("udp", address)
        conn, err := net.DialUDP("udp", nil, server)
        if err != nil {
                fmt.Println(err)
                return
        }
        fmt.Printf("The UDP server is %server\n", conn.RemoteAddr().String())
        defer conn.Close()

        for {
                reader := bufio.NewReader(os.Stdin)
                fmt.Print("Type message here: ")
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
                fmt.Printf("Answer: %s\n", string(buffer[0:n]))
        }
}
      