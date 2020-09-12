package utils

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
)

type Job interface {
	ExecTask() error
}

type worker struct {
	workerPool chan chan Job
	jobChannel chan Job
	quit       chan bool
}

type dispatcher struct {
	maxWorkers int
	workerPool chan chan Job
	quit       chan bool
}

var (
	maxWorker = runtime.NumCPU()
	maxQueue  = 512
	jobQueue  chan Job
)

func newWorker(workPool chan chan Job) worker {
	return worker{
		workerPool: workPool,
		jobChannel: make(chan Job),
		quit:       make(chan bool),
	}
}

func (w worker) start() {
	go func() {
		for {
			w.workerPool <- w.jobChannel
			select {
			case job := <-w.jobChannel:
				if err := job.ExecTask(); err != nil {
					fmt.Printf("excute job failed with err: %v", err)
				}
			case <-w.quit:
				return
			}
		}
	}()
}

func (w worker) stop() {
	go func() {
		w.quit <- true
	}()
}

func newDispatcher(maxWorkers int) *dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &dispatcher{maxWorkers: maxWorkers, workerPool: pool, quit: make(chan bool)}
}

func (d *dispatcher) Run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := newWorker(d.workerPool)
		worker.start()
	}

	go d.dispatch()
}

func (d *dispatcher) stop() {
	go func() {
		d.quit <- true
	}()
}

func (d *dispatcher) dispatch() {
	for {
		select {
		case job := <-jobQueue:
			go func(job Job) {
				jobChannel := <-d.workerPool
				jobChannel <- job
			}(job)
		case <-d.quit:
			return
		}
	}
}

func GetGoroutineID() uint64 {
	b := make([]byte, 64)
	runtime.Stack(b, false)
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func InsertAsynTask(jobTask Job) {
	jobQueue <- jobTask
}

func GetWorkerQueueSize() int {
	return len(jobQueue)
}

func init() {
	runtime.GOMAXPROCS(maxWorker)
	jobQueue = make(chan Job, maxQueue)
	dispatcher := newDispatcher(maxWorker)
	dispatcher.Run()
}
