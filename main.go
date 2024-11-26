package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func main() {
	servers := []string{
		"datafeed1.cedrotech.com:81",
		"datafeed2.cedrotech.com:81",
	}

	username := "username"
	password := "password"
	ativo := "petr4"

	var conn net.Conn
	var err error

	// Try connecting to servers indefinitely
	for {
		for _, server := range servers {
			fmt.Printf("Attempting to connect to server: %s\n", server)
			conn, err = net.Dial("tcp", server)
			if err == nil {
				fmt.Printf("Connected to server: %s\n", server)
				break
			}
			fmt.Printf("Failed to connect to server: %s, error: %v\n", server, err)
		}

		if conn != nil && err == nil {
			break
		}

		fmt.Println("All servers failed. Retrying...")
		time.Sleep(2 * time.Second) // Optional sleep before retrying
	}
	defer conn.Close()

	// Initialize reader and writer for the connection
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Authenticate
	fmt.Fprintln(writer, username)
	writer.Flush()

	fmt.Fprintln(writer, password)
	writer.Flush()

	// Wait for "You are connected" confirmation
	for {
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading response: %v", err)
		}
		fmt.Print("Server: ", response)
		if strings.Contains(response, "You are connected") {
			fmt.Println("Authentication successful")
			break
		}
	}

	// Subscribe to asset using the SQT command
	sqtCommand := fmt.Sprintf("sqt %s", ativo)
	fmt.Fprintln(writer, sqtCommand)
	writer.Flush()
	fmt.Printf("Subscribed to asset: %s\n", ativo)

	// Read and process incoming messages
	go func() {
		for {
			response, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("Error reading data: %v", err)
			}
			if strings.TrimSpace(response) == "" {
				continue
			}
			fmt.Println("Message: ", response)
		}
	}()

	// Keep the connection alive
	for {
		time.Sleep(1 * time.Second)
	}
}
