package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

func main() {

	flag.Parse()
	rootPath := flag.Arg(0)
	if rootPath == "" {
		fmt.Println("Debe proporcionar una ruta de inicio.")
		os.Exit(1)
	}

	var wg sync.WaitGroup
	rootPath = rootPath + "/maildir"
	entries, err := os.ReadDir(rootPath)
	if err != nil {
		fmt.Printf("Error al leer el directorio: %v\n", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			wg.Add(1)
			go func(p string) {
				defer wg.Done()
				processDirectory(p)
			}(rootPath + "/" + entry.Name())
		}
	}

	if err != nil {
		fmt.Printf("Error al leer la ruta: %v\n", err)
		os.Exit(1)
	}

	wg.Wait()
}

type EmailData struct {
	Subject string `json:"subject"`
	From    string `json:"from"`
	To      string `json:"to"`
	Body    string `json:"body"`
}

func processDirectory(path string) {
	fmt.Printf("Procesando directorio: %s\n", path)

	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error al leer el directorio: %v\n", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			filePath := filepath.Join(path, entry.Name())
			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error al leer el archivo: %v\n", err)
				continue
			}
			emailData := parseEmailData(string(fileContent))
			jsonData, err := json.MarshalIndent(emailData, "", "    ")
			if err != nil {
				fmt.Println("Error al convertir a JSON:", err)
				return
			}
			sendToZincSearch(jsonData)
		} else {
			processDirectory(filepath.Join(path, entry.Name()))
		}
	}
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

const zincURL = "http://localhost:4080/api/emails/_doc"

func sendToZincSearch(jsonData []byte) error {
	req, err := http.NewRequest("POST", zincURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+basicAuth("", ""))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ZincSearch respondió con código de estado: %d", resp.StatusCode)
	}

	return nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
