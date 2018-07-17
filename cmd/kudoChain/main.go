package main

import (
	"kudoChain/internal/blockchain"
	"kudoChain/internal/commandManager"
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
	node := network.NewNode(serverPort)

	defer node.CloseConnections()
	go node.StartServer()

	commandManager := commandManager.NewCommandManager()
	commandManager.RegisterCommand("connect", node.CreateConnection)
	commandManager.RegisterCommand("listConnections", node.ListConnections)
	commandManager.RegisterCommand("sendBlock", node.SendBlock)

	commandManager.UserInput()
}
