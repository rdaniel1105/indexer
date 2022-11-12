package helpers

var (
	emailBodies []string
)

// RepeatedEmailChecker makes sure that email has not been created before, so it is not uploaded twice.
func RepeatedEmailChecker(newBody string) bool {

	for _, body := range emailBodies {
		if body == newBody {

			return true
		}
	}

	emailBodies = append(emailBodies, newBody)

	return false
}
