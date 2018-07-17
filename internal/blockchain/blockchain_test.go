package blockchain

import "testing"

type hashTestPair struct {
	block        Block
	expectedHash string
}

var tests = []hashTestPair{
	{Block{1, "", "timestamp", "", "Thomas", ""}, "df9f931e333586ff29d2e7e1700db2728f06158d0863e6445ca610c8d5f96e30"}}

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
	block1 := Block{1, "", "timestamp", "", "Thomas", ""}
	block2 := Block{1, "", "timestamp", "", "Thomas", ""}
	block1.ComputeHash()
	block2.ComputeHash()
	if block1.Hash != block2.Hash {
		t.Error("Block 1 : ", block1.ToString(), "Block 2 : ", block2.ToString(), "Got Different Hashes")
	}
}
