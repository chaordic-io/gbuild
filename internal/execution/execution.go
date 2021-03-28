package execution

import (
	"fmt"
	"os"
	"sync"

	"github.com/chaordic-io/gbuild/internal/config"
)

func RunTask(task config.Task, wg *sync.WaitGroup) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()
	fmt.Println("Done Running the task")

	panic("foo")
	wg.Done()
}

func RunPlan(tasks []config.Task, planName string) error {

	return nil
}
