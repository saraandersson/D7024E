package main

import (
        "fmt"
        "math/rand"
        "net"
        "os"
        "strconv"
        "strings"
        "time"
        "bufio"
)

func random(min, max int) int {
        return rand.Intn(max-min) + min
}

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
        rand.Seed(time.Now().Unix())

        for {
                n, addr, err := connection.ReadFromUDP(buffer)
                fmt.Print("-> ", string(buffer[0:n-1]))

                if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
                        fmt.Println("Exiting UDP server!")
                        return
                }
                reader := bufio.NewReader(os.Stdin)
                fmt.Print(">> ")
                text, _ := reader.ReadString('\n')
                data := []byte(text + "\n")
                //data := []byte(strconv.Itoa(random(1, 1001)))
                fmt.Printf("data: %s\n", string(data))
                _, err = connection.WriteToUDP(data, addr)
                if err != nil {
                        fmt.Println(err)
                        return
                }
        }
}
    