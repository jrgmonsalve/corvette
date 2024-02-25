package application

import "fmt"

type ZincSearchUploader struct {
}

func NewZincSearchUploader() *ZincSearchUploader {
	return &ZincSearchUploader{}
}

func (zsu *ZincSearchUploader) Upload() error {
	fmt.Println("Start uploading to ZincSearch")
	fmt.Println("End uploading to ZincSearch")
	return nil
}
