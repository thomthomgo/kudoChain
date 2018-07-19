package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Id            int
	PreviousHash  string
	Timestamp     string
	PreviousOwner string
	NewOwner      string
	Hash          string
	PreviousBlock *Block
}

func NewBlock(previousBlock Block, PreviousOwner, NewOwner string) Block {
	newBlock := Block{previousBlock.Id + 1, previousBlock.Hash, time.Now().Format("2006-01-02 15:04:05"), PreviousOwner, NewOwner, "", &previousBlock}
	newBlock.ComputeHash()
	return newBlock
}

func (chain *Block) AddBlock(args []string) {
	//Should add checks to ensure owner can share kudo
	previousOwner := args[0]
	newOwner := args[1]
	newBlock := NewBlock(*chain, previousOwner, newOwner)
	*chain = newBlock
}

func (block Block) ToString() string {
	return strconv.Itoa(block.Id) + " " + block.PreviousHash + " " + block.Timestamp +
		" " + block.PreviousOwner + " " + block.NewOwner
}

func (block *Block) ComputeHash() {
	HashString := strings.Replace(block.ToString(), " ", "", -1)
	Hash := sha256.New()
	Hash.Write([]byte(HashString))
	sha := Hash.Sum(nil)
	block.Hash = hex.EncodeToString(sha)
}

func CheckBlocks(previousBlock, newBlock Block) bool {
	if newBlock.PreviousHash != previousBlock.Hash {
		return false
	}
	return true
}

func ChooseLongerChain(chain1, chain2 []Block) []Block {
	if len(chain2) > len(chain1) {
		return chain2
	}
	return chain1
}

func (chain *Block) Print(args []string) {
	log.Printf("%v", chain)
	json, err := json.MarshalIndent(&chain, "", "		")
	if err != nil {
		log.Printf("Could not marshal chain : %v", err)
		return
	}
	log.Printf("%s", json)
}
