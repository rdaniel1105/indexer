package models

// Email struct is the structure the emails will have
type Email struct {
	MessageID               string `json:"Message-ID"`
	Date                    string `json:"Date"`
	From                    string `json:"From"`
	To                      string `json:"To"`
	Subject                 string `json:"Subject"`
	Cc                      string `json:"Cc"`
	MimeVersion             string `json:"Mime-Version"`
	ContentType             string `json:"Content-Type"`
	ContentTransferEncoding string `json:"Content-Transfer-Encoding"`
	Bcc                     string `json:"Bcc"`
	XFrom                   string `json:"X-From"`
	XTo                     string `json:"X-To"`
	Xcc                     string `json:"X-cc"`
	Xbcc                    string `json:"X-bcc"`
	XFolder                 string `json:"X-Folder"`
	XOrigin                 string `json:"X-Origin"`
	XFileName               string `json:"X-FileName"`
	Body                    string `json:"Body"`
}

// EmailFields represents the fields that the email will contain.
var EmailFields = []string{
	"Message-ID",
	"Date",
	"From",
	"To",
	"Subject",
	"Cc",
	"Mime-Version",
	"Content-Type",
	"Content-Transfer-Encoding",
	"Bcc",
	"X-From",
	"X-To",
	"X-cc",
	"X-bcc",
	"X-Folder",
	"X-Origin",
	"X-FileName",
}
