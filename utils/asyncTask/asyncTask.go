package asyncTask

import (
	. "GoServer/middleWare/dataBases/redis"
	. "GoServer/utils/threadWorker"
)

type AsyncTaskType uint64

const (
	UNKNOWN_ASYNC_TASK AsyncTaskType = iota
	XASYNC_CREATE_THIRD_USER
)

type AsyncTaskFunc func(task *AsyncTaskEntity)

type AsyncTaskEntity struct {
	Type   AsyncTaskType
	Lock   *RedisLock
	Func   AsyncTaskFunc
	Param  map[string]interface{}
	Entity interface{}
}

func New() *AsyncTaskEntity {
	return &AsyncTaskEntity{
		Type:  UNKNOWN_ASYNC_TASK,
		Lock:  nil,
		Func:  nil,
		Param: nil,
		Entry: nil,
	}
}

func (task *AsyncTaskEntity) InsertWorkerQueue() {
	var work Job = task
	InsertAsyncTask(work)
}

func (task *AsyncTaskEntity) ExecTask() {
	if task == nil {
		return
	}

	if task.Func != nil {
		task.Func(task)
		return
	}
}
