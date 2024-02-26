package helpers

import (
	"errors"
	"flag"
)

func GetArgumentFromCLI(error_message string) (string, error) {
	flag.Parse()
	argument := flag.Arg(0)
	if argument == "" {
		return "", errors.New(error_message)
	}
	return argument, nil
}
