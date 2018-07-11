package commandManager

import (
	"bufio"
	"kudoChain/internal/network"
	"log"
	"os"
	"strings"
	"time"
)

func UserInput() {
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
		network.CreateConnection(strings.TrimRight(commandWithArgs[1], "\n"))
	case "listConnections":
		network.ListConnections()
	case "sendBlock":
		network.SendBlock()
	default:
		log.Printf("Unknown command:%v", command)
	}

}
