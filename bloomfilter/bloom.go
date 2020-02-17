package bloomfilter

import "github.com/wuduozhi/go-datastructures/bitarray"

type BloomFilter struct {
	m uint
	k uint
	b *bitarray.BitArray
}

