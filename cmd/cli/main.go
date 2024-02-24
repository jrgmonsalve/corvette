package main

import (
	"github.com/jrgmonsalve/corvette/cmd/cli/internal/application"
	"github.com/jrgmonsalve/corvette/cmd/cli/pkg/performance"
)

func main() {

	stopCPUProfile := performance.StartCPUProfile("cpu_profile.prof")
	defer stopCPUProfile()

	recolector := application.NewEmailRecolector()
	indexer := application.NewIndexer(recolector)
	indexer.Start()

	performance.WriteMemProfile("mem_profile.prof")

}
