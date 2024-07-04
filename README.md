# ZkAddressStore

## Overview

ZkAddressStore is a project that demonstrates the use of Bloom filters to share sets of blockchain addresses while preserving privacy. It provides tools for generating Ethereum addresses, encoding them into a Bloom filter, and checking addresses against the filter.

## Key Features

- Generate large sets of Ethereum addresses
- Encode addresses into a space-efficient Bloom filter
- Check if an address is potentially in the set
- Preserve privacy while sharing address sets

## How It Works

The project uses a Bloom filter, which is saved in a `.gob` file. This data structure allows for efficient storage and querying of large sets of data with a controllable false-positive rate.

For example, a set of 1 million Ethereum addresses results in a `.gob` file of approximately 3.5MB (uncompressed). The actual size depends on the number of addresses and the chosen false-positive rate.

The `.gob` file can be easily shared across data pipelines. With the addition of key-pair encryption (not implemented in this version), confidentiality can be further enhanced.

## Installation

```bash
git clone https://github.com/your-username/ZkAddressStore.git
cd ZkAddressStore
go mod tidy
```

## Usage

### Step 1: Generate Ethereum Addresses

Generate a set of Ethereum addresses:

```bash
go run evmaddress_generator/main.go -n 1000000
```

This generates 1 million Ethereum addresses and stores them in `address.txt`.

### Step 2: Build the Bloom Filter

Create a Bloom filter from the generated addresses:

```bash
go run encoder/main.go -n 1000000 -p 0.000001
```

- `-n`: Number of entries (should match the number of generated addresses)
- `-p`: False positive rate

This creates a `bloomfilter.gob` file containing the Bloom filter.

### Step 3: Check Addresses

Use the checker to verify if addresses are in the set:

```bash
go run checker/main.go -f bloomfilter.gob
```

This starts an interactive shell where you can enter addresses to check against the Bloom filter.

## Example: Large-Scale Usage

Here's an example of building a Bloom filter with 24 million Ethereum addresses:

```bash
# Assuming you have a file 'eth_all.csv' with 24M addresses, and the max number of addresses is 1000000000
go run encoder/main.go -n 1000000000 -p 0.000001 -input ~/path/to/eth_all.csv

# Check addresses
go run checker/main.go -f bloomfilter.gob
```

In this example, the resulting `bloomfilter.gob` file is about 3.3GB.

## Performance

- Encoding 24M addresses takes approximately 20 seconds on a standard machine.
- The Bloom filter provides constant-time complexity for both adding and checking addresses.

## Limitations

- Bloom filters have a false-positive rate but no false negatives.
- The current implementation does not include encryption for the `.gob` file.

## Future Improvements

- Implement encryption for the `.gob` file to enhance security.
- Optimize the Bloom filter parameters for different use cases.
- Add a web interface for easier interaction with the Bloom filter.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[Specify your license here]