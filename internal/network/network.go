package network

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kudoChain/internal/blockchain"
	"log"
	"net"
	"strings"
	"time"
)

type Node struct {
	Port            string
	openConnections map[string]net.Conn
	chain           *blockchain.Block
}

func NewNode(port string, chain *blockchain.Block) *Node {
	openConnections := make(map[string]net.Conn)
	return &Node{port, openConnections, chain}
}

func (n Node) StartServer() {
	log.Printf("Starting local node on port %v...", n.Port)
	listenTcp, err := net.Listen("tcp", n.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer listenTcp.Close()
	log.Printf("Local node successfully started. Listening connections.")
	for {
		connection, err := listenTcp.Accept()
		if err != nil {
			log.Panic(err)
		}
		go n.handleConnection(connection)
	}
}

func (n Node) CreateConnection(args []string) error {
	address := args[0]
	tcpAddress, err := net.ResolveTCPAddr("tcp", address)

	if err != nil {
		log.Printf("Could not resolve %v", address)
		return err
	}

	log.Printf("Trying to connect to %v...", tcpAddress.String())

	connection, err := net.DialTCP("tcp", nil, tcpAddress)

	if err != nil {
		log.Printf("%v is unreachable", address)
		return err
	}
	n.openConnections[connection.RemoteAddr().String()] = connection

	log.Printf("Managed to connect to %v.", connection.RemoteAddr().String())

	time.Sleep(2 * time.Second)

	io.WriteString(connection, fmt.Sprintln(n.Port))

	go n.receiveMessage(connection)
	return nil
}

func (n Node) handleConnection(connection net.Conn) {
	connectionAddress := connection.RemoteAddr().String()
	log.Printf("Handling incoming connection from %v", connectionAddress)

	bufReader := bufio.NewReader(connection)
	for {
		port, err := bufReader.ReadString('\n')
		if err != nil {
			log.Print(err)
			return
		}
		port = strings.TrimRight(port, "\n")
		n.openConnections[strings.Split(connectionAddress, ":")[0]+":"+port] = connection
		go n.receiveMessage(connection)
		return
	}

}

func (n Node) receiveMessage(connection net.Conn) {
	buf := bufio.NewReader(connection)
	for {
		message, err := buf.ReadBytes(10)
		if err != nil {
			log.Printf("Lost connection")
			break
		}
		log.Printf("Received message from %v", connection.LocalAddr().String())
		receivedChain, err := receiveChain(message)
		if err != nil {
			log.Printf("%v", err)
			continue
		}
		log.Printf("Received message is of type Block")
		*n.chain = receivedChain // should check longer chain
		n.chain.Print([]string{})
	}
}

func receiveChain(message []byte) (blockchain.Block, error) {
	var receivedChain blockchain.Block
	if err := json.Unmarshal(message, &receivedChain); err != nil {
		return receivedChain, errors.New("Received message is not of type Block")
	}
	return receivedChain, nil
}

func (n Node) SendBlock(args []string) error {
	for connectionAddress, connection := range n.openConnections {
		jsonEncoder := json.NewEncoder(connection)
		err := jsonEncoder.Encode(n.chain)
		if err != nil {
			return err
		}
		log.Printf("Sent block to : %v.", connectionAddress)
	}
	return nil
}

func (n Node) CloseConnections() {
	for connectionAddress, connection := range n.openConnections {
		log.Printf("Closing connection %v", connectionAddress)
		connection.Close()
	}
}

func (n Node) ListConnections(args []string) error {
	i := 1
	for connectionAddress := range n.openConnections {
		log.Printf("Connection %v : %v", i, connectionAddress)
		i++
	}
	return nil
}
