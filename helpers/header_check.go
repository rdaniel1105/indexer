package helpers

import (
	"example/indexer/models"
	"strings"
)

// CheckElementInSlice checks if a header is being repeated in the email body.
func CheckElementInSlice(emailFieldSlice []string, field string) bool {
	result := false
	for _, x := range emailFieldSlice {
		if strings.Contains(x, field) {
			result = true

			break
		}
	}

	return result
}

// EmailHeaderCheck checks email format (in case of having multiline headers)
func EmailHeaderCheck(body string) string {
	correctedEmail := ""
	lines := strings.Split(strings.TrimRight(body, "\n"), "\n")

	numOfHeaders, emailFieldSlice := GetHeaderNumber(lines)

	for i, line := range lines {
		check := CheckElementInSlice(emailFieldSlice, line)

		if i == 0 {
			correctedEmail += line
			continue
		}

		if !check && i < numOfHeaders && len(line) != 0 {
			line = strings.Replace(line, "", " ", 1)
		}

		correctedEmail += line + "\n"
	}

	return correctedEmail
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
			if !check {
				emailFieldSlice = append(emailFieldSlice, line)
				counter = i
				break
			}
		}
	}

	return counter, emailFieldSlice
}
