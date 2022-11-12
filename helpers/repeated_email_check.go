package helpers

var (
	emailBodies []string
)

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
