# go worker
A `goWorker` is a simple task queue implemented by List


# Usage

```golang
// new worker pool
worker := NewWorker(5)

// new task
task := NewTask("[custom task id]", func(taskId string, taskState *TaskState) {
    // ........
})

// worker.StopTask("[task id]")
// Attention! This function will only change the task state to stop and will not actively stop the task, so please handle the stop logic yourself in the callback function

// push to worker pool
worker.Push(task)

// wait
worker.Wait()
```