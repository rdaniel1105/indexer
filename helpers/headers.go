package helpers

import (
	"example/indexer/models"
	"strings"
)

// EmailHeaderCheck checks the email's format (in case of having multiline headers).
func EmailHeaderCheck(body string) string {
	correctedEmail := ""
	lines := strings.Split(strings.TrimRight(body, "\n"), "\n")

	numberOfHeaders, emailFieldSlice := GetHeaderNumber(lines)

	for i, line := range lines {
		check := CheckElementInSlice(emailFieldSlice, line)

		if i == 0 {
			correctedEmail += line
			continue
		}

		if !check && i < numberOfHeaders && len(line) != 0 {
			line = strings.Replace(line, "", " ", 1)
		}

		correctedEmail += line + "\n"
	}

	return correctedEmail
}

// GetHeaderNumber gets the correct number of headers using CheckElementInSlice() to make sure we do not re-set a
// header when having a header in the email body.
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

// CheckElementInSlice checks if a header is being repeated in the email body.
func CheckElementInSlice(emailFieldSlice []string, field string) bool {
	repeated := false
	for _, x := range emailFieldSlice {
		if strings.Contains(x, field) {
			repeated = true

			break
		}
	}

	return repeated
}
