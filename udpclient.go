package main

import (
        "bufio"
        "fmt"
        "net"
        "os"
        "time"
        "strconv"
)

func main() {
        //address := os.Getenv("ADDRESS")
        //port := os.Getenv("PORT")
        go mainServer(8000) 
        <- time.After(10*time.Second)
        server, err := net.ResolveUDPAddr("udp", "172.0.0.1:8000")
        conn, err := net.DialUDP("udp", nil, server)
        if err != nil {
                fmt.Println(err)
                return
        }
        fmt.Printf("The UDP server is %s\n", conn.RemoteAddr().String())
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

func mainServer(port int) {
    //port_input := os.Getenv("PORT")
    //PORT := ":" + port_input
    s, err := net.ResolveUDPAddr("udp", strconv.Itoa(port))
    if err != nil {
            fmt.Println(err)
            return
    }
    connection, err := net.ListenUDP("udp", s)
    if err != nil {
            fmt.Println(err)
            return
    }
    defer connection.Close()
    buffer := make([]byte, 1024)

    for {
            n, addr, err := connection.ReadFromUDP(buffer)
            fmt.Print("Message: ", string(buffer[0:n-1]))
           // reader := bufio.NewReader(os.Stdin)
           // fmt.Print("Type answer here: ")
           // text, _ := reader.ReadString('\n')
            data := []byte(" world " + "\n")
            _, err = connection.WriteToUDP(data, addr)
            if err != nil {
                    fmt.Println(err)
                    return
            }
    }
}
      