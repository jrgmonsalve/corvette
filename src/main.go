package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type EmailData struct {
	Subject string `json:"subject"`
	From    string `json:"from"`
	To      string `json:"to"`
	Body    string `json:"body"`
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env:", err)
		return
	}

	flag.Parse()
	rootPath := flag.Arg(0)
	if rootPath == "" {
		log.Println("You must specify the path.")
		os.Exit(1)
	}
	c := make(chan int, 5)
	var wg sync.WaitGroup
	rootPath = rootPath + "/maildir"

	entries, err := os.ReadDir(rootPath)
	if err != nil {
		log.Printf("Fail to read directory: %v\n", err)
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			c <- 1
			wg.Add(1)
			fullPath := filepath.Join(rootPath, entry.Name())
			go processDirectory(fullPath, &wg, c)
		}
	}

	wg.Wait()
}

func processDirectory(path string, wg *sync.WaitGroup, c chan int) {
	log.Printf("Processing directory: %s\n", path)
	defer wg.Done()
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Printf("Fail to read directory: %v\n", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			fullFilePath := filepath.Join(path, entry.Name())
			log.Println("Processing file: ", fullFilePath)
			fileContent, err := os.ReadFile(fullFilePath)
			if err != nil {
				log.Printf("Fail to read file: %v\n", err)
				continue
			}
			emailData := parseEmailData(string(fileContent))
			jsonData, err := json.MarshalIndent(emailData, "", "    ")
			if err != nil {
				log.Println("Error marshaling JSON:", err)
				return
			}
			err = sendToZincSearch(jsonData)
			//TODO: retry depending on error
			if err != nil {
				log.Println("Error sending data to ZincSearch:", err)
				return
			}
		} else {
			processDirectory(filepath.Join(path, entry.Name()), wg, c)
		}
	}
	<-c
}

func parseEmailData(content string) EmailData {
	var emailData EmailData

	re := regexp.MustCompile(`X-FileName:.*\n`)

	match := re.FindStringIndex(content)
	if match == nil {
		fmt.Println("Tag 'X-FileName:' no encontrado.")
		return emailData
	}

	startIdx := match[0]
	endIdx := match[1]

	if match == nil {
		fmt.Println("Tag 'X-FileName:' no encontrado.")
		return emailData
	}

	reSubject := regexp.MustCompile(`Subject: (.+)`)
	reFrom := regexp.MustCompile(`From: (.+)`)
	reTo := regexp.MustCompile(`To: (.+)`)

	emailData.Subject = findValue(reSubject, content[:startIdx])
	emailData.From = findValue(reFrom, content[:startIdx])
	emailData.To = findValue(reTo, content[:startIdx])

	emailData.Body = strings.TrimSpace(content[endIdx:])

	return emailData
}

func findValue(re *regexp.Regexp, content string) string {
	match := re.FindStringSubmatch(content)
	if len(match) >= 2 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

func sendToZincSearch(jsonData []byte) error {
	apiURL := os.Getenv("API_URL")
	apiUser := os.Getenv("API_USER")
	apiPassword := os.Getenv("API_PASSWORD")

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+basicAuth(apiUser, apiPassword))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {

		return fmt.Errorf("fail to send data. unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
