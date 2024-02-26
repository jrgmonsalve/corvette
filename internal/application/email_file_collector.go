package application

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/jrgmonsalve/corvette/cmd/cli/internal/domain"
	"github.com/jrgmonsalve/corvette/cmd/cli/pkg/helpers"
)

type EmailFileCollector struct{}

func NewEmailFileCollector() *EmailFileCollector {
	return &EmailFileCollector{}
}

func (efr *EmailFileCollector) Collect(emailChannel chan<- domain.Email, wg *sync.WaitGroup) error {
	defer wg.Done()
	path, err := helpers.GetArgumentFromCLI("You must provide the root path")
	if err != nil {
		return err
	}

	maxDepthDirectory := helpers.GetIntValueFromEnv("MAX_DEPTH_DIRECTORY", 3)

	err = efr.recursiveFileReader(path, 0, maxDepthDirectory, emailChannel)
	if err != nil {
		return err
	}

	return nil
}

func (efr *EmailFileCollector) recursiveFileReader(path string, currentDepth, maxDepth int, emailChannel chan<- domain.Email) error {

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
			efr.recursiveFileReader(fullPath, currentDepth+1, maxDepth, emailChannel)
		} else if patronFileName.MatchString(entrie.Name()) {
			log.Println("Processing file: ", path, entrie.Name())
			fileContent, err := os.ReadFile(fullPath)
			if err != nil {
				log.Printf("Fail to read file: %v\n", err)
				continue
			}
			email := domain.Email{}
			err = email.MappingFromString(string(fileContent))
			emailChannel <- email
			if err != nil {
				log.Printf("Fail to map email: %v\n", err)
				continue
			}

		}
	}
	close(emailChannel)
	return nil
}
