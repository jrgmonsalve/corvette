package main

type EmailData struct {
	Subject string `json:"subject"`
	From    string `json:"from"`
	To      string `json:"to"`
	Body    string `json:"body"`
}

type EmailDataBulk struct {
	Index   string      `json:"index"`
	Records []EmailData `json:"records"`
}
