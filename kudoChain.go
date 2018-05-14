package main

import (
	"strconv"
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
}

func computeHash(block Block) string {
	hashString := strconv.Itoa(block.id) + block.previousHash + block.timestamp + block.previousOwner + block.newOwner
	return hashString
}

func main() {

}
