# ZkAddressStore

## what does it do?
This project showcases using a Bloom filter to share sets of addresses while preserving privacy.

The data structure is saved in a `.gob` file. For a set of 1 million Ethereum addresses, the resulting file is 30MB without compression. When compressed using gzip, it shrinks to around 13MB.

This .gob file can be easily shared across the data pipeline, and with the help of key pair, confidentiality can be achieved as well.

## Step 1

Generate 1M EVM addresses

```bash
# generate 1M EVM addresses, resume is stored in address.txt
go run evmaddress_generator/main.go -n 1000000

```

## Step 2

Build a bloom filter

```bash
# use the content in address.txt to build the bloom filter, one entry per line, and the result is stored in a .gob file.
go run encoder/main.go
ls -al bloomfilter.gob
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
