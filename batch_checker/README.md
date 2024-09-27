# A Batch Checker for bloomfilter

```bash
cat btc_tocheck.txt | pa-batch-checker -f=./bloomfilter.gob > /tmp/missing.txt
```

Where btc_tocheck.txt is a file with one address per line
