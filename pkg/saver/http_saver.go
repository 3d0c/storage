package saver

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/3d0c/storage/pkg/config"
	"github.com/3d0c/storage/pkg/database/models"
)

// HTTPSaver struct
type HTTPSaver struct {
	objectID string
	nodes    config.Nodes
	chunk    Chunk
}

// NewHTTPSaver constructor
func NewHTTPSaver(objectID string, size int, nodes config.Nodes) (*HTTPSaver, error) {
	httpSaver := &HTTPSaver{
		objectID: objectID,
		nodes:    nodes,
		chunk:    NewChunk(size),
	}

	return httpSaver, nil
}

// Save interface implementation
func (s *HTTPSaver) Save(r io.Reader) error {
	var (
		pm  *models.Parts
		err error
	)

	if pm, err = models.NewPartsModel(); err != nil {
		return fmt.Errorf("error initializing Parts model - %s", err)
	}

	fmt.Printf("=== chunk: %v\n", s.chunk)

	for i := 0; i < ChunkCount(); i++ {
		chunkSize := s.chunk.Size
		// First chunk will always be larger
		// TODO For better distribution it can be randomly choosen
		if i == 0 {
			chunkSize += s.chunk.Modulo
		}

		// TODO Circular buffer might be an option for better distribution

		node, nodeID := s.nodes.Pick()
		node = fmt.Sprintf("%s/%s", node, s.objectID)

		// DEBUG
		fmt.Printf("=== Saving part #%d of length %d to node %s\n", i, chunkSize, node)

		pr, pw := io.Pipe()

		go func() {
			defer pw.Close()

			if _, err := io.CopyN(pw, r, int64(chunkSize)); err != nil {
				panic(fmt.Sprintf("error copying into pipe #%d - %s", i, err))
			}
		}()

		if err = s.save(node, pr); err != nil {
			return fmt.Errorf("error saving part #%d - %s", i, err)
		}
		pr.Close()

		// TODO ROLLBACK to be implemented
		if err = pm.Add(s.objectID, nodeID, i); err != nil {
			return fmt.Errorf("error adding part - %s", err)
		}
	}

	return nil
}

// TODO add read/write timeouts
func (s *HTTPSaver) save(url string, payload io.Reader) error {
	var (
		client      = new(http.Client)
		req         *http.Request
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*20)
		resp        *http.Response
		err         error
	)
	defer cancel()

	if req, err = http.NewRequestWithContext(ctx, "PUT", url, payload); err != nil {
		return fmt.Errorf("error creating http request - %s", err)
	}

	if resp, err = client.Do(req); err != nil {
		return fmt.Errorf("error doing request %v - %s", req, err)
	}

	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		return fmt.Errorf("non 2xx response status code")
	}

	return nil
}
