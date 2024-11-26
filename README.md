# Cedro Crystal Socket Client

This Go program connects to the Cedro Crystal system using a TCP socket, authenticates the user, and subscribes to real-time quotes for a specific asset (`petr4` by default) using the `sqt` command.

## Features

- Automatically attempts connection to multiple servers in sequence.
- Retries connections indefinitely if all servers fail.
- Authenticates using user credentials and optional Software Key.
- Subscribes to and processes real-time data for a specified asset.
- Prints incoming messages to the console in real-time.

## Prerequisites

Ensure you have the following:
- Go 1.18 or higher installed on your system.
- Cedro Crystal credentials (username, password, optional Software Key).
- A list of valid server addresses with access to Cedro Crystal.

## Installation

1. Clone this repository and navigate to the project folder:

```
   git clone https://github.com/yourusername/cedro-crystal-client.git
   cd cedro-crystal-client
```

2. Update the `main.go` file with your Cedro Crystal credentials, server list, and the asset you wish to monitor.

3. Build and run the program:

   go run main.go

## Configuration

### Server List

The program uses a list of servers for connection attempts. Update the `servers` list in `main.go` with valid server addresses:

servers := []string{
    "server1-address:81",
    "server2-address:81",
    "server3-address:81",
}

### User Credentials

Replace the placeholders for `username` and `password` in `main.go` with your Cedro Crystal credentials. Leave the `softwareKey` variable empty if not required.

### Asset Subscription

Specify the asset to monitor by updating the `ativo` variable in `main.go`. For example:

#### ativo := "petr4"


## Usage

Run the program to establish a connection, authenticate, and subscribe to the desired asset. The program will:
- Attempt connection to servers in the order provided.
- Authenticate with your Cedro Crystal credentials.
- Subscribe to the real-time quotes for the configured asset.
- Print incoming messages to the console.

## Commands

The program uses the `sqt` command to subscribe to real-time data for a specific asset.

**Syntax**:
`sqt <asset>`

Add `N` to the command for a one-time snapshot without monitoring.

**Example**:
`>sqt petr4`

## Message Format

The Cedro Crystal system sends messages in the following format:

### Header

`T:<asset>:<time>`

- `T`: Message type.
- `<asset>`: The asset being monitored.
- `<time>`: Timestamp of the message.

### Body

Consists of one or more key-value pairs:

`<index>:<value>`

Each message ends with the `!` character.

## Example Output

Attempting to connect to server: server1-address:81
Connected to server: server1-address:81
Server: You are connected
Authentication successful
Subscribed to asset: petr4
Message: T:petr4:12:30:45
Message: 1:1234.56!
Message: 2:7890!

## Error Handling

The program automatically retries connection attempts if a server fails. If all servers are unreachable, it continues retrying indefinitely with a brief delay between attempts. Errors during communication terminate the program with an appropriate error message.

## Contributing

Contributions are welcome! If you have ideas, bug reports, or improvements, feel free to open an issue or submit a pull request.