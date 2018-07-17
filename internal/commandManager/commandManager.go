package commandManager

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)

type CommandManager struct {
	commands          map[string]func(args []string)
	terminationSignal chan bool
}

func NewCommandManager() *CommandManager {
	commands := make(map[string]func(args []string))
	terminationSignal := make(chan bool)
	return &CommandManager{commands, terminationSignal}
}

func (manager CommandManager) RegisterCommand(name string, function func(args []string)) {
	manager.commands[name] = function
}

func (manager CommandManager) UserInput() {
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
			go manager.manageCommand(input)
		case <-time.After(1 * time.Second):
		case shouldTerminate := <-manager.terminationSignal:
			if shouldTerminate {
				log.Printf("Shutting down...")
				return
			}
		}
	}
}

func (manager CommandManager) manageCommand(fullCommand string) {

	fullCommand = strings.TrimRight(fullCommand, "\n")
	split := strings.Split(fullCommand, " ")
	command := split[0]
	args := split[1:]

	switch command {
	case "quit":
		manager.terminationSignal <- true
	default:
		if manager.commands[command] != nil {
			manager.commands[command](args)
		} else {
			log.Printf("Unknown command:%v", command)
		}
	}

}
