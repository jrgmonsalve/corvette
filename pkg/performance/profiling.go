package performance

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

func StartCPUProfile(filename string) func() {
	cpuFile, err := os.Create(filename)
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	return func() {
		pprof.StopCPUProfile()
		cpuFile.Close()
	}
}

func WriteMemProfile(filename string) {
	memFile, err := os.Create(filename)
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer memFile.Close()
	runtime.GC() // recolect the garbage to get up-to-date statistics
	if err := pprof.WriteHeapProfile(memFile); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}
