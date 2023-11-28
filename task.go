package worker

import (
	"strings"

	"github.com/google/uuid"
)

type TaskState int8

const (
	// 准备
	Ready TaskState = iota
	// 运行
	Run
	// 停止
	Stop
	// 完成
	Finish
	// 异常
	Error
)

type Task struct {
	Id      string
	State   TaskState
	Handler func(taskId string, taskState *TaskState)
}

// 新建任务实例
func NewTask(Id string, Fn func(taskId string, taskState *TaskState)) *Task {
	if Id == "" {
		Id = strings.ReplaceAll(uuid.NewString(), "-", "")
	}
	return &Task{
		Id:      Id,
		State:   Ready,
		Handler: Fn,
	}
}

// 启动任务
func (t *Task) Start() {
	t.State = Run
	t.Handler(t.Id, &t.State)
}
