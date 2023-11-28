package worker

import (
	"fmt"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	worker := NewWorker(5)

	for i := 1; i <= 20; i++ {
		task := NewTask(fmt.Sprint(i), func(taskId string, taskState *TaskState) {
			for taskId == "2" {
				time.Sleep(2 * time.Second)
				if *taskState == Stop {
					break
				}
			}

			fmt.Printf("task %s finish; workers: %d, tasks: %d\n", taskId, worker.GetWorkers(), worker.GetTasksLen())
		})
		worker.Push(task)
	}

	go func() {
		time.Sleep(10 * time.Second)
		worker.StopTask(fmt.Sprint(2))
	}()

	worker.Wait()
}
