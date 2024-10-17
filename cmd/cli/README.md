# CLI tool for working with address bloom filters

## usage

### Encoder to build a bloomfilter from a list of addresses
```bash
pa-cli encode --input ./addresses.txt --output ./bloomfilter.gob -n 1000000 -p 0.000001
```

- `-input`: Input file path, one address per line
- `-output`: Output file path, it is a binary bloomfilter file, the content is not human readable.
- `-n`: Number of entries (should match the number of generated addresses)
- `-p`: False positive rate. e.g. 0.000001 is 1 in a million.

### Console based interactive client for testing bloomfilter

```bash
pa-cli check -f=./bloomfilter.gob
```

### A Batch Checker for bloomfilter

```bash
cat btc_tocheck.txt | pa-cli batch-check -f=./bloomfilter.gob > /tmp/missing.txt
```

Where btc_tocheck.txt is a file with one address per line


### A Batch Checker for bloomfilter

```bash
pa-cli generate-addresses --output ./addresses.txt -n 1000000
```

Where btc_tocheck.txt is a file with one address per line