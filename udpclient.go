package main

import (
        "bufio"
        "fmt"
        "net"
        "os"
        "strings"
)

func main() {
        con := os.Getenv("ADDRESS")
        s, err := net.ResolveUDPAddr("udp4", con)
        c, err := net.DialUDP("udp4", nil, s)
        if err != nil {
                fmt.Println(err)
                return
        }

        fmt.Printf("The UDP server is %s\n", c.RemoteAddr().String())
        defer c.Close()

        for {
                reader := bufio.NewReader(os.Stdin)
                fmt.Print(">> ")
                text, _ := reader.ReadString('\n')
                data := []byte(text + "\n")
                _, err = c.Write(data)
                if err != nil {
                        fmt.Println(err)
                        return
                }

                buffer := make([]byte, 1024)
                n, _, err := c.ReadFromUDP(buffer)
                if err != nil {
                        fmt.Println(err)
                        return
                }
                fmt.Printf("Reply: %s\n", string(buffer[0:n]))
        }
}
      