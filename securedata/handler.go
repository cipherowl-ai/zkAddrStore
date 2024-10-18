package securedata

import (
	"io"
)

// SecureDataHandler defines an interface for securely handling data (encryption/decryption).
type SecureDataHandler interface {
	Writer(output io.Writer) (io.WriteCloser, error) // Encrypts and returns an io.Writer
	Reader(input io.Reader) (io.Reader, error)       // Decrypts and returns an io.Reader
}
