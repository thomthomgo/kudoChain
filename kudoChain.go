package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//TODO :
//choose longer chain

type Block struct {
	id            int
	previousHash  string
	timestamp     string
	previousOwner string
	newOwner      string
	hash          string
}

func (block Block) ToString() string {
	return strconv.Itoa(block.id) + " " + block.previousHash + " " + block.timestamp +
		" " + block.previousOwner + " " + block.newOwner
}

func (block *Block) computeHash() {
	hashString := strings.Replace(block.ToString(), " ", "", -1)
	hash := sha256.New()
	hash.Write([]byte(hashString))
	sha := hash.Sum(nil)
	block.hash = hex.EncodeToString(sha)
}

func checkBlocks(previousBlock, newBlock Block) bool {
	if newBlock.previousHash != previousBlock.hash {
		return false
	}
	return true
}

func newBlock(previousBlock Block, previousOwner, newOwner string) Block {
	newBlock := Block{previousBlock.id + 1, previousBlock.hash, time.Now().String(), previousOwner, newOwner, ""}
	newBlock.computeHash()
	return newBlock
}

func main() {
	genesis := Block{1, "", time.Now().String(), "", "Thomas", ""}
	genesis.computeHash()
	fmt.Printf(genesis.hash)
}
