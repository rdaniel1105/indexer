package helpers

import (
	"example/indexer/models"
	"fmt"
	"log"
	"strings"
)

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

// GetHeaderNumber gets the correct number of headers using CheckElementInSlice() to make sure they're not repeated.
func GetHeaderNumber(lines []string) (int, []string) {
	counter := 0
	emailFieldSlice := make([]string, 0)

	for i, line := range lines {
		for j := 0; j < len(models.EmailFields); j++ {
			if !strings.Contains(line, models.EmailFields[j]) {
				continue
			}

			check := CheckElementInSlice(emailFieldSlice, models.EmailFields[j])
			if check == "F" {
				emailFieldSlice = append(emailFieldSlice, line)
				counter = i
				break
			}
		}
	}

	return counter, emailFieldSlice
}