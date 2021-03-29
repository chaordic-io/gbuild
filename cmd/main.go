package main

import (
	"flag"
	"fmt"
	"os"
	"time"

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
	flag.Parse()
	fmt.Printf("Running target execution plan '%v' on file %v..\n\n", target, fileName)
	conf, err := config.LoadConfig(fileName)
	if err != nil {
		fmt.Printf("Could not read config file %v, reason: %v exiting\n\n", fileName, err.Error())
		os.Exit(1)
	}
	targets, err := config.GetTargetsForPlan(conf, target)
	if err != nil {
		fmt.Printf("Could not get targets for %v, reason: %v exiting\n\n", target, err.Error())
		os.Exit(1)
	}
	_, err = execution.RunPlan(targets)
	if err != nil {
		fmt.Printf("Could execute plan, reason: %v exiting\n\n", err.Error())
		os.Exit(1)
	}
	elapsed := time.Since(start)
	fmt.Printf("Build completed successfully after %v\n\n", elapsed)

}
