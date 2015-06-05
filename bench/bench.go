package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/lazywei/mockingbird"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()

	fileContent, err := ioutil.ReadFile("./samples/ABAP/24-game-solve.abap")
	if err != nil {
		panic(err)
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	defer timeTrack(time.Now(), "ExtractTokens")
	mockingbird.ExtractTokens(string(fileContent))
}
