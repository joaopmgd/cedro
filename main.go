package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	errorIndex2NotFound = fmt.Errorf("index 2 not found in message")
)

func main() {
	// Replace with actual host, port, username, and password
	host := "datafeed1.cedrotech.com"
	port := "81"
	username := "username"
	password := "password"
	tickerMessage := "sqt petr4"

	// Connect to the socket
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to socket: %v", err))
	}
	defer conn.Close()

	fmt.Println("Connected to the socket.")
	reader := bufio.NewReader(conn)

	// Main loop to handle connection flow
	for {
		message := readNonIgnoredMessage(reader)

		// Connection flow and message handling
		switch {
		case message == "Connecting...":
			// Respond with an empty message
			sendMessage(conn, "")

		case message == "Welcome to Cedro Crystal":
			// Ignore this message

		case message == "Username:":
			// Respond with the username
			sendMessage(conn, username)

		case message == "Password:":
			// Respond with the password
			sendMessage(conn, password)

		case message == "You are connected":
			// Respond with the listening command
			sendMessage(conn, tickerMessage)
			fmt.Println("Authentication successful. Listening for messages...")

		default:
			// Handle and format index-value messages
			if strings.HasPrefix(message, "T:WDOF25:") {
				ticker, parsedTime, value, err := parseMessage(message)
				if err == nil {
					fmt.Printf("%s (%s): %f\n", ticker, parsedTime.Format("15:04:05"), value)
				} else {
					if err == errorIndex2NotFound {
						continue
					}
					fmt.Printf("Skipping message: %s (Error: %v)\n", message, err)
				}
			}
		}
	}
}

// sendMessage sends a message to the socket connection
func sendMessage(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message + "\n"))
	if err != nil {
		panic(fmt.Sprintf("Failed to send message: %v", err))
	}
}

// readNonIgnoredMessage reads a message from the connection and ignores blank lines, "SYN", or newline-only messages
func readNonIgnoredMessage(reader *bufio.Reader) string {
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(fmt.Sprintf("Error reading from socket: %v", err))
		}
		line = strings.TrimSpace(line)
		if line == "" || line == "SYN" {
			// Ignore empty messages, break lines, and SYN
			continue
		}
		return line
	}
}

// parseMessage parses the incoming message and extracts the ticker, formatted time, and value for index 2.
func parseMessage(message string) (string, time.Time, float64, error) {
	// Messages are separated by ':', and we are looking for "2:value" pairs.
	parts := strings.Split(message, ":")
	if len(parts) < 4 {
		return "", time.Time{}, 0, fmt.Errorf("invalid message format")
	}

	// Extract the ticker and time
	ticker := parts[1]
	rawTime := parts[2]

	// Parse the time into time.Time format (assuming BRT time zone)
	parsedTime, err := time.Parse("150405", rawTime)
	if err != nil {
		return "", time.Time{}, 0, fmt.Errorf("invalid time format: %v", err)
	}

	// Search for index 2 and extract its value
	for i := 3; i < len(parts)-1; i += 2 {
		index := parts[i]
		if index == "2" && i+1 < len(parts) {
			valueStr := parts[i+1]
			value, err := strconv.ParseFloat(valueStr, 64)
			if err != nil {
				return "", time.Time{}, 0, fmt.Errorf("invalid value for index 2: %v", err)
			}
			return ticker, parsedTime, value / 1000, nil
		}
	}

	return "", time.Time{}, 0, errorIndex2NotFound
}
