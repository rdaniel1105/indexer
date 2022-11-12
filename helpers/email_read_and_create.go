package helpers

import (
	"example/indexer/models"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/mail"
	"strings"
)

// ReadAndCreateEmailStruct reads the text file content and creates the corresponding structure.
func ReadAndCreateEmailStruct(root string) (models.Email, bool) {
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

	repeated := RepeatedEmailChecker(email.Body)
	if repeated {
		return email, true
	}

	return email, false
}
