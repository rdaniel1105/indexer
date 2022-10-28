package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/mail"
	"os"
	"strings"
)

type Email struct {
	MessageID               string
	Date                    string
	From                    string
	To                      string
	Subject                 string
	Cc                      string
	MimeVersion             string
	ContentType             string
	ContentTransferEncoding string
	Bcc                     string
	XFrom                   string
	XTo                     string
	Xcc                     string
	Xbcc                    string
	XFolder                 string
	XOrigin                 string
	XFileName               string
	Body                    string
}

func IOReadDir(root string) ([]string, error) {
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

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func headerNumber(lines []string, fields []string) (int, []string) {
	Counter := 0
	b := make([]string, 0)
	for i, line := range lines {
		for j := 0; j < len(fields); j++ {
			if strings.Contains(line, fields[j]) {
				check := checkSlice(b, fields[j])
				if check == "F" {
					b = append(b, line)
					Counter = i
					break
				}
			}
		}
	}
	//fmt.Println(Counter)
	//fmt.Println(b)
	return Counter, b
}

func checkSlice(slice []string, field string) string {
	var result string = "F"
	for _, x := range slice {
		if strings.Contains(x, field) {
			result = "T"
			break
		}
	}
	return result
}

func emailHeaderCheck(body string, fields []string) string {
	fmt.Println("entro")
	finalString := ""

	lines := strings.Split(strings.TrimRight(body, "\n"), "\n")
	Counter, b := headerNumber(lines, fields)
	for i, line := range lines {
		check := checkSlice(b, line)
		if check == "F" && i < Counter && len(line) != 0 {
			line = strings.Replace(line, "", " ", 1)
			if i == 0 {
				finalString += line
				break
			}
			finalString += line + "\n"
		} else {
			finalString += line + "\n"
		}

	}
	//fmt.Println(finalString)
	return finalString
}

func ReadAndApplyEmailFormat(root string, fields []string) Email {
	emailContent, err := ioutil.ReadFile(root)
	if err != nil {
		fmt.Println("File reading error", err)
	}

	correctedEmail := emailHeaderCheck(string(emailContent), fields)
	contentReader := strings.NewReader(correctedEmail)

	emailMessage, err := mail.ReadMessage(contentReader)
	if err != nil {
		fmt.Println("error aca")
		log.Fatal(err)
	}

	header := emailMessage.Header

	bodyReader, err := io.ReadAll(emailMessage.Body)
	if err != nil {
		log.Fatal(err)
	}
	enron := Email{header.Get(fields[0]),
		header.Get(fields[1]),
		header.Get(fields[2]),
		header.Get(fields[3]),
		header.Get(fields[4]),
		header.Get(fields[5]),
		header.Get(fields[6]),
		header.Get(fields[7]),
		header.Get(fields[8]),
		header.Get(fields[9]),
		header.Get(fields[10]),
		header.Get(fields[11]),
		header.Get(fields[12]),
		header.Get(fields[13]),
		header.Get(fields[14]),
		header.Get(fields[15]),
		header.Get(fields[16]),
		string(bodyReader)}
	return enron
}

func WriteEmailInJDSON(enron Email) string {
	b, _ := json.Marshal(enron)
	mydir, err := os.Getwd()
	fullDirectory := mydir + "/emails.jdson"
	if err != nil {
		fmt.Println(err)
	}
	if _, err := os.Stat(fullDirectory); err == nil {
		// path/to/whatever exists
		f, err := os.OpenFile(fullDirectory, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
		}
		data := "\n" + `{ "index" : { "_index" : "myindex" } }` + "\n" + string(b)

		defer f.Close()

		_, err2 := f.WriteString(data)

		if err2 != nil {
			log.Fatal(err2)
		}
		return "Done!"
	} else {
		f, err1 := os.Create("emails.jdson")

		if err1 != nil {
			log.Fatal(err1)
		}
		data := `{ "index" : { "_index" : "myindex" } }` + "\n" + string(b)

		defer f.Close()

		_, err2 := f.WriteString(data)

		if err2 != nil {
			log.Fatal(err2)
		}
		return "Done!"
	}
}

func FileChecker(root string, fields []string, files []string) string {
	for _, file := range files {
		fileRoot := root + "/" + file
		dirCheck, _ := isDirectory(fileRoot)
		if !dirCheck {
			fmt.Println(fileRoot)
			//checkedEmail := emailHeaderCheck(fileRoot, fields)
			fullEmail := ReadAndApplyEmailFormat(fileRoot, fields)
			message := WriteEmailInJDSON(fullEmail)
			fmt.Println(message)
		} else {
			subFiles, err := IOReadDir(fileRoot)
			if err != nil {
				log.Fatal(err)
			}
			FileChecker(fileRoot, fields, subFiles)
		}
	}
	return "All files done!"
}

func main() {

	root := "../../Enron emails/enron_mail_20110402/maildir"

	files, err := IOReadDir(root)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(files)

	fields := []string{
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

	message := FileChecker(root, fields, files)
	fmt.Println(message)

}
