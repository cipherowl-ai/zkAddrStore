package main

import (
	"addressdb/lib"
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	// Define the command-line flag
	filename := flag.String("f", "bloomfilter.gob", "Path to the .gob file containing the Bloom filter")

	// Parse the command-line flags
	flag.Parse()

	start := time.Now()
	// Open the serialized Bloom filter file
	filter, err := lib.BloomFilterFromFile(*filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(-1)

	}
	elapsed := time.Since(start)
	fmt.Printf("> Time taken to load bloomfilter: %v\n", elapsed)

	// Create a scanner to read input from standard input
	scanner := bufio.NewScanner(os.Stdin)

	// measure the time it takes to check the bloom filter
	start = time.Now()

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
		// only handle non-existing entries
		if !filter.TestString(input) {
			fmt.Println("NOT in set:", input)
		}
	}
	elapsed = time.Since(start)
	fmt.Printf("> Time taken to check bloomfilter: %v\n", elapsed)

	if err := scanner.Err(); err != nil && err != io.EOF {
		fmt.Fprintf(os.Stderr, "Error reading from standard input: %v\n", err)
	}
}
