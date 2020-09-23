package main

import (
  "bufio"
  "fmt"
  "os"
  "strings"
)

func main() {

  reader := bufio.NewReader(os.Stdin)
  fmt.Println("Type operation below, you can choose between following: store, find, put, get, exit")
  fmt.Println("---------------------")

  for {
    fmt.Print("-> ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimRight(text, "\n")
	fmt.Println(text)
	// convert CRLF to LF
	
	switch text {
		case "store":
			fmt.Println("enter store")
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

  }

}