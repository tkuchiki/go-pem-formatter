package pemformatter

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
)

func removeWhitespace(data string) string {
	re := regexp.MustCompile(`\s`)
	return re.ReplaceAllString(data, "")
}

func parseData(data string) (string, string, string, error) {
	re := regexp.MustCompile(`(?m)(-----BEGIN\s(?:.+)-----)(?:\s)?([\S\s]+)(?:\s)?(-----END\s(?:.+)-----)`)

	group := re.FindStringSubmatch(data)

	if len(group) < 1 {
		return "", "", "", fmt.Errorf("invalid data")
	}

	return group[1], group[2], group[3], nil
}

func format(header, body, footer string) string {
	var newline string
	if runtime.GOOS == "windows" {
		newline = "\r\n"
	} else {
		newline = "\n"
	}

	runes := []rune(removeWhitespace(body))
	pems := make([]string, 0)
	pems = append(pems, header)
	splitLen := 64
	runeLen := len(runes)
	for i := 0; i < runeLen; i += splitLen {
		if i+splitLen < len(runes) {
			pems = append(pems, string(runes[i:(i+splitLen)]))
		} else {
			pems = append(pems, string(runes[i:]))
		}
	}
	pems = append(pems, footer)

	return strings.Join(pems, newline)
}

func Format(data string) (string, error) {
	header, body, footer, err := parseData(data)
	if err != nil {
		return "", err
	}

	return format(header, body, footer), nil
}
