package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

//struct Block
//calculate Hash
//generateBlock (from old block)
//check Block
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
	return strconv.Itoa(block.id) + block.previousHash + block.timestamp + block.previousOwner + block.newOwner
}

func (block *Block) computeHash() {
	hashString := block.ToString()
	hash := sha256.New()
	hash.Write([]byte(hashString))
	sha := hash.Sum(nil)
	block.hash = hex.EncodeToString(sha)
}

func main() {
	genesis := Block{1, "", time.Now().String(), "", "Thomas", ""}
	genesis.computeHash()
	fmt.Printf(genesis.hash)

}
