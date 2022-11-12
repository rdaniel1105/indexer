package helpers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	errNewReq       = errors.New("NewRequest")
	errDoReq        = errors.New("Do(req)")
	errReadingBody  = errors.New("reading body from request")
	errCloseResBody = errors.New("closing response body")
)

// BulkData indexes the data to the database
func BulkData(query string) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	zincSearchURL := os.Getenv("ZincSearchURL")
	admin := os.Getenv("ADMIN")
	password := os.Getenv("PASSWORD")

	payload := strings.NewReader(query)

	req, err := http.NewRequest(http.MethodPost, zincSearchURL, payload)
	if err != nil {
		log.Fatal(errNewReq, err)
	}

	req.SetBasicAuth(admin, password)
	req.Header.Set("Content-Type", "application-ndjson")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(errDoReq, err)
	}

	defer closeResponseBody(resp)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(errReadingBody, err)
	}

	log.Println(resp.StatusCode)
	fmt.Println("Uploaded!")
	fmt.Println(string(body))
}

func closeResponseBody(response *http.Response) {
	err := response.Body.Close()
	if err != nil {
		fmt.Println(errCloseResBody, err)
	}
}
