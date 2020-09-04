package main
import (
    "fmt"
    "net"
    "time"
)

func main() {
    go mainServer() //Gör egen tråd
    <- time.After(5*time.Second)
    conn, err := net.Dial("udp", "127.0.0.1:1234")
    if err != nil {
        fmt.Printf("ERROR: %v", err)
        return
    }
    fmt.Printf("Hit kommer jag!")
    fmt.Fprintf(conn, "Hello")
    conn.Close()
}

func mainServer() {
    p := make([]byte, 2048)
    addr := net.UDPAddr{
        Port: 1234,
        IP: net.ParseIP("127.0.0.1"),
    }
    ser, err := net.ListenUDP("udp", &addr)
    if err != nil {
        fmt.Printf("ERROR %v\n", err)
        return
    }
    for {
        _,remoteaddr,err := ser.ReadFromUDP(p)
        fmt.Print(p)
        if err != nil {
            fmt.Printf("ERROR %v\n", err)
            return
        }
        go sendResponse(ser, remoteaddr)
    }
}
func sendResponse(conn *net.UDPConn, addr *net.UDPAddr) {
    _,err := conn.WriteToUDP([]byte("World"), addr)
    if err != nil {
        fmt.Printf("ERROR SERVER: %v", err)
    }
    fmt.Printf("Response sending")
}