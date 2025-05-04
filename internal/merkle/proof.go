package merkle

import (
	"crypto/sha256"
	"encoding/binary"
	"hash"
)

// A simple Bloom Filter for AMQ support
type BloomFilter struct {
	Bits   []bool
	Size   uint
	K      int
	Hasher hash.Hash
}

func NewBloomFilter(size uint, k int) *BloomFilter {
	return &BloomFilter{
		Bits:   make([]bool, size),
		Size:   size,
		K:      k,
		Hasher: sha256.New(),
	}
}

func (bf *BloomFilter) Add(item string) {
	for i := 0; i < bf.K; i++ {
		index := bf.hash(item, i) % bf.Size
		bf.Bits[index] = true
	}
}

func (bf *BloomFilter) MightContain(item string) bool {
	for i := 0; i < bf.K; i++ {
		index := bf.hash(item, i) % bf.Size
		if !bf.Bits[index] {
			return false
		}
	}
	return true
}

func (bf *BloomFilter) hash(item string, i int) uint {
	bf.Hasher.Reset()
	bf.Hasher.Write([]byte(item))
	bf.Hasher.Write([]byte{byte(i)})
	sum := bf.Hasher.Sum(nil)
	return uint(binary.BigEndian.Uint32(sum[:4]))
}
