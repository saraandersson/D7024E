package main

import (
	"fmt"
    "net"
)

func main() {
	p := make([]byte, 2048)
    addr := net.UDPAddr{
        Port: 8001,
        IP: net.ParseIP("127.0.0.1"),
    }
    ser, err := net.ListenUDP("udp", &addr)
    if err != nil {
        fmt.Printf("ERROR %v\n", err)
        return
    }
    for {
        _,remoteaddr,err := ser.ReadFromUDP(p)
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