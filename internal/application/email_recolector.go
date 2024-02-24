package application

import "fmt"

type EmailRecolector struct {
}

func NewEmailRecolector() *EmailRecolector {
	return &EmailRecolector{}
}

func (er *EmailRecolector) Collect() error {
	fmt.Println("Collecting emails")
	return nil
}
