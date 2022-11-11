package helpers

import (
	"encoding/json"
	"example/indexer/models"
)

// WriteEmailInNDJSON writes the email with the ndjson format to the file
func WriteEmailInNDJSON(fullEmail models.Email) string {
	jsonEmail, _ := json.Marshal(fullEmail)
	data := string(jsonEmail)

	// myDirectory, err := os.Getwd()
	// fullDirectory := myDirectory + "/emails1.ndjson"
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// if _, err := os.Stat(fullDirectory); err == nil {
	// 	file, err := os.OpenFile(fullDirectory, os.O_APPEND|os.O_WRONLY, 0644)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}

	// 	defer file.Close()

	// 	_, err = file.WriteString(data)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	//BulkData(data)
	// 	return "Done!"
	// }

	// file, err := os.Create("emails1.ndjson")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer file.Close()

	// _, err = file.WriteString(data)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	return data
}
