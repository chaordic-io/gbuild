package main

import (
	"flag"
	"os"
	"time"

	"github.com/chaordic-io/gbuild/internal/common"
	"github.com/chaordic-io/gbuild/internal/config"
	"github.com/chaordic-io/gbuild/internal/execution"
)

var target string
var fileName string

func init() {
	flag.StringVar(&target, "t", "build", "Define target execution plan")
	flag.StringVar(&fileName, "f", ".gbuild.yaml", "File to run")
}

func main() {
	start := time.Now()
	log := common.OSLog{}
	flag.Parse()
	log.Printf("Running target execution plan '%v' on file %v..\n\n", target, fileName)
	conf, err := config.LoadConfig(fileName, log)
	if err != nil {
		log.Printf("Could not read config file %v, reason: %v exiting\n\n", fileName, err.Error())
		os.Exit(1)
	}
	targets, err := config.GetTargetsForPlan(conf, target, log)
	if err != nil {
		log.Printf("Could not get targets for %v, reason: %v exiting\n\n", target, err.Error())
		os.Exit(1)
	}
	_, err = execution.RunPlan(targets, log)
	if err != nil {
		log.Printf("Could execute plan, reason: %v exiting\n\n", err.Error())
		os.Exit(1)
	}
	elapsed := time.Since(start)
	log.Printf("Build completed successfully after %v\n\n", elapsed)

}
