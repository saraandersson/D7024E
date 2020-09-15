package d7024e

import (
	"fmt"
	//"bufio"
	"net"
	//"os"
	"time"
	"strconv"
)

type Network struct {
	contact *Contact
}

func Listen(ip string, port int) {
	port2 := ":" + strconv.Itoa(port)
    s, err := net.ResolveUDPAddr("udp4", port2)
    if err != nil {
            fmt.Println(err)
            return
    }
    connection, err := net.ListenUDP("udp4", s)
    if err != nil {
            fmt.Println(err)
            return
    }
    defer connection.Close()
    buffer := make([]byte, 1024)

    for {
            n, addr, err := connection.ReadFromUDP(buffer)
            data := []byte("Hello from " + ip + "\n")
            _, err = connection.WriteToUDP(data, addr)
            if err != nil {
                    fmt.Println(err)
                    return
            }
    }
}

func (network *Network) SendPingMessage(contact *Contact) {
	fmt.Println("Kommer till ping!")
	go Listen(contact.Address, 1234) //Gör egen tråd
	<- time.After(1*time.Second)
	server, err := net.ResolveUDPAddr("udp4", contact.Address)
	conn, err := net.DialUDP("udp4", nil, server)
	if err != nil {
			fmt.Println(err)
			return
	}
	fmt.Printf("The UDP server is %s\n", conn.RemoteAddr().String())
	defer conn.Close()

	for {
			/*reader := bufio.NewReader(os.Stdin)
			fmt.Printf("Type message here: ")
			message, _ := reader.ReadString('\n')*/
			data := []byte("Ping " + "\n")
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
			fmt.Printf("Answer: %s", string(buffer[0:n]))
	}
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func CreateNetwork(contact *Contact) Network {
	network := Network{}
	network.contact=contact
	return network
}
