package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

//TODO : handle reception of block

type Block struct {
	Id            int
	PreviousHash  string
	Timestamp     string
	PreviousOwner string
	NewOwner      string
	Hash          string
}

func (block Block) ToString() string {
	return strconv.Itoa(block.Id) + " " + block.PreviousHash + " " + block.Timestamp +
		" " + block.PreviousOwner + " " + block.NewOwner
}

func (block *Block) computeHash() {
	HashString := strings.Replace(block.ToString(), " ", "", -1)
	Hash := sha256.New()
	Hash.Write([]byte(HashString))
	sha := Hash.Sum(nil)
	block.Hash = hex.EncodeToString(sha)
}

func checkBlocks(previousBlock, newBlock Block) bool {
	if newBlock.PreviousHash != previousBlock.Hash {
		return false
	}
	return true
}

func newBlock(previousBlock Block, PreviousOwner, NewOwner string) Block {
	newBlock := Block{previousBlock.Id + 1, previousBlock.Hash, time.Now().String(), PreviousOwner, NewOwner, ""}
	newBlock.computeHash()
	return newBlock
}

func chooseLongerChain(chain1, chain2 []Block) []Block {
	if len(chain2) > len(chain1) {
		return chain2
	}
	return chain1
}

var (
	terminationSignal chan bool           = make(chan bool)
	openConnections   map[string]net.Conn = make(map[string]net.Conn)
	blockChain        []Block             = make([]Block, 0)
	serverPort        string
)

func main() {
	serverPort = os.Args[1]
	//clientPort := os.Args[2]
	defer closeConnections()
	go server(serverPort)
	userInput()
}

func userInput() {
	ch := make(chan string)

	go func(ch chan string) {
		inputReader := bufio.NewReader(os.Stdin)
		for {
			s, err := inputReader.ReadString('\n')
			if err != nil {
				close(ch)
				return
			}
			ch <- s
		}
	}(ch)

	for {
		select {
		case input := <-ch:
			go manageCommand(input)
		case <-time.After(1 * time.Second):
		case shouldTerminate := <-terminationSignal:
			if shouldTerminate {
				log.Printf("Shutting down...")
				return
			}
		}

	}
}

func manageCommand(command string) {
	commandWithArgs := strings.Split(command, " ")
	switch strings.TrimRight(commandWithArgs[0], "\n") {
	case "quit":
		terminationSignal <- true
	case "connect":
		createConnection(strings.TrimRight(commandWithArgs[1], "\n"))
	case "listConnections":
		listConnections()
	case "sendBlock":
		sendBlock()
	default:
		log.Printf("Unknown command:%v", command)
	}

}

func createConnection(address string) net.Conn {
	tcpAddress, err := net.ResolveTCPAddr("tcp", address)

	if err != nil {
		log.Printf("Could not resolve %v", address)
		//SHOULD HANDLE
	}
	log.Printf("Trying to connect to %v...", tcpAddress.String())
	connection, err := net.DialTCP("tcp", nil, tcpAddress)
	if err != nil {
		log.Printf("%v is unreachable", address)
		return nil
	}
	log.Printf("Managed to connect to %v.", connection.RemoteAddr().String())
	openConnections[connection.RemoteAddr().String()] = connection

	time.Sleep(2 * time.Second)

	io.WriteString(connection, fmt.Sprintln(serverPort))

	go readIncomingMessage(connection)
	return connection
}

func handleConnection(connection net.Conn) {
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
		openConnections[strings.Split(connectionAddress, ":")[0]+":"+port] = connection
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
		var block Block

		if err := json.Unmarshal(message, &block); err != nil {
			log.Printf("Could not unmarshal received bytes to block %v", err)
			continue
		}
		log.Printf("Received block : %v", block)
	}
}

func server(port string) {
	log.Printf("Starting server...")

	listenTcp, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	defer listenTcp.Close()
	for {
		connection, err := listenTcp.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(connection)
	}
}

func sendBlock() {
	block1 := Block{1, "Message", "from", "server", "", ""}

	for connectionAddress, connection := range openConnections {
		log.Printf("Sending block to : %v...", connectionAddress)
		jsonEncoder := json.NewEncoder(connection)
		jsonEncoder.Encode(block1)
	}
}

func closeConnections() {
	for connectionAddress, connection := range openConnections {
		log.Printf("Closing connection %v", connectionAddress)
		connection.Close()
	}
}

func listConnections() {
	i := 1
	for connectionAddress := range openConnections {
		log.Printf("Connection %v : %v", i, connectionAddress)
		i++
	}

}
