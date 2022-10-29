package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strings"
)

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

var emailFields []string = []string{
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

func DirectoryReader(root string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}

func DirectoryChecker(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

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

func CheckElementInSlice(slice []string, field string) string {
	var result string = "F"
	for _, x := range slice {
		if strings.Contains(x, field) {
			result = "T"
			break
		}
	}
	return result
}

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

func ReadAndApplyEmailFormat(root string) Email {
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

func WriteEmailInJDSON(fullEmail Email) string {
	jsonEmail, _ := json.Marshal(fullEmail)
	myDirectory, err := os.Getwd()
	fullDirectory := myDirectory + "/emails1.ndjson"
	if err != nil {
		fmt.Println(err)
	}
	if _, err := os.Stat(fullDirectory); err == nil {
		// path/to/whatever exists
		file, err := os.OpenFile(fullDirectory, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
		}
		data := "\n" + `{ "index" : { "_index" : "myindex" } }` + "\n" + string(jsonEmail)

		defer file.Close()

		_, err2 := file.WriteString(data)

		if err2 != nil {
			log.Fatal(err2)
		}
		//BulkData(data)
		return "Done!"
	}

	file, err1 := os.Create("emails1.ndjson")

	if err1 != nil {
		log.Fatal(err1)
	}
	data := `{ "index" : { "_index" : "myindex" } }` + "\n" + string(jsonEmail)

	defer file.Close()

	_, err2 := file.WriteString(data)

	if err2 != nil {
		log.Fatal(err2)
	}
	//BulkData(data)
	return "Done!"
}

func FileChecker(root string, files []string) string {
	for _, file := range files {
		fileRoot := root + "/" + file
		directoryCheck, _ := DirectoryChecker(fileRoot)
		if !directoryCheck {
			fmt.Println(fileRoot)
			fullEmail := ReadAndApplyEmailFormat(fileRoot)
			confirmationMessage := WriteEmailInJDSON(fullEmail)
			fmt.Println(confirmationMessage)
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

func BulkData(query string) string {
	payload := strings.NewReader(query)

	req, err := http.NewRequest("POST", "http://localhost:4080/api/_bulk", payload)
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
	return string(body)
}

func main() {

	root := "../../Enron emails/enron_mail_20110402/maildir"

	files, err := DirectoryReader(root)
	if err != nil {
		log.Fatal(err)
	}

	message := FileChecker(root, files)
	fmt.Println(message)
}
