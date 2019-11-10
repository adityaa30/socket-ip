package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

// Message holds the information which will be send to the Node `name`
type Message struct {
	name string
	data []byte
}

// Client holds information about the socket connection and data to be sent.
type Client struct {
	socket  net.Conn
	name    string
	message chan Message
}

// ClientManager structure will hold all of the available clients, received data,
// and potential incoming or terminatng clients.
type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
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
		case conn := <-manager.register:
			manager.clients[conn] = true

			address := conn.socket.RemoteAddr().String()
			fmt.Println("Added new connection", address, "Name: ", conn.name)

		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.message)
				delete(manager.clients, conn)

				address := conn.socket.RemoteAddr().String()
				fmt.Println("A connection is terminated", address, "Name: ", conn.name)
			}

		case message := <-manager.broadcast:
			for conn := range manager.clients {
				if conn.name == message.name {
					select {
					case conn.message <- message:
					default:
						close(conn.message)
						delete(manager.clients, conn)
					}
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

			var name, str string
			fmt.Sscanf(string(message), "%s %q", &name, &str)
			manager.broadcast <- Message{
				name: name,
				data: []byte(str),
			}
		}
	}
}

// If the client has data to be send and there are no errors, that data
// will be sent to the client in question.
func (manager *ClientManager) send(client *Client) {
	defer client.socket.Close()
	for {
		select {
		case message, ok := <-client.message:
			if !ok {
				return
			}
			client.socket.Write(message.data)
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
		broadcast:  make(chan Message),
		register:   make(chan *Client, 1),
		unregister: make(chan *Client, 1),
	}

	go manager.start()

	clientNames := [7]string{"A", "B", "C", "D", "E", "F", "G"}

	for i := 0; i < len(clientNames); i++ {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		client := &Client{
			socket:  connection,
			name:    clientNames[i],
			message: make(chan Message, 1),
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

	fmt.Println("Your ip: ", conn.LocalAddr().String())
	fmt.Println("<Category> <Node Name> \"<Message>\"")
	fmt.Println("Category:\n1 = Send file\n2 = Send text")
	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')

		var category int // 1 => Send a file; 2 => Send text
		var name, str string

		fmt.Sscanf(string(message), "%d %s %q", &category, &name, &str)
		if category == 1 {
			data, err := ioutil.ReadFile(str)
			if err != nil {
				fmt.Println("File reading error", err)
			} else {
				filePath := str
				str := string(data)
				chunkSize := 1024 / 8
				chunks := len(str) / chunkSize
				fmt.Printf("Read the file '%s' and divided into %d chunks.\n", filePath, chunks+1)
				for i := 0; i < chunks+1; i++ {
					fmt.Printf("Sending chunk %d\n", i)
					message := string(name + " \"" + str[(i*chunkSize):min((i+1)*chunkSize, len(str))] + "\"")
					conn.Write([]byte(strings.TrimRight(message, "\n")))
					time.Sleep(300 * time.Millisecond)
				}
			}
		} else {
			message := string(name + " \"" + str + "\"")
			conn.Write([]byte(strings.TrimRight(message, "\n")))
		}

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
