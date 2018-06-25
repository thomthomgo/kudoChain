package main

import (
	"testing"
	"time"
)

type HashTestPair struct {
	block        Block
	expectedHash string
}

var tests = []HashTestPair{
	{Block{1, "", "Timestamp", "", "Thomas", ""}, "a2af678f097b53af79ffaa59b7d858531e4d0d83f65737e21ca003b28d18a4a6"}}

func TestComputeHash(t *testing.T) {
	for _, pair := range tests {
		pair.block.computeHash()
		if pair.block.Hash != pair.expectedHash {
			t.Error("For : ", pair.block.ToString(),
				"\n , expected :", pair.expectedHash,
				"\n, but was :", pair.block.Hash)
		}
	}
}

func TestHashUnicity(t *testing.T) {
	block1 := Block{1, "", "Timestamp", "", "Thomas", ""}
	block2 := Block{1, "", "Timestamp", "", "Thomas", ""}
	block1.computeHash()
	block2.computeHash()
	if block1.Hash != block2.Hash {
		t.Error("Block 1 : ", block1.ToString(), "Block 2 : ", block2.ToString(), "Got Different Hashes")
	}
}

func TestBlockCreation(t *testing.T) {
	block1 := Block{1, "", "Timestamp", "", "Thomas", ""}
	block2 := newBlock(block1, "Thomas", "Bill")

	if !checkBlocks(block1, block2) {
		t.Error("Block 1 : ", block1.ToString(), "Block 2 : ", block2.ToString(), "Are not well created")
	}
}

func TestCheckBlock(t *testing.T) {
	block1 := Block{1, "", "Timestamp", "", "Thomas", ""}
	block1.computeHash()
	block2 := newBlock(block1, "Thomas", "Bill")
	block3 := Block{3, "", "Timestamp2", "", "Joe", ""}
	block3.computeHash()

	if !checkBlocks(block1, block2) {
		t.Error("Block 1 : ", block1.ToString(), "Block 2 : ", block2.ToString(), "are not valid")
	}

	if checkBlocks(block1, block3) {
		t.Error("Block 1 : ", block1.ToString(), "Block 3 : ", block3.ToString(), "are considered valid but should not be")
	}
}

func TestChooseLongerChain(t *testing.T) {
	block1 := Block{1, "", time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC).String(), "", "Thomas", ""}
	block1.computeHash()
	block2 := newBlock(block1, "Thomas", "Joe")
	block3 := newBlock(block2, "Joe", "Bill")

	chain1 := [2]Block{block1, block2}
	chain2 := [3]Block{block1, block2, block3}

	longerChain := chooseLongerChain(chain1[:], chain2[:])

	for i := range chain2 {
		if chain2[i] != longerChain[i] {
			t.Error("Chain 2 should be longer than chain 1 : ", len(chain1), " blocks in chain 1 vs ", len(chain2), "in chain 2")
		}
	}
}
