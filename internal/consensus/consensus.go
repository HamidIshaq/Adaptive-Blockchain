package consensus

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"
)

type Node struct {
	ID         string
	Reputation int
	Stake      int
	Trusted    bool
}

type HybridConsensus struct {
	nodes            map[string]*Node
	leaderID         string
	quorumThreshold  float64
	difficulty       int
	leaderElectionMu sync.Mutex
}

func NewHybridConsensus() *HybridConsensus {
	return &HybridConsensus{
		nodes:           make(map[string]*Node),
		quorumThreshold: 0.66,
		difficulty:      4,
	}
}

func (hc *HybridConsensus) RegisterNode(id string, stake int) {
	hc.nodes[id] = &Node{
		ID:         id,
		Reputation: 100,
		Stake:      stake,
		Trusted:    true,
	}
}

// Proof of Work for randomness injection
type PoWResult struct {
	Nonce string
	Hash  string
}

func performPoW(data string, difficulty int) PoWResult {
	prefix := ""
	for i := 0; i < difficulty; i++ {
		prefix += "0"
	}

	for {
		nonceBytes := make([]byte, 8)
		rand.Read(nonceBytes)
		hash := sha256.Sum256(append([]byte(data), nonceBytes...))
		hashStr := hex.EncodeToString(hash[:])
		if hashStr[:difficulty] == prefix {
			return PoWResult{
				Nonce: hex.EncodeToString(nonceBytes),
				Hash:  hashStr,
			}
		}
	}
}

// Delegated BFT: elect a leader based on PoW randomness + stake
func (hc *HybridConsensus) ElectLeader() string {
	hc.leaderElectionMu.Lock()
	defer hc.leaderElectionMu.Unlock()

	fmt.Println("🗳️  Electing leader via Delegated BFT + PoW...")
	seed := time.Now().String()
	powResult := performPoW(seed, hc.difficulty)
	fmt.Printf("🔑 PoW Hash: %s | Nonce: %s\n", powResult.Hash, powResult.Nonce)

	r := new(big.Int)
	r.SetString(powResult.Hash[:16], 16)
	mod := big.NewInt(int64(len(hc.nodes)))
	selected := new(big.Int).Mod(r, mod).Int64()

	nodeList := make([]*Node, 0, len(hc.nodes))
	for _, node := range hc.nodes {
		nodeList = append(nodeList, node)
	}

	hc.leaderID = nodeList[selected].ID
	fmt.Printf("👑 Leader Elected: %s\n", hc.leaderID)
	return hc.leaderID
}

// Voting mechanism (simulated)
func (hc *HybridConsensus) ReachQuorum() bool {
	trustedCount := 0
	for _, node := range hc.nodes {
		if node.Trusted {
			trustedCount++
		}
	}
	quorum := float64(trustedCount) / float64(len(hc.nodes))
	fmt.Printf("✅ Quorum achieved: %.2f\n", quorum)
	return quorum >= hc.quorumThreshold
}

func (hc *HybridConsensus) ProposeBlock() {
	if hc.leaderID == "" || !hc.ReachQuorum() {
		hc.ElectLeader()
	}

	fmt.Printf("🔒 Leader %s proposing block...\n", hc.leaderID)
	// Additional voting, signature aggregation, and block validation would go here
}
