package securedata

import (
	"io"
)

// SecureDataHandler defines an interface for securely handling data (encryption/decryption).
type SecureDataHandler interface {
	Writer(output io.Writer) (io.WriteCloser, error)  // Encrypts and returns an io.Writer
	Reader(input io.Reader) (VerifyDataReader, error) // Decrypts and returns an io.Reader
}

// VerifyDataReader is used for reading data that should be verified with a signature.
// It is needed because entire data needs to be read before verifying the signature, but bloom.BloomFilter
// does not always read the entire data it writes.
type VerifyDataReader interface {
	Read(b []byte) (n int, err error)
	// VerifySignature is used to verify that the embedded signatures are valid.
	// This method needs to be called once all the data has been read.
	VerifySignature() error
}
