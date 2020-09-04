package main
import (
	"fmt"
	"net"
	"time"
)

func main() {
	go mainServer() //Starting mainServer as a thred
	<- time.After(5*time.Second)  //Timer, continue after 5 second
    conn, err := net.Dial("udp", "127.0.0.1:1234") //Sending UDP request
    if err != nil {
        fmt.Printf("Error: Error when sending request, %v", err)
        return
	}
	fmt.Printf("Sending request.")
    fmt.Fprintf(conn, "Hello")  //Sending "Hello" as a request
    conn.Close()
}

func mainServer() {
	p := make([]byte, 2048)
    addr := net.UDPAddr{
        Port: 1234,
        IP: net.ParseIP("127.0.0.1"),
    }
    ser, err := net.ListenUDP("udp", &addr) //Listen after UDP request
    if err != nil {
        fmt.Printf("Error: Error when reciving request, %v", err)
        return
	}
    for {
		_,remoteaddr,err := ser.ReadFromUDP(p) //Reading request message
		fmt.Print(p) 
        if err != nil {
            fmt.Printf("Error: Error when reading request, %v", err)
            return
        }
        go sendResponse(ser, remoteaddr)  //Call sendResponse function as a thread
    }
}
func sendResponse(conn *net.UDPConn, addr *net.UDPAddr) {
    _,err := conn.WriteToUDP([]byte("World"), addr)
    if err != nil {
        fmt.Printf("Error: Error when sending response, %v", err)
	}
	fmt.Printf("Response sending")
}