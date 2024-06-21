package main

import (
	"bufio"
	"encoding/gob"
	"flag"
	"fmt"
	"os"

	"github.com/bits-and-blooms/bloom/v3"
)

func main() {
	// Define CLI flags
	nFlag := flag.Uint("n", 10000000, "number of elements expected")
	pFlag := flag.Float64("p", 0.00001, "false positive probability")
	inputFile := flag.String("input", "addresses.txt", "input file path")
	outputFile := flag.String("output", "bloomfilter.gob", "output file path")

	// Parse the flags
	flag.Parse()

	// Use the parsed values to create a Bloom filter
	filter := bloom.NewWithEstimates(uint(*nFlag), *pFlag)

	// Open an input file
	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read strings line by line and add them to the Bloom filter
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		filter.AddString(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from file:", err)
		return
	}

	// Serialize the Bloom filter to a file
	f, err := os.Create(*outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer f.Close()

	encoder := gob.NewEncoder(f)
	if err := encoder.Encode(filter); err != nil {
		fmt.Println("Error encoding bloom filter:", err)
		return
	}

	fmt.Println("Bloom filter has been serialized successfully.")
}
