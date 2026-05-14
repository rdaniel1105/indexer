package main

import (
	"encoding/json"
	"errors"
	"github.com/rdaniel1105/indexer/helpers"
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

const (
	chunkSize            = 5000
	maxChunkBytes        = 10 * 1024 * 1024 // 10 MB
	bufferSize           = 50000
	maxConcurrentSenders = 10

	jumpLineByte = 1

	jumpLine = "\n"
)

var (
	errGoDotenvLoad = errors.New("loading godotenv failed")

	indexName = "emails"
)

type chunkState struct {
	counter    int
	bytesCount int
	builder    strings.Builder
}

func (cs *chunkState) reset() {
	cs.counter = 0
	cs.bytesCount = 0
	cs.builder.Reset()
}

func initEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(errGoDotenvLoad, err)
	}

	indexName = helpers.IndexName()
}

func emailsChunkSender(emailSender chan string, wg *sync.WaitGroup) {
	chunkState := chunkState{
		counter:    0,
		bytesCount: 0,
		builder:    strings.Builder{},
	}

	limiter := make(chan struct{}, maxConcurrentSenders)

	for email := range emailSender {
		emailSize := len(email) + jumpLineByte

		if chunkState.bytesCount+emailSize > maxChunkBytes {
			sendChunk(getChunk(chunkState), wg, limiter)

			chunkState.reset()
		}

		chunkState.builder.WriteString(email)
		chunkState.builder.WriteString(jumpLine)

		chunkState.bytesCount += emailSize
		chunkState.counter++

		if chunkState.counter == chunkSize {
			sendChunk(getChunk(chunkState), wg, limiter)

			chunkState.reset()
		}
	}

	if chunkState.counter > 0 {
		sendChunk(getChunk(chunkState), wg, limiter) // no need to reset state, because it's the last chunk
	}
}

func getChunk(chunkState chunkState) string {
	chunk := chunkState.builder.String()

	return chunk
}

func sendChunk(chunk string, wg *sync.WaitGroup, limiter chan struct{}) {
	wg.Add(1)
	limiter <- struct{}{}

	go bulkUpload(chunk, wg, limiter)
}

func bulkUpload(data string, wg *sync.WaitGroup, limiter chan struct{}) {
	defer wg.Done()
	defer func() { <-limiter }()

	if err := helpers.BulkData(data); err != nil {
		log.Printf("bulk upload error: %v", err)

		return
	}
}

// fileChecker checks if the directory contains either a file or another directory
// in case it finds a file, it will process it and send it to the emailSender channel
func fileChecker(root string, files []string, emailSender chan string) {
	for _, file := range files {
		fileRoot := root + "/" + file

		directoryCheck, err := helpers.DirectoryChecker(fileRoot)
		if err != nil {
			fmt.Println(err)
		}

		if !directoryCheck {
			processEmail(fileRoot, emailSender)
			continue
		}

		subDirectories, err := helpers.DirectoryReader(fileRoot)
		if err != nil {
			fmt.Println(err)
		}

		fileChecker(fileRoot, subDirectories, emailSender)
	}
}

func processEmail(fileRoot string, emailSender chan string) {
	fullEmail, err := helpers.CreateEmailStruct(fileRoot)
	if err != nil {
		fmt.Println(err)

		return
	}

	jsonEmail, err := json.Marshal(fullEmail)
	if err != nil {
		fmt.Println(err)

		return
	}

	// Fall back to the file path when the email has no Message-ID header
	// (rare, but the Enron corpus contains some malformed messages).
	docID := fullEmail.MessageID
	if docID == "" {
		docID = fileRoot
	}

	actionJSON, err := json.Marshal(map[string]map[string]string{
		"index": {"_index": indexName, "_id": docID},
	})
	if err != nil {
		fmt.Println(err)

		return
	}

	emailSender <- string(actionJSON) + jumpLine + string(jsonEmail)
}

func main() {
	flag.Parse()

	dirPath := flag.Arg(0)
	if dirPath == "" {
		log.Fatal("usage: indexer <path-to-maildir>")
	}

	initEnv()

	files, err := helpers.DirectoryReader(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	emailSender := make(chan string, bufferSize)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		emailsChunkSender(emailSender, &wg)
	}()

	fileChecker(dirPath, files, emailSender)

	close(emailSender)
	wg.Wait()

	fmt.Printf("data indexed\n")
}
