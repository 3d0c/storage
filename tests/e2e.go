package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/3d0c/storage/pkg/apiserver/handlers"
	"github.com/3d0c/storage/pkg/config"
	"github.com/3d0c/storage/pkg/node"
	"github.com/3d0c/storage/pkg/utils"
)

func main() {
	var (
		nodes    []string
		obtained []byte
		status   int
		err      error
	)

	if nodes, err = startNodes(); err != nil {
		fmt.Printf("error starting nodes - %s", err)
		os.Exit(-1)
	}

	payload := utils.RandomBytes()
	objectID := utils.RandomInt()
	expectedHash := utils.MakeSHA256Hash(payload)

	// no needs to run proxy API server, just use Put handler with corresponding request
	fileHandler := handlers.FileHandler(config.ProxyConfig{
		Nodes: nodes,
	})

	url := fmt.Sprintf("localhost/%d", objectID)
	put := httptest.NewRequest("PUT", url, bytes.NewReader(payload))

	if status, err = fileHandler.Put(nil, put); err != nil || status != http.StatusOK {
		fmt.Printf("Test failed. Error doing proxy request - %d, %s\n", status, err)
		os.Exit(-1)
	}

	get := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()

	if status, err = fileHandler.Get(w, get); err != nil || status != http.StatusOK {
		fmt.Printf("Test failes. Error getting data - %d, %s\n", status, err)
		os.Exit(-1)
	}

	resp := w.Result()

	if obtained, err = io.ReadAll(resp.Body); err != nil {
		fmt.Printf("Test failes. Error reading data - %s\n", err)
		os.Exit(-1)
	}

	obtainedHash := utils.MakeSHA256Hash(obtained)

	if expectedHash != obtainedHash {
		fmt.Printf("Test failed with different hashed %s != %s\n", expectedHash, obtainedHash)
	}

	fmt.Printf("HAPPY PASS\n")
	os.Exit(0)
}

func startNodes() ([]string, error) {
	var (
		count       = utils.RandomSeed(5, 9)
		nodesURL    = make([]string, 0, count)
		ctx, cancel = context.WithCancel(context.Background())
		apiSrv      *node.APIHTTPServer
		err         error
	)
	defer cancel()

	for i := 0; i < count; i++ {
		cfg := config.NodeConfig{}
		cfg.Address = fmt.Sprintf(":900%d", i)
		cfg.StorageDir = fmt.Sprintf("/tmp/node-%d", i)

		if apiSrv, err = node.NewAPIHTTPServer(cfg); err != nil {
			return nil, fmt.Errorf("error initializing node API server - %s", err)
		}

		apiSrv.Run(ctx)

		nodesURL = append(nodesURL, cfg.Address)
	}

	return nodesURL, nil
}
