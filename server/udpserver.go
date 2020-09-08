package main

import (
        "fmt"
        "net"
        "os"
        "bufio"
)

func main() {
        port_input := os.Getenv("PORT")
        PORT := ":" + port_input
        s, err := net.ResolveUDPAddr("udp4", PORT)
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
                fmt.Print("Message: ", string(buffer[0:n-1]))
                reader := bufio.NewReader(os.Stdin)
                fmt.Print("Type answer here: ")
                text, _ := reader.ReadString('\n')
                data := []byte(text + "\n")
                _, err = connection.WriteToUDP(data, addr)
                if err != nil {
                        fmt.Println(err)
                        return
                }
        }
}
    