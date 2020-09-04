package main
import (
    "fmt"
    "net"
    "time"
)

func main() {
    go mainServer() //Gör egen tråd
    <- time.After(10*time.Second)
    conn, err := net.Dial("udp", "127.0.0.1:1234")
    if err != nil {
        fmt.Printf("ERROR: %v", err)
        return
    }
    fmt.Printf("Hit kommer jag!")
    fmt.Fprintf(conn, "Hello")
    conn.Close()
}

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr) {
    _,err := conn.WriteToUDP([]byte("From server: Hello I got your mesage "), addr)
    if err != nil {
        fmt.Printf("Couldn't send response %v", err)
    }
}


func mainServer() {
    p := make([]byte, 2048)
    addr := net.UDPAddr{
        Port: 1234,
        IP: net.ParseIP("127.0.0.1"),
    }
    ser, err := net.ListenUDP("udp", &addr)
    if err != nil {
        fmt.Printf("Some error %v\n", err)
        return
    }
    for {
        _,remoteaddr,err := ser.ReadFromUDP(p)
        fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
        if err !=  nil {
            fmt.Printf("Some error  %v", err)
            continue
        }
        go sendResponse(ser, remoteaddr)
    }
}