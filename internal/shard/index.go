package shard

// ShardIndex keeps metadata about each shard
type ShardIndex struct {
	ID           string
	Transactions []string
	CreatedAt    int64
}

// NewShardIndex initializes a new shard index
func NewShardIndex(id string, txs []string, timestamp int64) *ShardIndex {
	return &ShardIndex{
		ID:           id,
		Transactions: txs,
		CreatedAt:    timestamp,
	}
}
