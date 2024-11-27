package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	// Dial a connection to the server at localhost:8080
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		fmt.Println("Error dialing server:", err)
		os.Exit(1)
	}
	defer conn.Close() // Ensure the connection is closed when we're done

	// Print that the connection was successful
	fmt.Println("Connected to the server. Connection will stay open for 5 seconds...")

	// Sleep for 5 seconds to keep the connection open
	time.Sleep(5 * time.Second)

	// After 5 seconds, the connection will automatically be closed
	fmt.Println("Closing the connection after 5 seconds")
}
