# ZkAddressStore

## Status Of The Project

The project is currently in the beta phrase, it will contain breaking changes in the future.

## Overview

ZkAddressStore demonstrates the use of Bloom filters to share sets of blockchain addresses while preserving privacy. It provides tools for generating Ethereum addresses, encoding them into a Bloom filter, and checking addresses against the filter.

## Key Features

- Generate large sets of Ethereum addresses
- Encode addresses into a space-efficient Bloom filter
- Check if an address is potentially in the set
- Preserve privacy while sharing address sets

## How It Works

The project uses a Bloom filter, stored in a `.gob` file. This data structure allows for efficient storage and querying of large sets of data with a controllable false-positive rate.

Example: A set of 1 million Ethereum addresses results in a `.gob` file of approximately 3.5MB (uncompressed). The actual size depends on the number of addresses and the chosen false-positive rate.

The `.gob` file can be easily shared across data pipelines. Adding key-pair encryption (not implemented in this version) would further enhance confidentiality.

## Installation

```bash
git clone https://github.com/your-username/ZkAddressStore.git
cd ZkAddressStore
go mod tidy
```

## Usage

### Step 1: Generate Ethereum Addresses (Optional)

```bash
go run evmaddress_generator/main.go -n 1000000
```

Generates 1 million Ethereum addresses and stores them in `address.txt`.

### Step 2: Build the Bloom Filter

```bash
go run encoder/main.go -n 1000000 -p 0.000001
```

- `-n`: Number of entries (should match the number of generated addresses)
- `-p`: False positive rate

Creates a `bloomfilter.gob` file containing the Bloom filter.

### Step 3: Use the Filter

Interactive mode:

```bash
go run checker/main.go -f bloomfilter.gob
```

Batch mode:

```bash
cat my_addresses.txt | go run batch_checker/main.go -f bloomfilter.gob
```

## Large-Scale Example

Building a Bloom filter with 24 million Ethereum addresses:

```bash
go run encoder/main.go -n 1000000000 -p 0.000001 -input ~/path/to/eth_all.csv

# Check addresses
go run checker/main.go -f bloomfilter.gob
```

Result: `bloomfilter.gob` file of about 3.3GB.

## Performance

- Encoding 24M addresses: ~20 seconds on a standard machine.
- Constant-time complexity for adding and checking addresses.

### Benchmarks

```bash
# Debug build
time target/debug/pa-encoder -input ~/Downloads/eth_all.csv
# 13.59s user 0.33s system 97% cpu 14.211 total

# Release build
time target/release/pa-encoder -input ~/Downloads/eth_all.csv
# 4.32s user 0.26s system 94% cpu 4.835 total
```

#### Large-scale performance

- Loading a 51GB Bloom filter (1.4 billion addresses): 18 seconds
- Checking 1173 addresses: ~2ms

#### Small-scale performance

```bash
go run encoder/main.go -n 100000 -p 0.0000001
```

Results in a ~450KB filter.

## Limitations

- False-positive rate, but no false negatives
- No encryption for the `.gob` file in current implementation
- Potential for brute-force recovery of entries if the value space is limited

## Future Improvements

- Implement encryption for `.gob` files
- Optimize Bloom filter parameters for various use cases
- Add a web interface for easier interaction

## Contributing

Contributions are welcome! Please submit a Pull Request.​​​​​​​​​​​​​​​​

Yes we are welcome to any suggestion including a better name for the project.