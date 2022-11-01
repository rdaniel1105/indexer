package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/joho/godotenv"
)

// Email struct is the structure the emails will have
type Email struct {
	MessageID               string `json:"Message-ID"`
	Date                    string
	From                    string
	To                      string
	Subject                 string
	Cc                      string
	MimeVersion             string `json:"Mime-Version"`
	ContentType             string `json:"Content-Type"`
	ContentTransferEncoding string `json:"Content-Transfer-Encoding"`
	Bcc                     string
	XFrom                   string `json:"X-From"`
	XTo                     string `json:"X-To"`
	Xcc                     string `json:"X-cc"`
	Xbcc                    string `json:"X-bcc"`
	XFolder                 string `json:"X-Folder"`
	XOrigin                 string `json:"X-Origin"`
	XFileName               string `json:"X-FileName"`
	Body                    string
}

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	memprofile = flag.String("memprofile", "", "write memory profile to `file`")

	emailFields []string = []string{
		"Message-ID",
		"Date",
		"From",
		"To",
		"Subject",
		"Cc",
		"Mime-Version",
		"Content-Type",
		"Content-Transfer-Encoding",
		"Bcc",
		"X-From",
		"X-To",
		"X-cc",
		"X-bcc",
		"X-Folder",
		"X-Origin",
		"X-FileName",
	}
)

// DirectoryReader reads the content inside the given path.
func DirectoryReader(root string) ([]string, error) {
	files := make([]string, 0)

	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}

	return files, nil
}

// DirectoryChecker checks if the path given is a directory.
func DirectoryChecker(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}

// GetHeaderNumber gets the correct number of headers using CheckElementInSlice() to make sure they're not repeated.
func GetHeaderNumber(lines []string) (int, []string) {
	counter := 0
	emailFieldSlice := make([]string, 0)

	for i, line := range lines {
		for j := 0; j < len(emailFields); j++ {
			if !strings.Contains(line, emailFields[j]) {
				continue
			}

			check := CheckElementInSlice(emailFieldSlice, emailFields[j])
			if check == "F" {
				emailFieldSlice = append(emailFieldSlice, line)
				counter = i
				break
			}
		}
	}

	return counter, emailFieldSlice
}

// CheckElementInSlice checks if a header is being repeated in the email body.
func CheckElementInSlice(slice []string, field string) string {
	result := "F"
	for _, x := range slice {
		if strings.Contains(x, field) {
			result = "T"
			break
		}
	}

	return result
}

// EmailHeaderCheck checks email format (in case of having multiline headers)
func EmailHeaderCheck(body string) string {
	correctedEmail := ""
	lines := strings.Split(strings.TrimRight(body, "\n"), "\n")
	counter, emailFieldSlice := GetHeaderNumber(lines)

	for i, line := range lines {
		check := CheckElementInSlice(emailFieldSlice, line)
		if check == "F" && i < counter && len(line) != 0 {
			line = strings.Replace(line, "", " ", 1)
			if i == 0 {
				correctedEmail += line
				break
			}
			correctedEmail += line + "\n"
		} else {
			correctedEmail += line + "\n"
		}

	}

	return correctedEmail
}

// ReadAndCreateEmailStruct reads the text file content and creates the corresponding structure.
func ReadAndCreateEmailStruct(root string) Email {
	emailContent, err := ioutil.ReadFile(root)
	if err != nil {
		fmt.Println("File reading error", err)
	}

	correctedEmail := EmailHeaderCheck(string(emailContent))
	contentReader := strings.NewReader(correctedEmail)
	emailMessage, err := mail.ReadMessage(contentReader)
	if err != nil {
		fmt.Println("error aca")
		log.Fatal(err)
	}

	header := emailMessage.Header

	body, err := io.ReadAll(emailMessage.Body)
	if err != nil {
		log.Fatal(err)
	}
	email := Email{header.Get(emailFields[0]),
		header.Get(emailFields[1]),
		header.Get(emailFields[2]),
		header.Get(emailFields[3]),
		header.Get(emailFields[4]),
		header.Get(emailFields[5]),
		header.Get(emailFields[6]),
		header.Get(emailFields[7]),
		header.Get(emailFields[8]),
		header.Get(emailFields[9]),
		header.Get(emailFields[10]),
		header.Get(emailFields[11]),
		header.Get(emailFields[12]),
		header.Get(emailFields[13]),
		header.Get(emailFields[14]),
		header.Get(emailFields[15]),
		header.Get(emailFields[16]),
		string(body)}
	return email
}

// WriteEmailInNDJSON writes the email with the ndjson format to the file
func WriteEmailInNDJSON(fullEmail Email) string {
	jsonEmail, _ := json.Marshal(fullEmail)
	// myDirectory, err := os.Getwd()
	// fullDirectory := myDirectory + "/emails1.ndjson"
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// if _, err := os.Stat(fullDirectory); err == nil {
	// 	// path/to/whatever exists
	// 	file, err := os.OpenFile(fullDirectory, os.O_APPEND|os.O_WRONLY, 0644)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	data := "\n" + `{ "index" : { "_index" : "myindex" } }` + "\n" + string(jsonEmail)

	// 	defer file.Close()

	// 	_, err2 := file.WriteString(data)
	// 	if err2 != nil {
	// 		log.Fatal(err2)
	// 	}
	// 	//BulkData(data)
	// 	return "Done!"
	// }

	// file, err1 := os.Create("emails1.ndjson")

	// if err1 != nil {
	// 	log.Fatal(err1)
	// }
	// data := `{ "index" : { "_index" : "myindex" } }` + "\n" + string(jsonEmail)

	// defer file.Close()

	// _, err2 := file.WriteString(data)

	// if err2 != nil {
	// 	log.Fatal(err2)
	// }
	data := `{ "index" : { "_index" : "myindex" } }` + "\n" + string(jsonEmail)

	return data
}

// FileChecker checks if the directory contains either a file or another directory
func FileChecker(root string, files []string) string {
	for _, file := range files {
		fileRoot := root + "/" + file
		directoryCheck, _ := DirectoryChecker(fileRoot)
		if !directoryCheck {
			fmt.Println(fileRoot)
			fullEmail := ReadAndCreateEmailStruct(fileRoot)
			data := WriteEmailInNDJSON(fullEmail)
			BulkData(data)
			fmt.Println("Done!")
		} else {
			subFiles, err := DirectoryReader(fileRoot)
			if err != nil {
				log.Fatal(err)
			}
			FileChecker(fileRoot, subFiles)
		}
	}
	return "All files done!"
}

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

	files, err := DirectoryReader(root)
	if err != nil {
		log.Fatal(err)
	}

	message := FileChecker(root, files)
	fmt.Println(message)

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
