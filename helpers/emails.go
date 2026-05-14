package helpers

import (
	"fmt"
	"io"
	"net/mail"
	"os"
	"path/filepath"
	"strings"

	"github.com/rdaniel1105/indexer/models"
)

// CreateEmailStruct reads an RFC 822 file and returns the parsed Email.
// Deduplication is handled downstream by ZincSearch (each document is
// keyed by Message-ID, so re-indexing the same email replaces rather than
// duplicates the existing record).
func CreateEmailStruct(path string) (*models.Email, error) {
	emailContent, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("file reading error: %w", err)
	}

	correctedEmail := emailHeaderCheck(string(emailContent))
	contentReader := strings.NewReader(correctedEmail)

	emailMessage, err := mail.ReadMessage(contentReader)
	if err != nil {
		return nil, fmt.Errorf("mail message reading error: %w", err)
	}

	header := emailMessage.Header

	body, err := io.ReadAll(emailMessage.Body)
	if err != nil {
		return nil, fmt.Errorf("mail body reading error: %w", err)
	}

	return &models.Email{
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
		Body:                    string(body),
	}, nil
}
