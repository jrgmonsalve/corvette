package helpers

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/jrgmonsalve/corvette/cmd/cli/internal/domain"
)

func RecursiveFileReader(path string, currentDepth, maxDepth int, emails *[]domain.Email, batchEmailsSize int) error {

	patronFileName := regexp.MustCompile(`^\d+_$`)

	if maxDepth >= 0 && currentDepth > maxDepth {
		return nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		log.Println("Error reading directory: ", err)
		return errors.New("Error reading directory: " + path)
	}
	for _, entrie := range entries {
		fullPath := filepath.Join(path, entrie.Name())
		if entrie.IsDir() {
			RecursiveFileReader(fullPath, currentDepth+1, maxDepth, emails, batchEmailsSize)
		} else if patronFileName.MatchString(entrie.Name()) {
			log.Println("Processing file: ", path, entrie.Name())
			fileContent, err := os.ReadFile(fullPath)
			if err != nil {
				log.Printf("Fail to read file: %v\n", err)
				continue
			}
			email := domain.Email{}
			err = email.MappingFromString(string(fileContent))
			if err != nil {
				log.Printf("Fail to map email: %v\n", err)
				continue
			}
			*emails = append(*emails, email)

			if len(*emails) >= batchEmailsSize {
				sendDataToAPI(*emails)
				*emails = []domain.Email{}
			}
		}
	}
	if len(*emails) > 0 {
		sendDataToAPI(*emails)
		*emails = []domain.Email{}
	}
	return nil
}

func sendDataToAPI(emails []domain.Email) {
	fmt.Println("Sending emails to API")
}
