package helpers

import (
	"encoding/json"
	"example/indexer/models"
)

// WriteEmailInNDJSON writes the email with the ndjson format to the file
func WriteEmailInNDJSON(fullEmail models.Email) string {
	jsonEmail, _ := json.Marshal(fullEmail)
	data := string(jsonEmail)

	data = "\n" + data

	return data
}
