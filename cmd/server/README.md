# GES - Good Enough Server

Command-line flags:

- `-f`: Path to the .gob file containing the Bloom filter (default: "bloomfilter.gob")
- `-p`: Port to listen on (default: 8080)
- `-r`: Rate limit for requests per second (default: 20)
- `-b`: Burst limit for rate limiting (default: 5)

## API Endpoints

1. Single Address Check (GET)

   ```
   GET /check?s=<address>
   ```

   Example:
   ```bash
   curl "http://localhost:8080/check?s=0xbd4649c52778bb9259d5cd38e97a936eab57a194"
   ```

   Response:
   ```json
   {"found": true}
   ```

2. Batch Address Check (POST)

   ```
   POST /checkBatch
   Content-Type: application/json

   {
     "addresses": ["address1", "address2", ...]
   }
   ```

   Example:
   ```bash
   curl -X POST -H "Content-Type: application/json" -d '{
     "addresses": [
       "0x8b063eEd5bD1628e2D2b02FfCce4917E11558E69",
       "0x6E0F47C6F0C0F97c42956ffb11650B63c97ec9Ea",
       "NOONONO",
       "BBBB"
     ]
   }' http://localhost:8080/checkBatch
   ```

   Response:
   ```json
   {
     "found": ["0x8b063eEd5bD1628e2D2b02FfCce4917E11558E69"],
     "notfound": ["0x6E0F47C6F0C0F97c42956ffb11650B63c97ec9Ea", "NOONONO", "BBBB"],
     "found_count": 1,
     "notfound_count": 3
   }
   ```

## Performance

The server is designed for high performance, especially for batch checks. While exact performance metrics can vary depending on hardware and network conditions, here are some general observations:

- Single address checks typically respond in under 1ms in docker container on macbook pro m3Max.
- Batch checks can process 1k of addresses in 3ms.

The use of a Bloom filter allows for extremely fast membership tests with a very low false positive rate, making this solution ideal for quickly checking large sets of addresses.

## Rate Limiting

The server implements configurable rate limiting to prevent abuse. By default, it allows 20 requests per second with a burst of 5 requests. These values can be adjusted using the `-r` and `-b` flags when starting the server.

## Logging

The server logs each request, including the method, URI, remote address, and processing time, which can be useful for monitoring and debugging.

## Graceful Shutdown

The server supports graceful shutdown, ensuring that ongoing requests are completed before the server stops when receiving an interrupt signal.

## Contributing

Please read CONTRIBUTING.md for details on our code of conduct and the process for submitting pull requests.

## License

see the LICENSE file for details.
