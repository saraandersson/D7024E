package main
import ( 
    "net" 
    "fmt" 
)

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr) {
    _,err := conn.WriteToUDP([]byte("World"), addr)
    if err != nil {
        fmt.Printf("Couldn't send response %v", err)
    }
}

func main() {
	/*p := make([]byte, 2048)
    addr := net.UDPAddr{
        Port: 1234,
        IP: net.ParseIP("127.0.0.1"),
    }*/
    err := net.Listen("udp", "127.0.0.1:1234")
    if err != nil {
        fmt.Printf("Some error1 %v\n", err)
        return
    } else {
        fmt.Printf("Ok!")
        return
    }
    /*for {
        _,remoteaddr,err := ser.ReadFromUDP(p)
        fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
        if err !=  nil {
            fmt.Printf("Some error2  %v", err)
            continue
        }
        go sendResponse(ser, remoteaddr)
    }*/
}