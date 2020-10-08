package cli

import (
  "bufio"
  "fmt"
  "os"
  "strings"
  "d7024e"
  "net"
  //"time"
  //"strconv"
)

func Cli(kademlia d7024e.Kademlia) {

	/*address := GetIPContainer() + ":" + "8080"
	contact := d7024e.NewContact(d7024e.NewKademliaID("1111111100000000000000000000000000000000"), address)
	routingtable := d7024e.NewRoutingTable(contact)
	network := d7024e.NewNetwork(contact, *routingtable)
	kademliaNetwork := d7024e.NewKademlia(&network, &contact, 20, 3)
	returnMessage := make(chan d7024e.Message)
	go network.Listen(address, 8080, returnMessage)
	<-time.After(2*time.Second) 
	bootstrapIP := "172.19.0.2"
	bootstrapAddress := bootstrapIP +":8000"
	bootstrapContact := d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), bootstrapAddress)
	routingtable.AddContact(contact)
	routingtable.AddContact(bootstrapContact)
	donePing := make(chan bool)
	//boostrapPortPing, _ := strconv.Atoi(defaultPort)
	go d7024e.SendPingMessage(&contact,&bootstrapContact,donePing)
	<- donePing
	kademliaNetwork.LookupContact(*contact.ID)
	fmt.Println("Lookup done!")*/
	
  reader := bufio.NewReader(os.Stdin)
  fmt.Println("Type operation below, you can choose between following: put, get, exit")
  fmt.Println("---------------------")

  for {
    fmt.Print("-> ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimRight(text, "\n")
	
	switch text {
		case "put":
			fmt.Print("Enter data to put: ")
			data, _ := reader.ReadString('\n')
			data = strings.TrimRight(data, "\n")
			sendData := []byte(data + "\n")
			donePut := make(chan bool)
			go kademlia.Store(sendData, donePut)
			select {
			case <-donePut:
				fmt.Println("Put is done!")
				break
			}
		case "get":
			fmt.Print("Enter key to get data: ")
			datakey, _ := reader.ReadString('\n')
			datakey = strings.TrimRight(datakey, "\n")
			file := kademlia.LookupData(datakey)
			fmt.Println("File found! Data: ")
			fmt.Println(file)
			fmt.Println("Get is done!")
			break
		case "exit":
			os.Exit(0)
		default:
			fmt.Println("Please type correct operation, you can choose between following: put, get, exit")
	}

  }

}

func GetIPContainer() string{
	containerHostname, _ := os.Hostname()
	addrs, _ := net.LookupHost(containerHostname)
	fmt.Println("Container IP address: " + addrs[0])
	return addrs[0]
}

func ExampleScanner_lines(done chan bool) {
	fmt.Println("Kommer till exempel skannern")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) // Println will add back the final '\n'
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	done <- true
}