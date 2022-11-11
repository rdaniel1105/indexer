package helpers

import (
	"encoding/json"
	"example/indexer/models"
)

// WriteEmailInNDJSON writes the email with the ndjson format to the file
func WriteEmailInNDJSON(fullEmail models.Email) string {
	jsonEmail, _ := json.Marshal(fullEmail)
	// myDirectory, err := os.Getwd()
	// fullDirectory := myDirectory + "/emails1.ndjson"
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// if _, err := os.Stat(fullDirectory); err == nil {
	// 	// path/to/whatever exists
	// 	file, err := os.OpenFile(fullDirectory, os.O_APPEND|os.O_WRONLY, 0644)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	data := "\n" + `{ "index" : { "_index" : "myindex" } }` + "\n" + string(jsonEmail)

	// 	defer file.Close()

	// 	_, err2 := file.WriteString(data)
	// 	if err2 != nil {
	// 		log.Fatal(err2)
	// 	}
	// 	//BulkData(data)
	// 	return "Done!"
	// }

	// file, err1 := os.Create("emails1.ndjson")

	// if err1 != nil {
	// 	log.Fatal(err1)
	// }
	// data := `{ "index" : { "_index" : "myindex" } }` + "\n" + string(jsonEmail)

	// defer file.Close()

	// _, err2 := file.WriteString(data)

	// if err2 != nil {
	// 	log.Fatal(err2)
	// }
	data := `{ "index" : { "_index" : "myindex" } }` + "\n" + string(jsonEmail)

	return data
}
