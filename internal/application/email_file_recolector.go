package application

import (
	"github.com/jrgmonsalve/corvette/cmd/cli/internal/domain"
	"github.com/jrgmonsalve/corvette/cmd/cli/pkg/helpers"
)

type EmailFileRecolector struct {
}

func NewEmailFileRecolector() *EmailFileRecolector {
	return &EmailFileRecolector{}
}

func (efr *EmailFileRecolector) Collect() error {
	path, err := helpers.GetCommandLineArgument("You must provide the root path")
	if err != nil {
		return err
	}

	var emails []domain.Email
	err = helpers.RecursiveFileReader(path, 0, 3, &emails, 4)
	if err != nil {
		return err
	}
	return nil
}
