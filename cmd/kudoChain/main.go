package main

import (
	"kudoChain/internal/blockchain"
	"kudoChain/internal/commandManager"
	"kudoChain/internal/network"
	"net"
	"os"
	"time"
)

var (
	terminationSignal chan bool           = make(chan bool)
	openConnections   map[string]net.Conn = make(map[string]net.Conn)
	blockChain        []blockchain.Block  = make([]blockchain.Block, 0)
)

func main() {
	chain := blockchain.Block{1, "", time.Now().Format("2006-01-02 15:04:05"), "", "Thomas", "", nil}
	chain.ComputeHash()

	serverPort := os.Args[1]
	node := network.NewNode(serverPort, &chain)

	defer node.CloseConnections()
	go node.StartServer()

	commandManager := commandManager.NewCommandManager()
	commandManager.RegisterCommand("connect", node.CreateConnection)
	commandManager.RegisterCommand("listConnections", node.ListConnections)
	commandManager.RegisterCommand("sendBlock", node.SendBlock)
	commandManager.RegisterCommand("addBlock", chain.AddBlock)
	commandManager.RegisterCommand("printChain", chain.Print)

	commandManager.Start()
}
