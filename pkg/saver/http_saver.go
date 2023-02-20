package saver

import (
	"fmt"
	"io"

	"github.com/3d0c/storage/pkg/config"
	"github.com/3d0c/storage/pkg/database/models"
	"github.com/3d0c/storage/pkg/utils"
)

// HTTPSaver struct
type HTTPSaver struct {
	objectID string
	cfg      config.ProxyConfig
	saved    []string
	chunk    Chunk
}

// NewHTTPSaver constructor
func NewHTTPSaver(objectID string, size int, c config.ProxyConfig) (*HTTPSaver, error) {
	httpSaver := &HTTPSaver{
		objectID: objectID,
		cfg:      c,
		saved:    make([]string, 0),
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

	if pm, err = models.NewPartsModel(s.cfg.Database); err != nil {
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

		node, nodeID := s.cfg.Nodes.Pick()
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

		if err = utils.Request("PUT", node, pr); err != nil {
			if err = s.rollback(); err != nil {
				return fmt.Errorf("error rollbacking transaction - %s", err)
			}
			return fmt.Errorf("error saving part #%d - %s", i, err)
		}
		pr.Close()

		if err = pm.Add(s.objectID, nodeID, i); err != nil {
			// Remove all saved chunks from nodes
			if err = s.rollback(); err != nil {
				return fmt.Errorf("error removing chunks - %s", err)
			}
			// Database TX rollback
			if err = pm.Rollback(); err != nil {
				return fmt.Errorf("error rollbacking transaction - %s", err)
			}
			return fmt.Errorf("error adding part - %s", err)
		}

		s.saved = append(s.saved, node)
	}

	if err = pm.Commit(); err != nil {
		return fmt.Errorf("error commiting transaction - %s", err)
	}

	return nil
}

func (s *HTTPSaver) rollback() error {
	for _, url := range s.saved {
		if err := utils.Request("DELETE", url, nil); err != nil {
			return fmt.Errorf("error deleting chunk from '%s' - %s", url, err)
		}
	}

	return nil
}
