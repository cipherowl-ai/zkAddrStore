# Encoder to build a bloomfilter from a list of addresses

## usage

```bash
go run encoder/main.go -input=./addresses.txt -output=./bloomfilter.gob -n 1000000 -p 0.000001
```

- `-input`: Input file path, one address per line
- `-output`: Output file path, it is a binary bloomfilter file, the content is not human readable.
- `-n`: Number of entries (should match the number of generated addresses)
- `-p`: False positive rate. e.g. 0.000001 is 1 in a million.
