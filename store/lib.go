package store

import (
	"encoding/gob"
	"os"

	"github.com/bits-and-blooms/bloom/v3"
)

// BloomFilterFromFile opens a .gob file and returns a bloom.BloomFilter
func BloomFilterFromFile(filename string) (*bloom.BloomFilter, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var filter *bloom.BloomFilter
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&filter); err != nil {
		return nil, err
	}
	return filter, nil
}
