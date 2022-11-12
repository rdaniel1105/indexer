package helpers

import (
	"encoding/json"
	"example/indexer/models"
	"fmt"
	"log"
	"os"
)

// WriteEmailInNDJSON writes the email with the ndjson format to the file
func WriteEmailInNDJSON(fullEmail models.Email) string {
	jsonEmail, _ := json.Marshal(fullEmail)
	data := string(jsonEmail)

	myDirectory, err := os.Getwd()
	filePath := myDirectory + "/emails.ndjson"
	if err != nil {
		fmt.Println(err)
	}

	if _, err := os.Stat(filePath); err == nil {
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
		}

		defer fileClose(file)

		data = "\n" + data

		_, err = file.WriteString(data)
		if err != nil {
			log.Fatal(err)
		}

		return data
	}

	file, err := os.Create("emails1.ndjson")
	if err != nil {
		log.Fatal(err)
	}

	defer fileClose(file)

	_, err = file.WriteString(data)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func fileClose(file *os.File) {
	err := file.Close()
	if err != nil {
		fmt.Println(errCloseResBody, err)
	}
}
