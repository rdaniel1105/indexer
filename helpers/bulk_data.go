package helpers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// BulkData indexes the data to the database
func BulkData(query string) {

	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	payload := strings.NewReader(query)
	zincSearchURL := os.Getenv("ZincSearchURL")

	req, err := http.NewRequest(http.MethodPost, zincSearchURL, payload)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth("admin", "Complexpass#123")
	req.Header.Set("Content-Type", "application-ndjson")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	log.Println(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Uploaded!")
	fmt.Println(string(body))
}
