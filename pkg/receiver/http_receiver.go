package receiver

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/3d0c/storage/pkg/config"
)

// HTTPReceiver struct
type HTTPReceiver struct {
	objectID string
	nodes    []string
}

// NewHTTPReceiver constructor
func NewHTTPReceiver(objectID string, nodes []int) (*HTTPReceiver, error) {
	var (
		result = &HTTPReceiver{
			objectID: objectID,
			nodes:    make([]string, 0, cap(nodes)),
		}
	)

	for _, v := range nodes {
		result.nodes = append(result.nodes, config.Proxy().Nodes.Get(v))
	}

	return result, nil
}

func (r *HTTPReceiver) Recv(dst io.Writer) error {
	var (
		err error
	)

	for _, node := range r.nodes {
		url := fmt.Sprintf("%s/%s", node, r.objectID)
		fmt.Printf("=== Getting %s\n", url)
		if err = r.copyPart(url, dst); err != nil {
			return fmt.Errorf("error copying from '%s' - %s", url, err)
		}
	}

	return nil
}

func (r *HTTPReceiver) copyPart(url string, dst io.Writer) error {
	var (
		client      = new(http.Client)
		req         *http.Request
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*20)
		resp        *http.Response
		err         error
	)
	defer cancel()

	if req, err = http.NewRequestWithContext(ctx, "GET", url, nil); err != nil {
		return fmt.Errorf("error creating reeiver http request - %s", err)
	}

	if resp, err = client.Do(req); err != nil {
		return fmt.Errorf("error doing receiver request %v - %s", req, err)
	}

	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		return fmt.Errorf("non 2xx response status code for '%s', StatusCode - %d", url, resp.StatusCode)
	}
	defer resp.Body.Close()

	if _, err = io.Copy(dst, resp.Body); err != nil {
		return err
	}

	return nil
}
