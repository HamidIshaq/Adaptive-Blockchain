package merkle

import (
	"fmt"
	"math/rand"
)

type MerkleTree struct {
	Transactions []string
	Filter       *BloomFilter
}

func NewMerkleTree() *MerkleTree {
	return &MerkleTree{
		Transactions: make([]string, 0),
		Filter:       NewBloomFilter(1024, 3),
	}
}

type AdaptiveMerkleForest struct {
	Shards map[string]*MerkleTree
}

func NewAdaptiveMerkleForest() *AdaptiveMerkleForest {
	return &AdaptiveMerkleForest{
		Shards: make(map[string]*MerkleTree),
	}
}

func (amf *AdaptiveMerkleForest) AddShard(id string) {
	amf.Shards[id] = NewMerkleTree()
	fmt.Printf("🌱 New shard created: %s\n", id)
}

func (amf *AdaptiveMerkleForest) AddTransaction(shardID, tx string) {
	tree, exists := amf.Shards[shardID]
	if !exists {
		tree = NewMerkleTree()
		amf.Shards[shardID] = tree
	}
	tree.Transactions = append(tree.Transactions, tx)
	tree.Filter.Add(tx)
}

func (amf *AdaptiveMerkleForest) RandomShardID() string {
	if len(amf.Shards) == 0 {
		return ""
	}
	keys := make([]string, 0, len(amf.Shards))
	for k := range amf.Shards {
		keys = append(keys, k)
	}
	return keys[rand.Intn(len(keys))]
}
