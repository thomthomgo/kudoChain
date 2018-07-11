package main

import (
	"kudoChain/internal/blockchain"
	"kudoChain/internal/network"
	"net"
	"os"
)

//TODO : handle reception of block

var (
	terminationSignal chan bool           = make(chan bool)
	openConnections   map[string]net.Conn = make(map[string]net.Conn)
	blockChain        []blockchain.Block  = make([]blockchain.Block, 0)
)

func main() {
	serverPort := os.Args[1]
	node := network.Node{Port: serverPort}

	defer node.closeConnections()
	go network.server(serverPort)
	commandManager.userInput()
}
