package helpers

var (
	emailIDs []string
)

// RepeatedEmailChecker makes sure that email has not been created before, so it is not uploaded twice.
func RepeatedEmailChecker(id string) bool {

	for _, emailID := range emailIDs {
		if emailID == id {

			return true
		}
	}
	emailIDs = append(emailIDs, id)

	return false
}
