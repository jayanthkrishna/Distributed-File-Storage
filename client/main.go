package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Dial a connection to the server at localhost:8080
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		fmt.Println("Error dialing server:", err)
		os.Exit(1)
	}
	defer conn.Close() // Ensure the connection is closed when we're done

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter some text: ")
		// Read a full line including spaces
		message, _ := reader.ReadString('\n')
		fmt.Println(message)
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error writing to connection:", err)
			continue
		}
	}
}
