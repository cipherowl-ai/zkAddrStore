package securedata

import (
	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/stretchr/testify/require"
	"testing"
)

func GenerateTestKeys(t *testing.T) []*crypto.Key {
	var keys []*crypto.Key
	keys = append(keys, generateKeyPair(t, "alice", "alice@alice.com")...)
	keys = append(keys, generateKeyPair(t, "bob", "bob@bob.com")...)
	keys = append(keys, generateKeyPair(t, "chad", "chad@chad.com")...)
	return keys
}

func generateKeyPair(t *testing.T, name, email string) []*crypto.Key {
	pgp := crypto.PGP()
	privKey, err := pgp.KeyGeneration().
		AddUserId(name, email).
		New().
		GenerateKey()
	require.NoError(t, err)
	pubKey, err := privKey.ToPublic()
	require.NoError(t, err)
	return []*crypto.Key{privKey, pubKey}
}
