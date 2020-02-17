package bloomfilter

import (
	"github.com/wuduozhi/go-datastructures/bitarray"
	"hash"
	"hash/fnv"
)

//type HashFunc func(key []byte) int

type BloomFilter struct {
	FilterSize  int
	HashFuncNum int
	ElementNum  int
	bitmap      bitarray.BitArray
	HashFunc    hash.Hash64
}

func (bf *BloomFilter) Add(key []byte) {
	h1, h2 := bf.getHash(key)

	for i := 0; i < bf.HashFuncNum; i++ {
		ind := (h1 + uint32(i)*h2) % uint32(bf.FilterSize)
		bf.bitmap.SetBit(uint64(ind))
	}
	bf.ElementNum++
}

func (bf *BloomFilter) Check(key []byte) bool {
	h1, h2 := bf.getHash(key)
	result := true
	for i := 0; i < bf.HashFuncNum; i++ {
		ind := (h1 + uint32(i)*h2) % uint32(bf.FilterSize)
		ok, _ := bf.bitmap.GetBit(uint64(ind))
		result = result && ok
	}
	return result
}

func (bf *BloomFilter) getHash(b []byte) (uint32, uint32) {
	bf.HashFunc.Reset()
	bf.HashFunc.Write(b)
	hash64 := bf.HashFunc.Sum64()
	h1 := uint32(hash64 & ((1 << 32) - 1))
	h2 := uint32(hash64 >> 32)
	return h1, h2
}

func NewBloomFilter(numHashFuncs, bfSize int) *BloomFilter {
	return &BloomFilter{
		FilterSize:  bfSize,
		HashFuncNum: numHashFuncs,
		ElementNum:  0,
		bitmap:      bitarray.NewBitArray(uint64(bfSize)),
		HashFunc:    fnv.New64(),
	}
}
