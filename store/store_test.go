package store

import (
	"addressdb/address"
	"addressdb/securedata"
	pgp "github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestBloomFilterStoreFunctional(t *testing.T) {
	addressHandler := &address.EVMAddressHandler{}

	aliceKeyPriv, aliceKeyPub, bobKeyPriv, bobKeyPub := generateTestKeys(t)
	aliceHandler, err := securedata.NewPGPSecureHandler(securedata.WithPrivateKey(aliceKeyPriv), securedata.WithPublicKey(bobKeyPub))
	require.NoError(t, err, "Failed to create Alice's secure data handler")

	bobHandler, err := securedata.NewPGPSecureHandler(securedata.WithPrivateKey(bobKeyPriv), securedata.WithPublicKey(aliceKeyPub))
	require.NoError(t, err, "Failed to create Bob's secure data handler")

	tests := []struct {
		name        string
		writer_opts []Option
		reader_opts []Option
	}{
		{"default", nil, nil},
		{"WithEstimates", []Option{WithEstimates(100, 0.0000001)}, nil},
		{"WithEncryption", []Option{WithSecureDataHandler(aliceHandler)}, []Option{WithSecureDataHandler(bobHandler)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf, err := NewBloomFilterStore(addressHandler, tt.writer_opts...)
			require.NoError(t, err, "Failed to create BloomFilterStore")

			addresses := []string{createAddress(), createAddress(), createAddress()}

			addAddressesToBloomFilter(t, bf, addresses)
			filePath := saveBloomFilterToFile(t, bf)
			defer os.Remove(filePath)

			bfReloaded, err2 := NewBloomFilterStoreFromFile(filePath, addressHandler, tt.reader_opts...)
			require.NoError(t, err2, "Failed to load Bloom filter from file")
			checkAddressesInBloomFilter(t, bfReloaded, addresses)
		})
	}
}

func addAddressesToBloomFilter(t *testing.T, bf *BloomFilterStore, addresses []string) {
	for _, addr := range addresses {
		if err := bf.AddAddress(addr); err != nil {
			t.Fatalf("Failed to add address %s: %v", addr, err)
		}
	}
}

func saveBloomFilterToFile(t *testing.T, bf *BloomFilterStore) string {
	filePath := os.TempDir() + "/bloomfilter.gob"
	if err := bf.SaveToFile(filePath); err != nil {
		t.Fatalf("Failed to save Bloom filter to file: %v", err)
	}
	return filePath
}

func checkAddressesInBloomFilter(t *testing.T, bf *BloomFilterStore, addresses []string) {
	for _, addr := range addresses {
		exists, err := bf.CheckAddress(addr)
		if err != nil {
			t.Fatalf("Error checking address %s: %v", addr, err)
		}
		if !exists {
			t.Errorf("Address %s should be present in the Bloom filter after reload", addr)
		}
	}
}

func createAddress() string {
	key, _ := crypto.GenerateKey()
	return crypto.PubkeyToAddress(key.PublicKey).Hex()
}

func generateTestKeys(t *testing.T) (*pgp.Key, *pgp.Key, *pgp.Key, *pgp.Key) {
	pgp := pgp.PGP()
	aliceKeyPriv, err := pgp.KeyGeneration().
		AddUserId("alice", "alice@alice.com").
		New().
		GenerateKey()
	require.NoError(t, err)
	aliceKeyPub, err := aliceKeyPriv.ToPublic()
	require.NoError(t, err)

	bobKeyPriv, err := pgp.KeyGeneration().
		AddUserId("bob", "bob@bob.com").
		New().
		GenerateKey()
	require.NoError(t, err)
	bobKeyPub, err := bobKeyPriv.ToPublic()
	require.NoError(t, err)

	return aliceKeyPriv, aliceKeyPub, bobKeyPriv, bobKeyPub
}
