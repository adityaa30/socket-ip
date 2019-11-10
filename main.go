package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

// Client holds information about the socket connection and data to be sent.
type Client struct {
	socket net.Conn
	data   chan []byte
}

// ClientManager structure will hold all of the available clients, received data,
// and potential incoming or terminatng clients.
type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// start() goroutine will run for lifespan of the server.
//
// reads data from the `register` channel, the connection will be stored and
// a status will be printed in the logs.
//
// If the `unregistered` channel has data and that data which represents a
// a connection, exists in our managed clients map, then the data channel for
// that connection will be closed and the connection will be removed from the
// list.
//
// If the `broadcast` channel has data it means that we've received a message
// this message should be sent to every other connection we're watching so this
// is done by looping through the available connections.
// If we can't send a message to a client, the client is closed and
// removed from the list of managed clients.
func (manager *ClientManager) start() {
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true

			address := connection.socket.RemoteAddr().String()
			fmt.Println("Added new connection", address)
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				close(connection.data)
				delete(manager.clients, connection)

				address := connection.socket.RemoteAddr().String()
				fmt.Println("A connection is terminated", address)
			}
		case message := <-manager.broadcast:
			for connection := range manager.clients {
				select {
				case connection.data <- message:
				default:
					close(connection.data)
					delete(manager.clients, connection)
				}
			}
		}
	}
}

// receive() is a goroutine and will exist for every connection that is
// established
//
// For aslong as the goroutine is available, it will be receiving data from a
// particular client.
//
// If everything went well and the message received wasn't empty
func (manager *ClientManager) receive(client *Client) {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)

		if err != nil {
			// the client will be unregistered & closed
			manager.unregister <- client
			client.socket.Close()
			break
		}

		// If message received is not empty, it will be added to broadcast channel
		// to be distributed to all client by the manager.
		if length > 0 {
			fmt.Println("Received:", string(message))
			manager.broadcast <- message
		}
	}
}

// If the client has data to be send and there are no errors, that data
// will be sent to the client in question.
func (manager *ClientManager) send(client *Client) {
	defer client.socket.Close()
	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			}
			client.socket.Write(message)
		}
	}
}

func startServerMode() {
	fmt.Println("Starting server...")
	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		fmt.Println(err)
	}

	manager := ClientManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 1),
		register:   make(chan *Client, 1),
		unregister: make(chan *Client, 1),
	}

	go manager.start()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		client := &Client{
			socket: connection,
			data:   make(chan []byte),
		}
		manager.register <- client

		go manager.receive(client)
		go manager.send(client)
	}
}

func (client *Client) receive() {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)

		// If there is read error due to some issue with server,
		// the connection will close
		if err != nil {
			client.socket.Close()
			break
		}

		if length > 0 {
			fmt.Println("Received:", string(message))
		}
	}
}

// While the messages are continuosly receiving, we will be continuosly
// allowing the user input for message sending. When enter key is pressed, text
// will be sent and the cycle will continue.
func startClientMode() {
	fmt.Println("Starting client...")
	conn, err := net.Dial("tcp", "localhost:12345")

	if err != nil {
		fmt.Println(err)
	}

	client := &Client{socket: conn}
	go client.receive()

	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		conn.Write([]byte(strings.TrimRight(message, "\n")))
	}
}

func main() {
	flagMode := flag.String("mode", "server", "Start in client or server mode")
	flag.Parse()

	if strings.ToLower(*flagMode) == "server" {
		startServerMode()
	} else {
		startClientMode()
	}
}
