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

func runTarget(task config.Target, wg *sync.WaitGroup, reads chan readOp, writes chan TargetResult) {
	start := time.Now()
	var waitTime time.Duration

	defer func() {
		if err := recover(); err != nil {
			elapsed := time.Now().Sub(start)
			writes <- TargetResult{&err, task, &waitTime, elapsed}
		}
	}()
	if task.DependsOn != nil && len(*task.DependsOn) > 0 {
		completed := false

		for !completed {
			read := readOp{
				resp: make(chan []TargetResult),
			}
			reads <- read
			resp := <-read.resp
			for _, t := range resp {
				for _, d := range *task.DependsOn {
					if d == t.Target.Name {
						completed = true
					}
				}
			}
		}
		time.Sleep(5 * time.Millisecond)
	}
	waitTime = time.Now().Sub(start)
	fmt.Println(task.Run)
	elapsed := time.Now().Sub(start)
	writes <- TargetResult{nil, task, &waitTime, elapsed}

	wg.Done()
}

func RunPlan(tasks []config.Target) ([]TargetResult, error) {
	reads := make(chan readOp)
	writes := make(chan TargetResult)
	var waitGroup sync.WaitGroup
	// Set number of effective goroutines we want to wait upon
	waitGroup.Add(len(tasks))

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

	for _, task := range tasks {
		go runTarget(task, &waitGroup, reads, writes)
	}
	waitGroup.Wait()
	read := readOp{
		resp: make(chan []TargetResult),
	}
	reads <- read
	resp := <-read.resp

	return resp, nil
}
