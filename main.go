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

	wg          sync.WaitGroup
	emailSender = make(chan string, 20000)
	done        = make(chan bool)
	// EmailReceiver = make(chan string, 20000)
	max = 20000
	i   = 0

	// emailsToUpload []string
)

// FileChecker checks if the directory contains either a file or another directory
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

			// emailsToUpload = append(emailsToUpload, helpers.WriteEmailInNDJSON(fullEmail))
			emailSender <- helpers.WriteEmailInNDJSON(fullEmail)

			fmt.Println("Done!")
		} else {
			subFiles, err := helpers.DirectoryReader(fileRoot)
			if err != nil {
				log.Fatal(err)
			}

			FileChecker(fileRoot, subFiles)
		}
	}

	return "All files done!"
}

// DataBulk asdasd
func DataBulk(done <-chan bool) {

	select {
	case <-done:
		break
	default:
		for i = 0; i < max; i++ {
			emailsChunk := ""

			for j := 0; j < max; j++ {
				emailsChunk += "\n" + <-emailSender
			}

			wg.Add(1)

			go func() {
				defer wg.Done()
				helpers.BulkData(emailsChunk)
			}()
			max += max
		}
		wg.Wait()
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

	root := "../enron_mail_20110402/maildir"
	// if len(os.Args) == 1 {
	// 	log.Fatal("No files to process")
	// 	return
	// }

	//root := os.Args[1]

	files, err := helpers.DirectoryReader(root)
	if err != nil {
		log.Fatal(err)
	}

	go DataBulk(done)

	message := FileChecker(root, files)
	close(emailSender)
	fmt.Println(message)

	close(done)

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}
}
