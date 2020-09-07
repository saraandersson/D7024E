package main

import (
	"fmt"
    "net"
    "time"
    "flag"
)

func main() {
    done := make(chan bool)
    numbPtr := flag.Int("numb", 42, "an int")
    flag.Parse()
    fmt.Println("numb:", *numbPtr)
	go mainServer(done) //Gör egen tråd
	<- time.After(1*time.Second)
    conn, err := net.Dial("udp", "127.0.0.1:1234")
    if err != nil {
        fmt.Printf("ERROR: %v", err)
        return
	}
	fmt.Printf("Hit kommer jag!")
    fmt.Fprintf(conn, "Hello")
    defer conn.Close()
    <-done
}

func mainServer(done chan bool) {
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
        if err != nil {
            fmt.Printf("ERROR %v\n", err)
            return
        }
        go sendResponse(ser, remoteaddr, done)
    }
}
func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, done chan bool) {
    _,err := conn.WriteToUDP([]byte("World"), addr)
    if err != nil {
        fmt.Printf("ERROR SERVER: %v", err)
    }
    fmt.Printf("Response sending")
    done <- true
}