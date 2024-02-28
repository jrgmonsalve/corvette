package api_zincsearch

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jrgmonsalve/corvette/cmd/cli/internal/domain"
	model_zincsearch "github.com/jrgmonsalve/corvette/cmd/cli/internal/infrastructure/models/zincsearch"
)

type ZincSearchService struct {
}

func NewZincSearchService() *ZincSearchService {
	return &ZincSearchService{}
}

func (zss *ZincSearchService) CreateBulk(emails []domain.Email) error {
	fmt.Println("Sending Creating bulk")
	fmt.Println("Emails:", len(emails))

	emailBulkRequest := model_zincsearch.EmailBulkRequest{
		IndexName: os.Getenv("INDEX_NAME"),
		Records:   emails,
	}
	jsonData := emailBulkRequest.MappingToJson()

	apiUser := os.Getenv("API_USER")
	apiPassword := os.Getenv("API_PASSWORD")
	apiURL := os.Getenv("API_URL") + "/_bulkv2"

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating request:", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+basicAuth(apiUser, apiPassword))

	// fmt.Println("Sending request to:", apiURL)
	// fmt.Println("Request body:", string(jsonData))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
	}
	defer resp.Body.Close()
	// fmt.Println("Response status:", resp.Status)
	if resp.StatusCode != http.StatusOK {
		log.Println("Error sending request:", resp.StatusCode)
	}
	senconds, _ := strconv.Atoi(os.Getenv("SECONDS_TO_SLEEP_BETWEEN_REQUEST"))
	time.Sleep(time.Duration(senconds) * time.Second)
	return nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
