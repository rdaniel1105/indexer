package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"

	"example/indexer/helpers"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	memprofile = flag.String("memprofile", "", "write memory profile to `file`")

	wg sync.WaitGroup

	emailSender = make(chan string, 20000)
	max         = 100
	emailsChunk string
)

// FileChecker checks if the directory contains either a file or another directory
// in case it finds a file, it will send it for it to be opened and read, then it will
// be written in ndjson format and returned to stack it and send it in a chunk.
func FileChecker(root string, files []string) string {
	for _, file := range files {
		fileRoot := root + "/" + file

		directoryCheck, err := helpers.DirectoryChecker(fileRoot)
		if err != nil {
			log.Fatal(err)
		}

		if !directoryCheck {
			fmt.Println(fileRoot)

			fullEmail, repeatedEmail := helpers.ReadAndCreateEmailStruct(fileRoot)
			if repeatedEmail {
				fmt.Println("Repeated!")
				continue
			}

			emailSender <- helpers.WriteEmailInNDJSON(fullEmail)

			fmt.Println("Done!")
		} else {
			subDirectories, err := helpers.DirectoryReader(fileRoot)
			if err != nil {
				log.Fatal(err)
			}

			FileChecker(fileRoot, subDirectories)
		}
	}

	return "All files done!"
}

// EmailsChunkSender receives the emails and then sends them ready to bulk in chunks of 100.
func EmailsChunkSender() {

	for {
		_, open := <-emailSender
		emailsChunk = ""

		for j := 0; j < max; j++ {
			if !open {
				wg.Add(1)

				go func() {
					defer wg.Done()
					helpers.BulkData(emailsChunk)
				}()

				break
			}

			emailsChunk += "\n" + <-emailSender
		}

		if !open {
			break
		}

		wg.Add(1)

		go func() {
			defer wg.Done()
			helpers.BulkData(emailsChunk)
		}()

	}

}

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}

		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}

		defer pprof.StopCPUProfile()
	}

	if len(os.Args) == 1 {
		log.Fatal("No files to process")
		return
	}

	root := os.Args[1] + "/maildir"

	files, err := helpers.DirectoryReader(root)
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		EmailsChunkSender()
	}()

	message := FileChecker(root, files)

	close(emailSender)

	fmt.Println(message)

	wg.Wait()

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}

		runtime.GC()
		if err = pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}

		f.Close()
		if err = pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not close the file: ", err)
		}
	}
}
