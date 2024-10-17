package address

import (
	"testing"
)

func TestEVMAddressHandler_Validate(t *testing.T) {
	handler := &EVMAddressHandler{}

	tests := []struct {
		name    string
		address string
		valid   bool
	}{
		{"lower case", "0x1234567890abcdef1234567890abcdef12345678", true},
		{"mixed case", "0X1234567890ABCDEF1234567890ABCDEF12345678", true},
		{"no prefix", "1234567890abcdef1234567890abcdef12345678", false},
		{"too short", "0x1234567890abcdef1234567890abcdef1234567", false},
		{"too long", "0x1234567890abcdef1234567890abcdef123456789", false},
	}

	for _, test := range tests {
		err := handler.Validate(test.address)
		if (err == nil) != test.valid {
			t.Errorf("Validate(%q) = %v; want valid = %v", test.address, err, test.valid)
		}
	}
}

func TestEVMAddressHandler_ToBytes(t *testing.T) {
	handler := &EVMAddressHandler{}

	tests := []struct {
		name    string
		address string
		bytes   []byte
		err     bool
	}{
		{"lower case", "0x1234567890abcdef1234567890abcdef12345678", []byte{0x12, 0x34, 0x56, 0x78, 0x90, 0xAB, 0xCD, 0xEF, 0x12, 0x34, 0x56, 0x78, 0x90, 0xAB, 0xCD, 0xEF, 0x12, 0x34, 0x56, 0x78}, false},
		{"mixed case", "0X1234567890ABCDEF1234567890ABCDEF12345678", []byte{0x12, 0x34, 0x56, 0x78, 0x90, 0xAB, 0xCD, 0xEF, 0x12, 0x34, 0x56, 0x78, 0x90, 0xAB, 0xCD, 0xEF, 0x12, 0x34, 0x56, 0x78}, false},
		{"odd hex length", "0x1234567890abcdef1234567890abcdef1234567", nil, true},
		{"invalid hex", "0xG234567890abcdef1234567890abcdef12345678", nil, true},
	}

	for _, test := range tests {
		result, err := handler.ToBytes(test.address)
		if (err != nil) != test.err {
			t.Errorf("ToBytes(%q) = %v; want error = %v", test.address, err, test.err)
		}
		if !test.err && !equalBytes(result, test.bytes) {
			t.Errorf("ToBytes(%q) = %v; want %v", test.address, result, test.bytes)
		}
	}
}

// Helper function to compare byte slices
func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
