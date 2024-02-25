package domain

import (
	"errors"
	"regexp"
	"strings"
)

type Email struct {
	Subject string `json:"subject"`
	From    string `json:"from"`
	To      string `json:"to"`
	Body    string `json:"body"`
}

func (email *Email) MappingFromString(content string) error {

	re := regexp.MustCompile(`X-FileName:.*\n`)
	match := re.FindStringIndex(content)
	if match == nil {
		return errors.New("tag 'X-FileName:' not found")
	}

	startIdx := match[0]
	endIdx := match[1]

	reSubject := regexp.MustCompile(`Subject: (.+)`)
	reFrom := regexp.MustCompile(`From: (.+)`)
	reTo := regexp.MustCompile(`To: (.+)`)

	email.Subject = findValue(reSubject, content[:startIdx])
	email.From = findValue(reFrom, content[:startIdx])
	email.To = findValue(reTo, content[:startIdx])

	email.Body = strings.TrimSpace(content[endIdx:])

	return nil
}

func findValue(re *regexp.Regexp, content string) string {
	match := re.FindStringSubmatch(content)
	if len(match) >= 2 {
		return strings.TrimSpace(match[1])
	}
	return ""
}
