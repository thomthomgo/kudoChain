package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
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

func NewBlock(previousBlock Block, PreviousOwner, NewOwner string) Block {
	newBlock := Block{previousBlock.Id + 1, previousBlock.Hash, time.Now().String(), PreviousOwner, NewOwner, ""}
	newBlock.ComputeHash()
	return newBlock
}

func ChooseLongerChain(chain1, chain2 []Block) []Block {
	if len(chain2) > len(chain1) {
		return chain2
	}
	return chain1
}
