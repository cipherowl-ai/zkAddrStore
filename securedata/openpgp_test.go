package securedata

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/stretchr/testify/require"
)

func generateTestKeys(t *testing.T) (*crypto.Key, *crypto.Key, *crypto.Key, *crypto.Key, *crypto.Key, *crypto.Key) {
	pgp := crypto.PGP()
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

	chadKeyPriv, err := pgp.KeyGeneration().
		AddUserId("chad", "chad@chad.com").
		New().
		GenerateKey()
	require.NoError(t, err)
	chadKeyPub, err := chadKeyPriv.ToPublic()
	require.NoError(t, err)

	return aliceKeyPriv, aliceKeyPub, bobKeyPriv, bobKeyPub, chadKeyPriv, chadKeyPub
}

func TestNewPGPSecureHandler(t *testing.T) {
	privKey, _, pubKey, _, _, _ := generateTestKeys(t)
	handler, err := NewPGPSecureHandler(WithPrivateKey(privKey), WithPublicKey(pubKey))
	require.NoError(t, err)
	require.NotNil(t, handler)
	assert.NotNil(t, handler.pgpHandle)
	assert.Equal(t, privKey, handler.privKey)
	assert.Equal(t, pubKey, handler.pubKey)
}

func TestWithPublicKeyPath(t *testing.T) {
	_, pubKey, _, _, _, _ := generateTestKeys(t)

	// Save the public key to a temporary file
	pubKeyData, err := pubKey.Armor()
	require.NoError(t, err)
	tmpFile, err := os.CreateTemp("", "pubkey_*.asc")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	_, err = tmpFile.Write([]byte(pubKeyData))
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	// Create OpenPGPSecureHandler using WithPublicKeyPath option
	handler, err := NewPGPSecureHandler(WithPublicKeyPath(tmpFile.Name()))
	require.NoError(t, err)
	require.NotNil(t, handler)
	assert.Equal(t, pubKey.GetFingerprint(), handler.pubKey.GetFingerprint())
}

func TestGpgKeys(t *testing.T) {
	pubkeyPath := "testdata/pubkey.asc"
	privkeyPath := "testdata/privkey.asc"

	// Create OpenPGPSecureHandler using WithPublicKeyPath option
	handler, err := NewPGPSecureHandler(WithPublicKeyPath(pubkeyPath), WithPrivateKeyPath(privkeyPath, "password123"))
	require.NoError(t, err)
	require.NotNil(t, handler)

	require.NoError(t, err)
	require.NotNil(t, handler)

	var buf bytes.Buffer
	writer, err := handler.Writer(&buf)
	require.NoError(t, err)

	_, err = writer.Write([]byte("test data"))
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	require.NoError(t, err)
	require.NotNil(t, handler)

	reader, err := handler.Reader(&buf)
	require.NoError(t, err)

	decryptedData, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, "test data", string(decryptedData))
}

func TestWithPrivateKeyPath(t *testing.T) {
	privKey, _, _, _, _, _ := generateTestKeys(t)
	tests := []struct {
		name       string
		privateKey *crypto.Key
		passphrase string
	}{
		{"with passphrase", privKey, "testpassphrase"},
		{"without passphrase", privKey, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var privKey *crypto.Key = tt.privateKey
			var err error
			if tt.passphrase != "" {
				pgp := crypto.PGP()
				privKey, err = pgp.LockKey(tt.privateKey, []byte(tt.passphrase))
				require.NoError(t, err)
			}

			// Save the private key to a temporary file
			privKeyData, err := privKey.Armor()
			require.NoError(t, err)
			tmpFile, err := os.CreateTemp("", "privkey_*.asc")
			require.NoError(t, err)
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.Write([]byte(privKeyData))
			require.NoError(t, err)
			require.NoError(t, tmpFile.Close())

			// Create OpenPGPSecureHandler using WithPrivateKeyPath option
			handler, err := NewPGPSecureHandler(WithPrivateKeyPath(tmpFile.Name(), tt.passphrase))
			require.NoError(t, err)
			require.NotNil(t, handler)
			assert.Equal(t, privKey.GetFingerprint(), handler.privKey.GetFingerprint())
		})

	}
}

func TestPGPSecureHandler_Writer(t *testing.T) {
	privKey, _, pubKey, _, _, _ := generateTestKeys(t)
	handler, err := NewPGPSecureHandler(WithPrivateKey(privKey), WithPublicKey(pubKey))
	require.NoError(t, err)
	require.NotNil(t, handler)

	var buf bytes.Buffer
	writer, err := handler.Writer(&buf)
	require.NoError(t, err)

	_, err = writer.Write([]byte("test data"))
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	require.NotEmpty(t, buf.Bytes())
}

func TestPGPSecureHandler_Reader(t *testing.T) {
	aliceKeyPriv, aliceKeyPub, bobKeyPriv, bobKeyPub, chadKeyPriv, _ := generateTestKeys(t)
	tests := []struct {
		name           string
		senderPrivKey  *crypto.Key
		receiverPubKey *crypto.Key
		readerPrivKey  *crypto.Key
		readerPubKey   *crypto.Key
		wantErr        bool
		wantReadErr    bool
	}{
		{"alice to bob, bob reads", aliceKeyPriv, bobKeyPub, bobKeyPriv, aliceKeyPub, false, false},
		{"alice to bob, chad reads", aliceKeyPriv, bobKeyPub, chadKeyPriv, aliceKeyPub, true, true},
		{"chad impersonates alice to bob, bob reads", chadKeyPriv, bobKeyPub, bobKeyPriv, aliceKeyPub, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendHandler, err := NewPGPSecureHandler(WithPrivateKey(tt.senderPrivKey), WithPublicKey(tt.receiverPubKey))
			require.NoError(t, err)
			require.NotNil(t, sendHandler)

			var buf bytes.Buffer
			writer, err := sendHandler.Writer(&buf)
			require.NoError(t, err)

			_, err = writer.Write([]byte("test data"))
			require.NoError(t, err)
			require.NoError(t, writer.Close())

			readHandler, err := NewPGPSecureHandler(WithPrivateKey(tt.readerPrivKey), WithPublicKey(tt.readerPubKey))
			require.NoError(t, err)
			require.NotNil(t, readHandler)

			reader, err := readHandler.Reader(&buf)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			decryptedData, err := io.ReadAll(reader)

			if tt.wantReadErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "test data", string(decryptedData))
			}
		})
	}
}
