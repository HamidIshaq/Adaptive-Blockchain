package shard

import (
	"adaptive-blockchain/internal/merkle"
	"fmt"
)

const MaxTransactionsPerShard = 5

func CheckRebalance(amf *merkle.AdaptiveMerkleForest) {
	for shardID, tree := range amf.Shards {
		if len(tree.Transactions) > MaxTransactionsPerShard {
			// Split
			fmt.Printf("🔀 Rebalancing shard %s...\n", shardID)

			mid := len(tree.Transactions) / 2
			leftTxs := tree.Transactions[:mid]
			rightTxs := tree.Transactions[mid:]

			newShardLeft := shardID + "-a"
			newShardRight := shardID + "-b"

			amf.AddShard(newShardLeft)
			amf.AddShard(newShardRight)

			for _, tx := range leftTxs {
				amf.AddTransaction(newShardLeft, tx)
			}
			for _, tx := range rightTxs {
				amf.AddTransaction(newShardRight, tx)
			}

			delete(amf.Shards, shardID)
			fmt.Printf("✅ Shard %s split into %s and %s\n", shardID, newShardLeft, newShardRight)
		}
	}
}
