package helpers

import (
	"example/indexer/models"
	"fmt"
	"io"
	"io/ioutil"
	"net/mail"
	"path/filepath"
	"strings"
)

var (
	emailBodies []string
)

// CreateEmailStruct reads the text file content and creates the corresponding email structure.
func CreateEmailStruct(path string) (models.Email, bool, error) {
	emailContent, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return models.Email{}, false, fmt.Errorf("file reading error: %w", err)
	}

	correctedEmail := EmailHeaderCheck(string(emailContent))
	contentReader := strings.NewReader(correctedEmail)
	emailMessage, err := mail.ReadMessage(contentReader)
	if err != nil {
		return models.Email{}, false, fmt.Errorf("mail message reading error: %w", err)
	}

	header := emailMessage.Header

	body, err := io.ReadAll(emailMessage.Body)
	if err != nil {
		return models.Email{}, false, fmt.Errorf("mail body reading error: %w", err)
	}

	email := models.Email{
		MessageID:               header.Get(models.EmailFields[0]),
		Date:                    header.Get(models.EmailFields[1]),
		From:                    header.Get(models.EmailFields[2]),
		To:                      header.Get(models.EmailFields[3]),
		Subject:                 header.Get(models.EmailFields[4]),
		Cc:                      header.Get(models.EmailFields[5]),
		MimeVersion:             header.Get(models.EmailFields[6]),
		ContentType:             header.Get(models.EmailFields[7]),
		ContentTransferEncoding: header.Get(models.EmailFields[8]),
		Bcc:                     header.Get(models.EmailFields[9]),
		XFrom:                   header.Get(models.EmailFields[10]),
		XTo:                     header.Get(models.EmailFields[11]),
		Xcc:                     header.Get(models.EmailFields[12]),
		Xbcc:                    header.Get(models.EmailFields[13]),
		XFolder:                 header.Get(models.EmailFields[14]),
		XOrigin:                 header.Get(models.EmailFields[15]),
		XFileName:               header.Get(models.EmailFields[16]),
		Body:                    string(body)}

	repeatedEmail := RepeatedEmailChecker(email.Body)
	if repeatedEmail {
		return email, true, nil
	}

	return email, false, nil
}

// RepeatedEmailChecker checks if the emails has been added already, to avoid duplicity.
func RepeatedEmailChecker(newBody string) bool {
	for _, body := range emailBodies {
		if body == newBody {
			return true
		}
	}

	emailBodies = append(emailBodies, newBody)

	return false
}
