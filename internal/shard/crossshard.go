package shard

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"adaptive-blockchain/internal/merkle"
)

type CrossShardMessage struct {
	FromShard     string
	ToShard       string
	TransactionID string
	Payload       string
	HashCommit    string // hash(Payload + TransactionID)
}

func NewCrossShardMessage(from, to, txID, payload string) *CrossShardMessage {
	commit := sha256.Sum256([]byte(payload + txID))
	return &CrossShardMessage{
		FromShard:     from,
		ToShard:       to,
		TransactionID: txID,
		Payload:       payload,
		HashCommit:    hex.EncodeToString(commit[:]),
	}
}

// ValidateCommitment checks cryptographic consistency before applying
func (msg *CrossShardMessage) ValidateCommitment() bool {
	expected := sha256.Sum256([]byte(msg.Payload + msg.TransactionID))
	return hex.EncodeToString(expected[:]) == msg.HashCommit
}

// ApplyCrossShardTransaction moves data between shards if commitment is valid
func ApplyCrossShardTransaction(amf *merkle.AdaptiveMerkleForest, msg *CrossShardMessage) error {
	if !msg.ValidateCommitment() {
		return errors.New("❌ invalid cryptographic commitment")
	}

	fromTree, ok1 := amf.Shards[msg.FromShard]
	toTree, ok2 := amf.Shards[msg.ToShard]

	if !ok1 || !ok2 {
		return errors.New("❌ source or target shard not found")
	}

	// Logically remove from source (optional, depends on design)
	// Only append to destination shard
	toTree.Transactions = append(toTree.Transactions, msg.Payload)

	if idx, ok := amf.ShardIndices[msg.ToShard]; ok {
		idx.Transactions = append(idx.Transactions, msg.Payload)
	}

	fmt.Printf("🔁 Cross-shard tx %s committed from %s → %s\n", msg.TransactionID, msg.FromShard, msg.ToShard)
	return nil
}
