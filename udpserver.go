package main
import (
    "net"
)

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:1234")
	if err != nil {
		// handle error
	}
	fmt.Fprintf(conn, "Hello")
	status, err := bufio.NewReader(conn).ReadString('\n')
}