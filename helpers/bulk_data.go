package helpers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	errNewReq       = errors.New("creating request failed")
	errDoReq        = errors.New("doing request failed")
	errReadingBody  = errors.New("reading body from request failed")
	errCloseResBody = errors.New("closing response body failed")

	_defaultZincSearchURL   = "http://localhost:4080"
	_defaultZincSearchIndex = "emails"

	httpClient = &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
			DisableKeepAlives:   false,
			DisableCompression:  false,
			ForceAttemptHTTP2:   true,
		},
	}

	admin    string
	password string

	zincSearchURL string

	setCredentialsOnce sync.Once
)

func setCredentials() {
	baseURL := os.Getenv("ZINCSEARCH_URL")
	if baseURL == "" {
		baseURL = _defaultZincSearchURL
	}

	// _bulk takes action metadata (including _index and _id) in the request
	// body, so the index name is no longer part of the URL path. The index
	// name lives in main.IndexName and is emitted per-document by the caller.
	zincSearchURL = fmt.Sprintf("%s/api/_bulk", baseURL)

	admin = os.Getenv("ZINCSEARCH_USERNAME")
	password = os.Getenv("ZINCSEARCH_PASSWORD")
}

// IndexName returns the ZincSearch index to write to, sourced from the
// ZINCSEARCH_INDEX env var with a sensible default. Exposed so that callers
// can construct per-document action lines for the _bulk endpoint.
func IndexName() string {
	if v := os.Getenv("ZINCSEARCH_INDEX"); v != "" {
		return v
	}
	return _defaultZincSearchIndex
}

// BulkData indexes the data to the database
func BulkData(query string) error {
	if query == "" {
		return nil
	}

	setCredentialsOnce.Do(setCredentials)

	req, err := http.NewRequest(http.MethodPost, zincSearchURL, strings.NewReader(query))
	if err != nil {
		return fmt.Errorf("%w: %v", errNewReq, err)
	}

	req.SetBasicAuth(admin, password)
	req.Header.Set("Content-Type", "application/x-ndjson")
	req.Header.Set("Connection", "keep-alive")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", errDoReq, err)
	}
	defer closeResponseBody(resp)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%w: %v", errReadingBody, err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("bulk request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if resp.StatusCode == http.StatusOK {
		log.Printf("Bulk upload successful: %s", string(body))
	}

	return nil
}

func closeResponseBody(response *http.Response) {
	err := response.Body.Close()
	if err != nil {
		fmt.Println(errCloseResBody, err)
	}
}
