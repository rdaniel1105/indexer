package helpers

import (
	"github.com/rdaniel1105/indexer/models"
	"strings"
)

const (
	jumpLine        = "\n"
	headerSeparator = ":"
)

// emailHeaderCheck fixes the email format in case it has headers in multiple lines.
func emailHeaderCheck(body string) string {
	var builder strings.Builder

	lines := strings.Split(strings.TrimRight(body, jumpLine), jumpLine)

	// find the end of the headers
	headerEndIndex := findHeaderEnd(lines)

	for i, line := range lines {
		processedLine := processLine(line, i, headerEndIndex)

		builder.WriteString(processedLine)

		// add jump line except for the last line
		if i < len(lines)-1 {
			builder.WriteString(jumpLine)
		}
	}

	return builder.String()
}

// findHeaderEnd finds the index where the headers end
func findHeaderEnd(lines []string) int {
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			return i
		}
	}

	return len(lines)
}

// processLine processes a single line of the email
func processLine(line string, currentIndex, headerEndIndex int) string {
	// if we are after the headers, return the line without modification
	if currentIndex >= headerEndIndex {
		return line
	}

	// if the line already starts with space or tab, no modify
	if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
		return line
	}

	// find the header separator ':'
	colonIndex := strings.Index(line, headerSeparator)
	if colonIndex == -1 {
		return " " + line
	}

	// extract and verify the header name
	headerName := strings.TrimSpace(line[:colonIndex])
	if _, exists := models.EmailFieldsMap[headerName]; exists {
		return line
	}

	return " " + line
}
