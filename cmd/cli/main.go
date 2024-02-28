package main

import (
	"github.com/jrgmonsalve/corvette/cmd/cli/internal/application"
	api_zincsearch "github.com/jrgmonsalve/corvette/cmd/cli/internal/infrastructure/services/zincsearch"
	"github.com/jrgmonsalve/corvette/cmd/cli/pkg/helpers"
	"github.com/jrgmonsalve/corvette/cmd/cli/pkg/performance"
)

func main() {

	stopCPUProfile := performance.StartCPUProfile("cpu_profile.prof")
	defer stopCPUProfile()
	helpers.LoadEnvFile()
	emailRepository := api_zincsearch.NewZincSearchService()
	collector := application.NewEmailFileCollector()
	indexer := application.NewIndexer(collector, emailRepository)

	indexer.Start()

	performance.WriteMemProfile("mem_profile.prof")

}
