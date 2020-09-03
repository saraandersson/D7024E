package main
import (
	"fmt"
	"net"
	"bufio"
)

func main() {
    p :=  make([]byte, 2048)
    conn, err := net.Dial("udp", "127.0.0.1:1234")
    if err != nil {
        fmt.Printf("Some error3 %v", err)
        return
    }
    fmt.Fprintf(conn, "Hi UDP Server?")
    _, err = bufio.NewReader(conn).Read(p)
    if err == nil {
        fmt.Printf("%s\n", p)
    } else {
        fmt.Printf("Some error4 %v\n", err)
    }
    conn.Close()
}