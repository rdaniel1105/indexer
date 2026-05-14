package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/rdaniel1105/indexer/models"
	"fmt"
	"io"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
)

var (
	emailHashes = make(map[string]struct{})

	hasher = sha256.New()
)

func hashContent(content string) string {
	hasher.Write([]byte(content))

	return hex.EncodeToString(hasher.Sum(nil))
}

// CreateEmailStruct reads the text file content and creates the corresponding email structure
func CreateEmailStruct(path string) (*models.Email, bool, error) {
	emailContent, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, false, fmt.Errorf("file reading error: %w", err)
	}

	correctedEmail := emailHeaderCheck(string(emailContent))
	contentReader := strings.NewReader(correctedEmail)

	emailMessage, err := mail.ReadMessage(contentReader)
	if err != nil {
		return nil, false, fmt.Errorf("mail message reading error: %w", err)
	}

	header := emailMessage.Header

	body, err := io.ReadAll(emailMessage.Body)
	if err != nil {
		return nil, false, fmt.Errorf("mail body reading error: %w", err)
	}

	email := &models.Email{
		MessageID:               header.Get("Message-ID"),
		Date:                    header.Get("Date"),
		From:                    header.Get("From"),
		To:                      header.Get("To"),
		Subject:                 header.Get("Subject"),
		Cc:                      header.Get("Cc"),
		MimeVersion:             header.Get("Mime-Version"),
		ContentType:             header.Get("Content-Type"),
		ContentTransferEncoding: header.Get("Content-Transfer-Encoding"),
		Bcc:                     header.Get("Bcc"),
		XFrom:                   header.Get("X-From"),
		XTo:                     header.Get("X-To"),
		Xcc:                     header.Get("X-cc"),
		Xbcc:                    header.Get("X-bcc"),
		XFolder:                 header.Get("X-Folder"),
		XOrigin:                 header.Get("X-Origin"),
		XFileName:               header.Get("X-FileName"),
		Body:                    string(body)}

	repeatedEmail := RepeatedEmailChecker(email.Body)
	if repeatedEmail {
		return nil, true, nil
	}

	return email, false, nil
}

// RepeatedEmailChecker ahora usa el hash del contenido
func RepeatedEmailChecker(newBody string) bool {
	hash := hashContent(newBody)

	_, exists := emailHashes[hash]
	if exists {
		return true
	}

	emailHashes[hash] = struct{}{}

	return false
}
