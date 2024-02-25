package helpers

import (
	"errors"
	"flag"
)

func GetCommandLineArgument(error_message string) (string, error) {
	flag.Parse()
	argument := flag.Arg(0)
	if argument == "" {
		return "", errors.New(error_message)
	}
	return argument, nil
}
