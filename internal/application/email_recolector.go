package application

import (
	"fmt"
)

type EmailRecolector struct {
}

func NewEmailRecolector() *EmailRecolector {
	return &EmailRecolector{}
}

func (er *EmailRecolector) Collect() error {
	fmt.Println("Start Collecting emails")
	for i := 0; i < 1000; i++ {
		if i%100 == 0 {
			fmt.Println("Collecting email ", i)
		}
	}
	fmt.Println("End Collecting emails")
	return nil
}
