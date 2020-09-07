package main

import (
	"fmt"
    "net"
    "time"
)

func main() {
    done := make(chan bool)
    //fmt.Println("port: ", os.Getenv("PORT"))
    portOwn := os.Getenv("PORTOWN")
    portSending := os.Getenv("PORTSENDING")
	go mainServer(done, portOwn) //Gör egen tråd
	<- time.After(1*time.Second)
    conn, err := net.Dial("udp", "127.0.0.1:" + portSending)
    if err != nil {
        fmt.Printf("ERROR: %v", err)
        return
	}
	fmt.Printf("Hit kommer jag!")
    fmt.Fprintf(conn, "Hello")
    defer conn.Close()
    <-done
}

func mainServer(done chan bool, port int) {
	p := make([]byte, 2048)
    addr := net.UDPAddr{
        Port: port,
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