package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"os"

	"github.com/bits-and-blooms/bloom/v3"
)

// this code snippet can be repurpose to make it be a http server
func main() {
	// Open the serialized Bloom filter file
	file, err := os.Open("bloomfilter.gob")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Decode the Bloom filter
	var filter *bloom.BloomFilter
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&filter); err != nil {
		fmt.Println("Error decoding bloom filter:", err)
		return
	}

	// Create a scanner to read input from standard input
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter strings to check. Type 'exit' to quit.")

	// Process input from user
	for {
		fmt.Print("Enter string: ")
		scanner.Scan()
		input := scanner.Text()
		if input == "exit" {
			break
		}
		if filter.TestString(input) {
			fmt.Println("Possibly in set.")
		} else {
			fmt.Println("Definitely not in set.")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from standard input:", err)
	}
}
