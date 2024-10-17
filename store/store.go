package store

import (
	"addressdb/address"
	"bufio"
	"errors"
	"fmt"
	"github.com/bits-and-blooms/bloom/v3"
	"os"
	"sync"
)

type BloomFilterStore struct {
	filter         *bloom.BloomFilter
	addressHandler address.AddressHandler
	mu             sync.RWMutex // Mutex to handle concurrent reloads.
}

// NewBloomFilterStore creates a new Bloom filter with file monitoring capabilities.
func NewBloomFilterStore(capacity uint, falsePositiveRate float64, addressHandler address.AddressHandler) (*BloomFilterStore, error) {
	if falsePositiveRate <= 0 || falsePositiveRate >= 1 {
		return nil, errors.New("invalid false positive rate: must be between 0 and 1")
	}

	filter := bloom.NewWithEstimates(capacity, falsePositiveRate)

	bf := &BloomFilterStore{
		filter:         filter,
		addressHandler: addressHandler,
	}

	return bf, nil
}

// NewBloomFilterStoreFromFile creates a new Bloom filter from a file.
func NewBloomFilterStoreFromFile(filePath string, addressHandler address.AddressHandler) (*BloomFilterStore, error) {
	bf := &BloomFilterStore{
		addressHandler: addressHandler,
	}

	if err := bf.LoadFromFile(filePath); err != nil {
		return nil, err
	}

	return bf, nil
}

// AddAddress inserts an address into the Bloom filter and encrypts the filter.
func (bf *BloomFilterStore) AddAddress(address string) error {
	// Validate and convert address
	if err := bf.addressHandler.Validate(address); err != nil {
		return err
	}

	addressBytes, err := bf.addressHandler.ToBytes(address)
	if err != nil {
		return err
	}

	bf.mu.Lock()
	defer bf.mu.Unlock()

	// Add to the Bloom filter
	bf.filter.Add(addressBytes)

	return nil
}

// CheckAddress decrypts the Bloom filter and checks if an address is in the filter.
func (bf *BloomFilterStore) CheckAddress(address string) (bool, error) {
	// Validate and convert address
	if err := bf.addressHandler.Validate(address); err != nil {
		return false, err
	}

	addressBytes, err := bf.addressHandler.ToBytes(address)
	if err != nil {
		return false, err
	}

	bf.mu.RLock()
	defer bf.mu.RUnlock()

	// Check if the address is in the filter
	return bf.filter.Test(addressBytes), nil
}

// LoadFromFile loads the Bloom filter from the specified file and decrypts it.
func (bf *BloomFilterStore) LoadFromFile(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("no file path specified for loading")
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}

	//TODO: Decrypt the data and validate the integrity of the decrypted data.
	r := bufio.NewReader(f)

	var filter bloom.BloomFilter
	if _, err := filter.ReadFrom(r); err != nil {
		return fmt.Errorf("failed to read Bloom filter: %v", err)
	}

	bf.mu.Lock()
	defer bf.mu.Unlock()

	bf.filter = &filter

	return nil
}

// SaveToFile saves the encrypted Bloom filter to the specified file.
func (bf *BloomFilterStore) SaveToFile(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("no file path specified for saving")
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	if _, err := bf.filter.WriteTo(w); err != nil {
		return err
	}

	return w.Flush()
}
