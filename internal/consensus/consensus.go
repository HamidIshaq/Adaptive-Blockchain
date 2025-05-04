package consensus

import "fmt"

type HybridConsensus struct{}

func NewHybridConsensus() *HybridConsensus {
	return &HybridConsensus{}
}

func (hc *HybridConsensus) ProposeBlock() {
	fmt.Println("🔒 Hybrid consensus proposing a block...")
}
