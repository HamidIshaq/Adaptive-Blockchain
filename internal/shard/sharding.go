package shard

import "fmt"

func SplitShard(shardID string) []string {
	fmt.Printf("Splitting shard: %s due to high load\n", shardID)
	// Logic to split based on load metrics
	return []string{"shardID-1", "shardID-2"}
}
