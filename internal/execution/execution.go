package execution

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/chaordic-io/gbuild/internal/config"
)

type readOp struct {
	resp chan []TargetResult
}

type TargetResult struct {
	Err     *interface{}
	Target  config.Target
	Wait    *time.Duration
	Elapsed time.Duration
}

func scheduleTarget(target config.Target, wg *sync.WaitGroup, reads chan readOp, writes chan TargetResult) {
	start := time.Now()
	var waitTime time.Duration

	defer func() {
		if err := recover(); err != nil {
			elapsed := time.Now().Sub(start)
			writes <- TargetResult{&err, target, &waitTime, elapsed}
		}
	}()
	if target.DependsOn != nil && len(*target.DependsOn) > 0 {
		completed := false

		for !completed {
			read := readOp{
				resp: make(chan []TargetResult),
			}
			reads <- read
			resp := <-read.resp
			for _, t := range resp {
				for _, d := range *target.DependsOn {
					if d == t.Target.Name {
						completed = true
					}
				}
			}
		}
		time.Sleep(5 * time.Millisecond)
	}
	waitTime = time.Now().Sub(start)
	//fmt.Println(target.Run)
	elapsed := time.Now().Sub(start)
	writes <- TargetResult{nil, target, &waitTime, elapsed}

	wg.Done()
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
					// TODO report more here
					fmt.Println(write.Err)
					os.Exit(1)
				}
			}
		}
	}()

	for _, target := range targets {
		go scheduleTarget(target, &waitGroup, reads, writes)
	}
	waitGroup.Wait()
	read := readOp{
		resp: make(chan []TargetResult),
	}
	reads <- read
	resp := <-read.resp

	return resp, nil
}
