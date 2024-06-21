# This is a project demosntrating using bloom filter to store lage number of private blockchain address

## Step 1

Generate 1M EVM addresses

```bash
# generate 1M EVM addresses, resume is stored in address.txt
go run evmaddress_generator/main.go -n 1000000

```

## Step 2

Build the bloom filter

```bash
# use the content in address.txt to build the bloom filter, 1 entry per line, and stored in bloomfilter.gob
go run encoder/main.go
ls -alh bloomfilter.gob
```

## Step 3

```bash
# load bloomfilter.gob, and start a shell to check if the address is in the bloom filter
> go run checker/main.go

Enter strings to check. Type 'exit' to quit.
Enter string: 0x6864A451C800D21B9c5673A8153E3aD47cEBc94d
Possibly in set.
Enter string: xxxx
Definitely not in set.
Enter string: 0x7e531DE4a901a88b2540a6973797d6AAA75F2fdF

```