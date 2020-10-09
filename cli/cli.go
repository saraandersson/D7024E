package cli

import (
  "bufio"
  "fmt"
  "os"
  "strings"
  "d7024e"
  //"net"
)

func Cli(kademlia d7024e.Kademlia) {

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
			sendData := []byte(data)
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