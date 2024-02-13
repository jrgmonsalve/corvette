package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"
	"syscall"
)

func main() {
	// Crear archivo de perfil de CPU
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Iniciar perfilado de CPU
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("Error starting CPU profile: ", err)
	}

	// Preparar manejo de señal de interrupción (Ctrl+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Usar una goroutine para manejar la señal
	go func() {
		<-c
		fmt.Println("\nCtrl+C pressed. Stopping CPU profile and exiting...")
		pprof.StopCPUProfile() // Detener el perfilado de CPU
		os.Exit(0)
	}()

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

	pprof.StopCPUProfile()
}
