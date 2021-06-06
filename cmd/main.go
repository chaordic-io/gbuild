package main

import (
	"flag"
	"os"
	"time"

	"github.com/chaordic-io/gbuild/internal"
)

var target string
var fileName string
var version bool

func init() {
	flag.StringVar(&target, "t", "build", "Define target execution plan")
	flag.StringVar(&fileName, "f", ".gbuild.yaml", "File to run")
	flag.BoolVar(&version, "v", false, "Print the installed gbuild version")
}

func main() {
	start := time.Now()
	log := internal.OSLog{}
	flag.Parse()
	if version {
		internal.PrintVersionInfo()
		os.Exit(0)
	}
	log.Printf("Running target execution plan '%v' on file %v..\n\n", target, fileName)
	conf, err := internal.LoadConfig(fileName, log)
	if err != nil {
		log.Printf("Could not read config file %v, reason: %v exiting\n\n", fileName, err.Error())
		os.Exit(1)
	}
	targets, err := internal.GetTargetsForPlan(conf, target, log)
	if err != nil {
		log.Printf("Could not get targets for %v, reason: %v exiting\n\n", target, err.Error())
		os.Exit(1)
	}

	err = internal.LoadCache(nil, &targets, nil)
	if err != nil {
		log.Printf("Failed to get cache, reason: %v\n\n", err.Error())
		os.Exit(1)
	}

	_, err = internal.RunPlan(targets, log)
	if err != nil {
		log.Printf("Error executing plan, reason: %v\n\n", err.Error())
		os.Exit(1)
	}

	err = internal.PutCache(nil, &targets, nil)
	if err != nil {
		log.Printf("Failed to put cache, reason: %v\n\n", err.Error())
		os.Exit(1)
	}
	elapsed := time.Since(start)
	log.Printf("Build completed successfully after %v\n\n", elapsed)
}
