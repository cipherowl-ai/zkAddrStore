package store

import (
	"addressdb/address"
	"github.com/ethereum/go-ethereum/crypto"
	"os"
	"testing"
)

func TestBloomFilterStoreFunctional(t *testing.T) {
	addressHandler := &address.EVMAddressHandler{} // Assuming you have a constructor for AddressHandler
	bf, err := NewBloomFilterStore(100, 0.1, addressHandler)
	if err != nil {
		t.Fatalf("Failed to create BloomFilterStore: %v", err)
	}

	addresses := []string{createAddress(), createAddress(), createAddress()}

	// Add addresses to the Bloom filter
	for _, addr := range addresses {
		if err := bf.AddAddress(addr); err != nil {
			t.Fatalf("Failed to add address %s: %v", addr, err)
		}
	}

	// Save the Bloom filter to a file
	filePath := os.TempDir() + "/bloomfilter.gob"
	if err := bf.SaveToFile(filePath); err != nil {
		t.Fatalf("Failed to save Bloom filter to file: %v", err)
	}

	// Create a new BloomFilterStore and load from the file
	bfReloaded, err := NewBloomFilterStoreFromFile(filePath, addressHandler)
	if err != nil {
		t.Fatalf("Failed to load Bloom filter from file: %v", err)
	}

	// Check if the addresses are present in the reloaded Bloom filter
	for _, addr := range addresses {
		exists, err := bfReloaded.CheckAddress(addr)
		if err != nil {
			t.Fatalf("Error checking address %s: %v", addr, err)
		}
		if !exists {
			t.Errorf("Address %s should be present in the Bloom filter after reload", addr)
		}
	}

	// Clean up the test file
	os.Remove(filePath)
}

func createAddress() string {
	key, _ := crypto.GenerateKey()
	return crypto.PubkeyToAddress(key.PublicKey).Hex()
}
