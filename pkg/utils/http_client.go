package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func Request(method string, url string, payload io.Reader) error {
	var (
		client      = new(http.Client)
		req         *http.Request
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*20)
		resp        *http.Response
		err         error
	)
	defer cancel()

	if req, err = http.NewRequestWithContext(ctx, method, url, payload); err != nil {
		return fmt.Errorf("error creating http request - %s", err)
	}

	if resp, err = client.Do(req); err != nil {
		return fmt.Errorf("error doing request %v - %s", req, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		return fmt.Errorf("non 2xx response status code")
	}

	return nil
}
