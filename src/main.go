package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
)

func main() {

	////Prceso de rendimiento de la aplicación/////////////
	cpu, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpu)
	defer pprof.StopCPUProfile()
	////////Fin prceso de rendimiento de la aplicación/////////

	LoadEnvFile()

	rootPath := GetRootPath()

	nodes := GetNodes(rootPath)
	totalNodes := len(nodes)
	increment := 0
	for _, node := range nodes {
		increment++
		prefixLog := fmt.Sprintf("Processing %d / %d directory: %s", increment, totalNodes, node.Name())
		if node.IsDir() {
			fullPath := filepath.Join(rootPath, node.Name())
			ProcessDirectory(fullPath+"/all_documents", prefixLog)
		}
	}

}
