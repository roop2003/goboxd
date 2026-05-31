//go:build load

package load_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/thesouldev/goboxd/server"
)

func TestHealthzLoad(t *testing.T) {
	testServer := httptest.NewServer(server.NewMux())
	defer testServer.Close()

	client := &http.Client{Timeout: 2 * time.Second}
	errs := make(chan error, 100)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			resp, err := client.Get(testServer.URL + "/healthz")
			if err != nil {
				errs <- err
				return
			}
			defer resp.Body.Close()

			if _, err := io.Copy(io.Discard, resp.Body); err != nil {
				errs <- err
				return
			}
			if resp.StatusCode != http.StatusOK {
				errs <- fmt.Errorf("status code = %d, want %d", resp.StatusCode, http.StatusOK)
			}
		}()
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Fatal(err)
	}
}
