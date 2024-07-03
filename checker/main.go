package main

import (
	"bufio"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/bits-and-blooms/bloom/v3"
)

func main() {
	// Define the command-line flag
	filename := flag.String("f", "bloomfilter.gob", "Path to the .gob file containing the Bloom filter")

	// Parse the command-line flags
	flag.Parse()

	// Open the serialized Bloom filter file
	file, err := os.Open(*filename)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", *filename, err)
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
	fmt.Println("Enter strings to check. Press Ctrl+D to exit.")

	// Process input from user
	for {
		fmt.Print("Enter string: ")
		if !scanner.Scan() {
			if scanner.Err() == nil {
				// This means we've reached EOF (Ctrl+D)
				fmt.Println("\nReached end of input. Exiting.")
			} else {
				fmt.Printf("Error reading input: %v\n", scanner.Err())
			}
			break
		}
		input := scanner.Text()
		if filter.TestString(input) {
			fmt.Println("Possibly in set.")
		} else {
			fmt.Println("Definitely not in set.")
		}
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		fmt.Printf("Error reading from standard input: %v\n", err)
	}
}
