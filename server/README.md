# Bloom Filter HTTP Server

This project implements an HTTP server that uses a Bloom filter to check if a given string is potentially in a pre-defined set. It's designed to be efficient and scalable, with features like rate limiting and graceful shutdown.

## Features

- Load a pre-generated Bloom filter from a .gob file
- HTTP endpoint to check if a string is in the set
- Rate limiting to prevent abuse
- Graceful shutdown
- Logging middleware
- Environment variable configuration

## Prerequisites

- Go 1.16 or later
- Git

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/bloom-filter-server.git
   cd bloom-filter-server
   ```

2. Install the required dependencies:
   ```
   go get github.com/gorilla/mux
   go get github.com/joho/godotenv
   go get golang.org/x/time/rate
   go get github.com/bits-and-blooms/bloom/v3
   ```

## Configuration

Create a `.env` file in the project root with the following content:

```
BLOOM_FILTER_FILE=path/to/your/bloomfilter.gob
PORT=8080
```

Adjust the `BLOOM_FILTER_FILE` path to point to your .gob file containing the Bloom filter.

## Running the Server

1. Make sure you have a valid Bloom filter .gob file at the location specified in your `.env` file.

2. Start the server:
   ```
   go run main.go
   ```

3. The server will start on the port specified in your `.env` file (default is 8080).

## Usage

To check if a string is in the set:

```
GET http://localhost:8080/check?s=yourstringhere
```

Example using curl:
```
curl "http://localhost:8080/check?s=0xbd4649c52778bb9259d5cd38e97a936eab57a194"
```

The server will respond with a JSON object:

```json
{
  "query": "example",
  "in_set": true,
  "message": "The string 'example' is possibly in the set."
}
```

## Development

### Project Structure

- `main.go`: The main application file containing the server logic.
- `.env`: Configuration file for environment variables.
- `bloomfilter.gob`: The pre-generated Bloom filter (not included in the repository).
