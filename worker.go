package worker

import (
	"container/list"
	"fmt"
	"os"
	"sync"
)

type Worker struct {
	// 最大并发数
	maxWorker uint
	// 任务数量计数
	metric *Metric
	// 任务队列
	tasks *list.List
	// 队列锁
	mutex     *sync.Mutex
	waitGroup *sync.WaitGroup

	// 异常回调
	panicHandler func(taskId string, err interface{})
}

func NewWorker(maxWorker uint) *Worker {
	worker := &Worker{
		maxWorker: maxWorker,
		tasks:     list.New(),
		metric:    NewMetric(),
		mutex:     &sync.Mutex{},
		waitGroup: &sync.WaitGroup{},
	}

	return worker
}

// 获取最大线程数
func (w *Worker) GetMaxWorkers() uint {
	return w.maxWorker
}

// 获取当前正在运行的线程数
func (w *Worker) GetWorkers() uint {
	return uint(w.metric.BusyWorkers())
}

// 获取任务列表
func (w *Worker) GetTasks() *list.List {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.tasks
}

// 获取队列大小
func (w *Worker) GetTasksLen() int {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.tasks.Len()
}

// 等待执行结束
func (w *Worker) Wait() {
	w.waitGroup.Wait()
}

// 从队列中取出一个任务
func (w *Worker) pullTask() *list.Element {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// 遍历队列，获取准备中的人物
	for t := w.tasks.Front(); t != nil; t = t.Next() {
		task := t.Value.(*Task)
		if task.State == Ready {
			return t
		}
	}
	return nil
}

// 删除一个任务
func (w *Worker) removeTask(e *list.Element) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.tasks.Remove(e)
}

// 执行任务
func (w *Worker) runTask() {
	// 任务 -1
	defer w.metric.DecBusyWorker()
	defer w.waitGroup.Done()

	// 队列有数据的情况下，一直获取数据
	for w.GetTasksLen() > 0 {
		// 从队列中取出一个任务
		element := w.pullTask()

		// 如果未取到准备中的任务，则循环获取
		if element == nil {
			continue
		}

		// 强制类型转换
		task := element.Value.(*Task)

		defer func() {
			// 执行结束移除任务
			w.removeTask(element)
			// 设置完成状态
			task.State = Finish
		}()

		// 执行任务
		task.Start()

		// 恢复panic
		if err := recover(); err != nil {
			task.State = Error

			if w.panicHandler != nil {
				w.panicHandler(task.Id, err)
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

// 推送任务到队列
// 返回Push后的队列长度
func (w *Worker) Push(task *Task) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// 向队列插入一个任务
	w.tasks.PushBack(task)

	// 当前运行的数量小于总数量
	// 如果任务池满, 则不再创建新的协程
	if w.GetWorkers() < w.GetMaxWorkers() {
		// 运行中任务 +1
		w.metric.IncBusyWorker()
		w.waitGroup.Add(1)

		// 启动新的worker
		go w.runTask()
	}
}

/*
停止指定任务。
  - 注意！该函数只会改变任务状态为stop，并不会主动停止任务，所以请在回调函数中自行处理停止逻辑
  - Attention! This function will only change the task state to stop and will not actively stop the task, so please handle the stop logic yourself in the callback function
*/
func (w *Worker) StopTask(taskId string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// 遍历队列，设置任务状态
	for t := w.tasks.Front(); t != nil; t = t.Next() {
		task := t.Value.(*Task)
		if task.Id == taskId {
			task.State = Stop
			break
		}
	}
}

// 停止所有任务
func (w *Worker) StopAllTask() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// 遍历队列，设置任务状态
	for t := w.tasks.Front(); t != nil; t = t.Next() {
		task := t.Value.(*Task)
		task.State = Stop
	}
}
