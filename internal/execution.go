package internal

import (
	"os"
	"os/exec"
	"sync"
	"time"
)

type readOp struct {
	resp chan []TargetResult
}

type TargetResult struct {
	Err     *error
	Target  Target
	Wait    *time.Duration
	Elapsed time.Duration
}

func scheduleTarget(target Target, waitGroup *sync.WaitGroup, retry int, reads chan readOp, writes chan TargetResult, log Log) {
	start := time.Now()

	if target.DependsOn != nil && len(*target.DependsOn) > 0 {
		completed := false
		read := readOp{
			resp: make(chan []TargetResult),
		}
		defer close(read.resp)
		for !completed {
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
	runTarget(target, waitGroup, retry, reads, writes, log, start)
}

func runTarget(target Target, waitGroup *sync.WaitGroup, retry int, reads chan readOp, writes chan TargetResult, log Log, start time.Time) {
	waitTime := time.Since(start)
	log.Printf("Target %v started.. Waited for %v\n", target.Name, waitTime)
	cmd := exec.Command("/bin/sh", "-c", target.Run)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if target.WorkDir != nil {
		if _, err := os.Stat(*target.WorkDir); os.IsNotExist(err) {
			writes <- TargetResult{&err, target, &waitTime, waitTime}
			waitGroup.Done()
			waitGroup.Done()
			return
		}
		cmd.Dir = *target.WorkDir
	}
	cmd.Start()

	go func() {
		cancelled := false
		read := readOp{
			resp: make(chan []TargetResult),
		}
		defer close(read.resp)
		for !cancelled {
			reads <- read
			resp := <-read.resp
			for _, t := range resp {
				if t.Target.Name == target.Name {
					cancelled = true
					waitGroup.Done()
					break
				} else if t.Err != nil {
					cancelled = true
					waitGroup.Done()
					cmd.Process.Kill()
					break
				}
			}
			if !cancelled {
				time.Sleep(5 * time.Millisecond)
			}
		}
	}()
	err := cmd.Wait()

	elapsed := time.Since(start)
	if err != nil {
		if target.MaxRetries != nil && *target.MaxRetries > retry {
			log.Printf("Target %v failed, retrying\n", target.Name)
			waitGroup.Add(1) // add to waitgroup on retry
			runTarget(target, waitGroup, retry+1, reads, writes, log, start)
		} else {
			log.Printf("Target %v failed after %v, reason: \n%v\n\n", target.Name, elapsed, err)
			writes <- TargetResult{&err, target, &waitTime, elapsed}
			waitGroup.Done()
		}
	} else {
		log.Printf("Target %v finished successfully after %v\n", target.Name, elapsed)
		writes <- TargetResult{nil, target, &waitTime, elapsed}
		waitGroup.Done()
	}

}

func RunPlan(targets []Target, log Log) ([]TargetResult, error) {
	reads := make(chan readOp)
	writes := make(chan TargetResult)

	defer close(reads)
	defer close(writes)
	var waitGroup sync.WaitGroup
	// Set number of effective goroutines we want to wait upon
	// this is x2, because we have a go routine watching if a target should be cancelled
	waitGroup.Add(len(targets) * 2)

	go func() {
		var state = []TargetResult{}
		// we wait for the last read to get results, so +1
		for len(state) < (len(targets) + 1) {
			select {
			case read := <-reads:
				read.resp <- state
			case write := <-writes:
				state = append(state, write)
			}
		}
	}()

	for _, target := range targets {
		go scheduleTarget(target, &waitGroup, 1, reads, writes, log)
	}
	waitGroup.Wait()

	read := readOp{
		resp: make(chan []TargetResult),
	}
	defer close(read.resp)
	reads <- read
	resp := <-read.resp
	var err error
	// TODO close channels cleanly

	for _, t := range resp {
		if t.Err != nil {
			err = *t.Err
		}
	}

	return resp, err
}
