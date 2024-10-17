package address

// AddressHandler defines the interface for address validation and conversion.
type AddressHandler interface {
	Validate(address string) error          // Validate the format of the address.
	ToBytes(address string) ([]byte, error) // Convert the address to a byte slice.
}
