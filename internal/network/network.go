package network

import (
	"bufio"
	"encoding/json"
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
}

func (n Node) StartServer() {
	log.Printf("Starting server...")
	listenTcp, err := net.Listen("tcp", n.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer listenTcp.Close()
	for {
		connection, err := listenTcp.Accept()
		if err != nil {
			log.Panic(err)
		}
		go n.handleConnection(connection)
	}
}

func (n Node) createConnection(address string) {

	tcpAddress, err := net.ResolveTCPAddr("tcp", address)

	if err != nil {
		log.Printf("Could not resolve %v", address)
	}

	log.Printf("Trying to connect to %v...", tcpAddress.String())

	connection, err := net.DialTCP("tcp", nil, tcpAddress)

	if err != nil {
		log.Printf("%v is unreachable", address)
	}

	log.Printf("Managed to connect to %v.", connection.RemoteAddr().String())
	n.openConnections[connection.RemoteAddr().String()] = connection

	time.Sleep(2 * time.Second)

	io.WriteString(connection, fmt.Sprintln(n.Port))

	go readIncomingMessage(connection)
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
		go readIncomingMessage(connection)
		return
	}

}

func readIncomingMessage(connection net.Conn) {
	buf := bufio.NewReader(connection)
	for {
		message, err := buf.ReadBytes(10)
		if err != nil {
			log.Printf("Lost connection")
			break
		}
		log.Printf("Received message from %v", connection.LocalAddr().String())
		var block blockchain.Block

		if err := json.Unmarshal(message, &block); err != nil {
			log.Printf("Could not unmarshal received bytes to block %v", err)
			continue
		}
		log.Printf("Received block : %v", block)
	}
}

func (n Node) SendBlock() {
	block := blockchain.Block{1, "Message", "from", "server", "", ""}
	for connectionAddress, connection := range n.openConnections {
		log.Printf("Sending block to : %v...", connectionAddress)
		jsonEncoder := json.NewEncoder(connection)
		jsonEncoder.Encode(block)
	}
}

func (n Node) CloseConnections() {
	for connectionAddress, connection := range n.openConnections {
		log.Printf("Closing connection %v", connectionAddress)
		connection.Close()
	}
}

func (n Node) ListConnections() {
	i := 1
	for connectionAddress := range n.openConnections {
		log.Printf("Connection %v : %v", i, connectionAddress)
		i++
	}

}
