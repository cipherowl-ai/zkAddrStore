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

	// read from standard input, till EOF, and check if the string is in the bloom filter

	for {
		if !scanner.Scan() {
			if scanner.Err() != nil {
				// print to stderr
				fmt.Fprintf(os.Stderr, "Error reading from standard input: %v\n", err)
			}
			break
		}
		input := scanner.Text()
		if filter.TestString(input) {
		} else {
			fmt.Println("NOT in set:", input)
		}
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		fmt.Printf("Error reading from standard input: %v\n", err)
	}
}
