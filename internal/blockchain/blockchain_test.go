package blockchain

import (
	"testing"
	"time"
)

type hashTestPair struct {
	block        Block
	expectedHash string
}

var tests = []hashTestPair{
	{Block{1, "", "timestamp", "", "Thomas", "", nil}, "df9f931e333586ff29d2e7e1700db2728f06158d0863e6445ca610c8d5f96e30"}}

func TestComputeHash(t *testing.T) {
	for _, pair := range tests {
		pair.block.ComputeHash()
		if pair.block.Hash != pair.expectedHash {
			t.Error("For : ", pair.block.ToString(),
				"\n , expected :", pair.expectedHash,
				"\n, but was :", pair.block.Hash)
		}
	}
}

func TestHashUnicity(t *testing.T) {
	block1 := Block{1, "", "timestamp", "", "Thomas", "", nil}
	block2 := Block{1, "", "timestamp", "", "Thomas", "", nil}
	block1.ComputeHash()
	block2.ComputeHash()
	if block1.Hash != block2.Hash {
		t.Error("Block 1 : ", block1.ToString(), "Block 2 : ", block2.ToString(), "Got Different Hashes")
	}
}

func TestBlockCreation(t *testing.T) {
	block1 := Block{1, "", "Timestamp", "", "Thomas", "", nil}
	block2 := NewBlock(block1, "Thomas", "Bill")

	if !CheckBlocks(block1, block2) {
		t.Error("Block 1 : ", block1.ToString(), "Block 2 : ", block2.ToString(), "Are not well created")
	}
}

func TestCheckBlock(t *testing.T) {
	block1 := Block{1, "", "Timestamp", "", "Thomas", "", nil}
	block1.ComputeHash()
	block2 := NewBlock(block1, "Thomas", "Bill")
	block3 := Block{3, "", "Timestamp2", "", "Joe", "", nil}
	block3.ComputeHash()

	if !CheckBlocks(block1, block2) {
		t.Error("Block 1 : ", block1.ToString(), "Block 2 : ", block2.ToString(), "are not valid")
	}

	if CheckBlocks(block1, block3) {
		t.Error("Block 1 : ", block1.ToString(), "Block 3 : ", block3.ToString(), "are considered valid but should not be")
	}
}

func TestChooseLongerChain(t *testing.T) {
	block1 := Block{1, "", time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC).String(), "", "Thomas", "", nil}
	block1.ComputeHash()
	block2 := NewBlock(block1, "Thomas", "Joe")
	block3 := NewBlock(block2, "Joe", "Bill")

	chain1 := [2]Block{block1, block2}
	chain2 := [3]Block{block1, block2, block3}

	longerChain := ChooseLongerChain(chain1[:], chain2[:])

	for i := range chain2 {
		if chain2[i] != longerChain[i] {
			t.Error("Chain 2 should be longer than chain 1 : ", len(chain1), " blocks in chain 1 vs ", len(chain2), "in chain 2")
		}
	}
}
