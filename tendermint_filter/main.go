package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"os"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/tendermint/tendermint/abci/server"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/version"
)

type SimpleApp struct {
	addresses [][]byte
}

var _ abcitypes.Application = (*SimpleApp)(nil)

func NewSimpleApp() *SimpleApp {
	return &SimpleApp{}
}

func (app *SimpleApp) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{
		Data:             "Simple App",
		Version:          version.ABCIVersion,
		AppVersion:       1,
		LastBlockHeight:  0,
		LastBlockAppHash: []byte{},
	}
}

func (app *SimpleApp) SetOption(req abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
	return abcitypes.ResponseSetOption{}
}

func (app *SimpleApp) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	return abcitypes.ResponseDeliverTx{Code: 0}
}

func (app *SimpleApp) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	// Assuming the transaction is just an address (20 bytes)
	if len(req.Tx) != 20 {
		return abcitypes.ResponseCheckTx{Code: 1, Log: "Invalid address length"}
	}

	// Add the address to our list
	app.addresses = append(app.addresses, req.Tx)

	// List out all addresses
	for i, addr := range app.addresses {
		fmt.Printf("Address %d: %X\n", i, addr)
	}

	return abcitypes.ResponseCheckTx{Code: 0}
}

func (app *SimpleApp) Commit() abcitypes.ResponseCommit {
	return abcitypes.ResponseCommit{}
}

func (app *SimpleApp) Query(req abcitypes.RequestQuery) abcitypes.ResponseQuery {
	return abcitypes.ResponseQuery{Code: 0}
}

func (app *SimpleApp) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	return abcitypes.ResponseInitChain{}
}

func (app *SimpleApp) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	return abcitypes.ResponseBeginBlock{}
}

func (app *SimpleApp) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
}

// New method to satisfy the Application interface
func (app *SimpleApp) ListSnapshots(req abcitypes.RequestListSnapshots) abcitypes.ResponseListSnapshots {
	return abcitypes.ResponseListSnapshots{}
}

// New method to satisfy the Application interface
func (app *SimpleApp) OfferSnapshot(req abcitypes.RequestOfferSnapshot) abcitypes.ResponseOfferSnapshot {
	return abcitypes.ResponseOfferSnapshot{}
}

// New method to satisfy the Application interface
func (app *SimpleApp) LoadSnapshotChunk(req abcitypes.RequestLoadSnapshotChunk) abcitypes.ResponseLoadSnapshotChunk {
	return abcitypes.ResponseLoadSnapshotChunk{}
}

// New method to satisfy the Application interface
func (app *SimpleApp) ApplySnapshotChunk(req abcitypes.RequestApplySnapshotChunk) abcitypes.ResponseApplySnapshotChunk {
	return abcitypes.ResponseApplySnapshotChunk{}
}

func main() {

	filename := flag.String("f", "bloomfilter.gob", "Path to the .gob file containing the Bloom filter")
	// Parse the command-line flags
	flag.Parse()
	// Open the serialized Bloom filter file
	file, err := os.Open(*filename)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", *filename, err)
		return
	}
	defer file.Close()

	// Decode the Bloom filter
	var filter *bloom.BloomFilter
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&filter); err != nil {
		fmt.Println("Error decoding bloom filter:", err)
		return
	}

	app := NewSimpleApp()

	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))

	server, err := server.NewServer(":26658", "socket", app)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating server: %v", err)
		os.Exit(1)
	}
	server.SetLogger(logger)

	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "error starting server: %v", err)
		os.Exit(1)
	}

	defer server.Stop()

	// Wait forever
	select {}
}
