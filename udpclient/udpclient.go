package main

import (
        "bufio"
        "fmt"
        "net"
        "os"
        "time"
        //"strconv"
)

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
      