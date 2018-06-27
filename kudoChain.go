package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

//In progress: fix bug local address
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
)

func main() {
	serverPort := os.Args[1]
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

	log.Printf("Trying to connect to %v", tcpAddress.String())
	connection, err := net.DialTCP("tcp", nil, tcpAddress)
	if err != nil {
		log.Printf("%v is unreachable", address)
		return nil
	}
	return connection
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

func handleConnection(connection net.Conn) {
	connectionAddress := connection.RemoteAddr().String()
	log.Printf("Handling incoming connection from %v", connectionAddress)
	//Add connection to list
	openConnections[connectionAddress] = connection
	//Just some POKing code -> marshal and send block to client
	block1 := Block{1, "Message", "from", "server", "", ""}
	jsonMsg, err := json.Marshal(block1)
	if err != nil {
		log.Printf("Could not marshal block")
	}
	jsonEncoder := json.NewEncoder(connection)
	jsonEncoder.Encode(jsonMsg)
}

func closeConnections() {
	for connectionAddress, connection := range openConnections {
		log.Printf("Closing connection %v", connectionAddress)
		connection.Close()
	}

}
