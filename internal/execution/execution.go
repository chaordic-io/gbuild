package execution

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/chaordic-io/gbuild/internal/config"
)

type readOp struct {
	resp chan []TargetResult
}

type TargetResult struct {
	Err     *error
	Target  config.Target
	Wait    *time.Duration
	Elapsed time.Duration
}

func scheduleTarget(target config.Target, retry int, wg *sync.WaitGroup, reads chan readOp, writes chan TargetResult) {
	start := time.Now()

	if target.DependsOn != nil && len(*target.DependsOn) > 0 {
		completed := false

		for !completed {
			read := readOp{
				resp: make(chan []TargetResult),
			}
			reads <- read
			resp := <-read.resp
			matches := 0
			for _, t := range resp {
				for _, d := range *target.DependsOn {
					if d == t.Target.Name {
						matches++
					}
				}
			}
			completed = matches == len(*target.DependsOn)
		}
		time.Sleep(5 * time.Millisecond)
	}
	waitTime := time.Since(start)
	fmt.Printf("Target %v started.. Waited for %v\n", target.Name, waitTime)
	err := RunTarget(target)
	if err != nil {
		if target.MaxRetries != nil && *target.MaxRetries > retry {
			fmt.Printf("Target %v failed, retrying\n", target.Name)
			scheduleTarget(target, retry+1, wg, reads, writes)
		} else {
			elapsed := time.Since(start)
			fmt.Printf("Target %v failed after %v\n", target.Name, elapsed)
			writes <- TargetResult{&err, target, &waitTime, elapsed}
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Target %v finished successfully after %v\n", target.Name, elapsed)
	writes <- TargetResult{nil, target, &waitTime, elapsed}

	wg.Done()
}

func RunTarget(target config.Target) error {
	cmd := exec.Command("/bin/sh", "-c", target.Run)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if target.WorkDir != nil {
		cmd.Dir = *target.WorkDir
	}
	err := cmd.Run()
	return err
}

func RunPlan(targets []config.Target) ([]TargetResult, error) {
	reads := make(chan readOp)
	writes := make(chan TargetResult)
	var waitGroup sync.WaitGroup
	// Set number of effective goroutines we want to wait upon
	waitGroup.Add(len(targets))

	go func() {
		var state = []TargetResult{}
		for {
			select {
			case read := <-reads:
				read.resp <- state
			case write := <-writes:
				state = append(state, write)

				if write.Err != nil {
					e := *write.Err
					fmt.Println(e.Error())
					os.Exit(1)
				}
			}
		}
	}()

	for _, target := range targets {
		go scheduleTarget(target, 0, &waitGroup, reads, writes)
	}
	waitGroup.Wait()
	read := readOp{
		resp: make(chan []TargetResult),
	}
	reads <- read
	resp := <-read.resp

	return resp, nil
}
