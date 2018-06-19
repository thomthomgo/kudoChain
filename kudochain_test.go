package main

import "testing"

type hashTestPair struct {
	block        Block
	expectedHash string
}

var tests = []hashTestPair{
	{Block{1, "", "timestamp", "", "Thomas", ""}, "df9f931e333586ff29d2e7e1700db2728f06158d0863e6445ca610c8d5f96e30"}}

func TestComputeHash(t *testing.T) {
	for _, pair := range tests {
		pair.block.computeHash()
		if pair.block.hash != pair.expectedHash {
			t.Error("For : ", pair.block.ToString(),
				"\n , expected :", pair.expectedHash,
				"\n, but was :", pair.block.hash)
		}
	}
}

func TestHashUnicity(t *testing.T) {
	block1 := Block{1, "", "timestamp", "", "Thomas", ""}
	block2 := Block{1, "", "timestamp", "", "Thomas", ""}
	block1.computeHash()
	block2.computeHash()
	if block1.hash != block2.hash {
		t.Error("Block 1 : ", block1.ToString(), "Block 2 : ", block2.ToString(), "Got Different Hashes")
	}
}

func TestBlockCreation(t *testing.T) {
	block1 := Block{1, "", "timestamp", "", "Thomas", ""}
	block2 := newBlock(block1, "Thomas", "Bill")

	if !checkBlocks(block1, block2) {
		t.Error("Block 1 : ", block1.ToString(), "Block 2 : ", block2.ToString(), "Are not well created")
	}
}

func TestCheckBlock(t *testing.T) {
	block1 := Block{1, "", "timestamp", "", "Thomas", ""}
	block1.computeHash()
	block2 := newBlock(block1, "Thomas", "Bill")
	block3 := Block{3, "", "timestamp2", "", "Joe", ""}
	block3.computeHash()

	if !checkBlocks(block1, block2) {
		t.Error("Block 1 : ", block1.ToString(), "Block 2 : ", block2.ToString(), "are not valid")
	}

	if checkBlocks(block1, block3) {
		t.Error("Block 1 : ", block1.ToString(), "Block 3 : ", block3.ToString(), "are considered valid but should not be")
	}

}
