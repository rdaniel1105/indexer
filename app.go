package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/mail"
	"os"
	"strings"

	// "io"

	// "io/ioutil"
	"log"
	//"net/mail"
	//"strings"
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

func ReadAndApplyEmailFormat(root string, fields []string) Email {
	emailContent, err := ioutil.ReadFile(root)
	if err != nil {
		fmt.Println("File reading error", err)
	}
	contentReader := strings.NewReader(string(emailContent))
	emailMessage, err := mail.ReadMessage(contentReader)
	if err != nil {
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
	if err != nil {
		fmt.Println(err)
	}
	if _, err := os.Stat(mydir); err == nil {
		// path/to/whatever exists
		f, err := os.OpenFile("emails.jdson", os.O_APPEND|os.O_WRONLY, 0644)
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
	}
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

func FileChecker(root string, fields []string, files []string) string {
	dirCheck, _ := isDirectory(root)
	if !dirCheck {
		fullEmail := ReadAndApplyEmailFormat(root, fields)
		message := WriteEmailInJDSON(fullEmail)
		fmt.Println(message)
	} else {
		subFiles, err := IOReadDir(root)
		if err != nil {
			log.Fatal(err)
		}
		for _, subFile := range subFiles {
			fmt.Println(subFile)
			subRoot := root + "/" + sub
			txtCheck2, _ := isDirectory(subRoot)

		}
	}
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

	for _, file := range files {
		subRoot := root + "/" + file
		dirCheck, _ := isDirectory(subRoot)
		if !dirCheck {
			fullEmail := ReadAndApplyEmailFormat(subRoot, fields)
			message := WriteEmailInJDSON(fullEmail)
			fmt.Println(message)
		} else {
			subFiles, err := IOReadDir(subRoot)
			if err != nil {
				log.Fatal(err)
			}
			for _, sub := range subFiles {
				subRoot1 := subRoot + "/" + sub
				txtCheck1, _ := isDirectory(subRoot1)
				if !txtCheck1 {
					fullEmail := ReadAndApplyEmailFormat(subRoot1, fields)
					message := WriteEmailInJDSON(fullEmail)
					fmt.Println(message)
				} else {
					subFiles1, err := IOReadDir(subRoot1)
					if err != nil {
						log.Fatal(err)
					}
					for _, sub1 := range subFiles1 {
						fmt.Println(sub1)
						subRoot2 := subRoot1 + "/" + sub1
						txtCheck2, _ := isDirectory(subRoot1)
						if !txtCheck2 {
							fullEmail := ReadAndApplyEmailFormat(subRoot2, fields)
							message := WriteEmailInJDSON(fullEmail)
							fmt.Println(message)
						}
						//else
					}
				}
			}
			fmt.Println(subFiles)
		}
	}

}
