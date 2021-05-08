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

func scheduleTarget(target Target, retry int, reads chan readOp, writes chan TargetResult, log Log) {
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
	waitTime := time.Since(start)
	log.Printf("Target %v started.. Waited for %v\n", target.Name, waitTime)
	err := runTarget(target, reads)
	elapsed := time.Since(start)
	if err != nil {
		if target.MaxRetries != nil && *target.MaxRetries > retry {
			log.Printf("Target %v failed, retrying\n", target.Name)
			scheduleTarget(target, retry+1, reads, writes, log)
		} else {
			log.Printf("Target %v failed after %v, reason: \n%v\n\n", target.Name, elapsed, err)
			writes <- TargetResult{&err, target, &waitTime, elapsed}
		}
	} else {
		log.Printf("Target %v finished successfully after %v\n", target.Name, elapsed)
		writes <- TargetResult{nil, target, &waitTime, elapsed}
	}

}

func runTarget(target Target, reads chan readOp) error {
	cmd := exec.Command("/bin/sh", "-c", target.Run)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if target.WorkDir != nil {
		if _, err := os.Stat(*target.WorkDir); os.IsNotExist(err) {
			return err
		}
		cmd.Dir = *target.WorkDir
	}
	cmd.Start()

	go func() {
		cancelled := false
		read := readOp{
			resp: make(chan []TargetResult),
		}

		// on failure, the for-loop below may try to msg reads, which has already been closed
		// this guards for this and closes the channel cleanly
		// this can be cleaned up by writing and retrying in runTarget instead
		defer func() {
			close(read.resp)
			if r := recover(); r != nil {
				cmd.Process.Kill()
			}
		}()
		for !cancelled {
			reads <- read
			resp := <-read.resp
			for _, t := range resp {
				if t.Target.Name == target.Name {
					cancelled = true
					break
				} else if t.Err != nil {
					cancelled = true
					cmd.Process.Kill()
					break
				}
			}
			time.Sleep(5 * time.Millisecond)
		}
	}()
	return cmd.Wait()
}

func RunPlan(targets []Target, log Log) ([]TargetResult, error) {
	reads := make(chan readOp)
	writes := make(chan TargetResult)

	defer close(reads)
	defer close(writes)
	var waitGroup sync.WaitGroup
	// Set number of effective goroutines we want to wait upon
	waitGroup.Add(len(targets))

	go func() {
		var state = []TargetResult{}
		// we wait for the last read to get results, so +1
		for len(state) < (len(targets) + 1) {
			select {
			case read := <-reads:
				read.resp <- state
			case write := <-writes:
				// waitGroup only until this condition
				if len(state) < len(targets) {
					waitGroup.Done()
				}
				state = append(state, write)
			}
		}
	}()

	for _, target := range targets {
		go scheduleTarget(target, 1, reads, writes, log)
	}
	waitGroup.Wait()

	read := readOp{
		resp: make(chan []TargetResult),
	}
	defer close(read.resp)
	reads <- read
	resp := <-read.resp
	var err error

	for _, t := range resp {
		if t.Err != nil {
			err = *t.Err
		}
	}

	return resp, err
}
