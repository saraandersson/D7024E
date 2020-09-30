package cli

import (
  "bufio"
  "fmt"
  "os"
  "strings"
  "d7024e"
)

func CliInput(kademlia *d7024e.Kademlia, done chan bool) {

  reader := bufio.NewReader(os.Stdin)
  fmt.Println("Type operation below, you can choose between following: store, find, put, get, exit")
  fmt.Println("---------------------")

  //for {
    fmt.Print("-> ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimRight(text, "\n")
	
	switch text {
		case "store":
			fmt.Println("enter store")
			fmt.Print("Enter data to store: ")
			data, _ := reader.ReadString('\n')
			data = strings.TrimRight(text, "\n")
			sendData := []byte(data + "\n")
			kademlia.Store(sendData)
			done <- true
		case "find":
			fmt.Println("enter find")
		case "put":
			fmt.Println("enter put")
		case "get":
			fmt.Println("enter get")
		case "exit":
			fmt.Println("enter exit")
		default:
			fmt.Println("Please type correct operation, you can choose between following: store, find, put, get, exit")
	}

  //}

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