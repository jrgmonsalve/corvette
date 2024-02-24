package main

import (
	"github.com/jrgmonsalve/corvette/cmd/cli/internal/application"
	"github.com/jrgmonsalve/corvette/cmd/cli/pkg/performance"
)

func main() {

	performance.StartRecord()

	recolector := application.NewEmailRecolector()
	indexer := application.NewIndexer(recolector)
	indexer.Start()

	performance.StopRecord()

}
