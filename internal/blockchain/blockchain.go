package blockchain

import (
	"fmt"
	"sync"
	"time"

	"adaptive-blockchain/internal/consensus"
	"adaptive-blockchain/internal/merkle"
	"adaptive-blockchain/internal/shard"
)

type Blockchain struct {
	MerkleForest    *merkle.AdaptiveMerkleForest
	ConsensusEngine *consensus.HybridConsensus
	mu              sync.Mutex
	blocks          int
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		MerkleForest:    merkle.NewAdaptiveMerkleForest(),
		ConsensusEngine: consensus.NewHybridConsensus(),
		blocks:          0,
	}
}

func (bc *Blockchain) Start() {
	for {
		time.Sleep(3 * time.Second) // Simulate block production

		bc.mu.Lock()

		// Simulate new transaction insert into random shard
		shardID := bc.MerkleForest.RandomShardID()
		if shardID == "" {
			shardID = "shard-0"
			bc.MerkleForest.AddShard(shardID)
		}

		bc.MerkleForest.AddTransaction(shardID, fmt.Sprintf("tx-%d", bc.blocks))

		// Check and rebalance shards if needed
		shard.CheckRebalance(bc.MerkleForest)

		bc.blocks++
		bc.mu.Unlock()

		fmt.Printf("⛓️  New block mined: %d | Total Shards: %d\n", bc.blocks, bc.ShardCount())
	}
}

func (bc *Blockchain) BlockHeight() int {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	return bc.blocks
}

func (bc *Blockchain) ShardCount() int {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	return len(bc.MerkleForest.Shards)
}
