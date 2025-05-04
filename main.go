package main

import (
	"adaptive-blockchain/internal/blockchain"
	"encoding/json"
	"fmt"
	"net/http"
)

var bc *blockchain.Blockchain

func main() {
	fmt.Println("🚀 Starting Adaptive Blockchain Node...")

	bc = blockchain.NewBlockchain()

	go bc.Start()

	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/verify", verifyHandler)

	fmt.Println("🌐 Web dashboard at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"blocks": bc.BlockHeight(),
		"shards": bc.ShardCount(),
	}
	json.NewEncoder(w).Encode(status)
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	tx := r.URL.Query().Get("tx")
	found := false
	for _, tree := range bc.MerkleForest.Shards {
		if tree.Filter.MightContain(tx) {
			found = true
			break
		}
	}
	resp := map[string]bool{"possible": found}
	json.NewEncoder(w).Encode(resp)
}
