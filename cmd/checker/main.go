package main

import (
	"addressdb/store"
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	// Define the command-line flag
	filename := flag.String("f", "bloomfilter.gob", "Path to the .gob file containing the Bloom filter")

	// Parse the command-line flags
	flag.Parse()

	filter, err := store.BloomFilterFromFile(*filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(-1)
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
