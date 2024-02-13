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
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func GetNodes(rootPath string) []os.DirEntry {
	var entries []os.DirEntry
	entries, err := os.ReadDir(rootPath)
	if err != nil {
		log.Printf("Fail to read directory: %v\n", err)
	}
	return entries
}

func GetRootPath() string {
	flag.Parse()
	rootPath := flag.Arg(0)
	if rootPath == "" {
		fmt.Println("Please provide a root path.")
		os.Exit(1)
	}
	return rootPath
}

func LoadEnvFile() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
}

func ProcessDirectory(path string, prefixLog string) {

	channelSize := getIntValueFromEnv("CHANNEL_SIZE", 10)

	c := make(chan int, channelSize)
	var wg sync.WaitGroup

	nodes := GetNodes(path)

	var batchEmailData []EmailData
	increment := 0

	for _, node := range nodes {
		increment++
		log.Println(prefixLog+"(files ", increment, "/", len(nodes), ")")
		fullPath := filepath.Join(path, node.Name())

		if !node.IsDir() {
			// log.Println("Processing file: ", fullPath)

			fileContent, err := os.ReadFile(fullPath)
			if err != nil {
				log.Printf("Fail to read file: %v\n", err)
				continue
			}
			emailData := parseEmailData(string(fileContent))
			batchEmailData = append(batchEmailData, emailData)

			if len(batchEmailData)%getIntValueFromEnv("BATCH_SIZE", 10) == 0 {
				// fmt.Println("into to requsts.............")
				bulk := EmailDataBulk{
					Index:   os.Getenv("INDEX_NAME"),
					Records: batchEmailData,
				}
				jsonData, err := json.Marshal(bulk)
				if err != nil {
					log.Println("Error marshaling JSON:", err)
					continue
				}
				c <- 1
				wg.Add(1)
				apiURL := os.Getenv("API_URL") + "/" + "_bulkv2"
				batchEmailData = batchEmailData[:0]
				go insertDocument(apiURL, jsonData, &wg, c)
			}

		} //else {
		// ProcessDirectory(fullPath)
		//}

	}
	wg.Wait()
}

func getIntValueFromEnv(key string, defaultValue int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		value = defaultValue
		fmt.Println("Fail to convert "+key+", set default value to:"+strconv.Itoa(defaultValue), err)
	}
	return value
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

func insertDocument(apiURL string, jsonData []byte, wg *sync.WaitGroup, c chan int) {
	defer wg.Done()

	apiUser := os.Getenv("API_USER")
	apiPassword := os.Getenv("API_PASSWORD")

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
	<-c
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
