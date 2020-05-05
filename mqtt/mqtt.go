package mqtt

import (
	. "GoServer/utils"
	"bytes"
	"fmt"
	M "github.com/eclipse/paho.mqtt.golang"
	"runtime"
	"strconv"
	"time"
)

type Job interface {
	ExecTask() error
}

type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	Quit       chan bool
}

type Dispatcher struct {
	MaxWorkers int
	WorkerPool chan chan Job
	Quit       chan bool
}

type MqMsg struct {
	Topic   string
	Payload []byte
}

var (
	MaxWorker = runtime.NumCPU()
	MaxQueue  = 512
	JobQueue  chan Job
)

var client M.Client = nil

func NewWorker(workPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workPool,
		JobChannel: make(chan Job),
		Quit:       make(chan bool),
	}
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel
			select {
			case job := <-w.JobChannel:
				if err := job.ExecTask(); err != nil {
					fmt.Printf("excute job failed with err: %v", err)
				}
			case <-w.Quit:
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{MaxWorkers: maxWorkers, WorkerPool: pool, Quit: make(chan bool)}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}

	go d.Dispatch()
}

func (d *Dispatcher) Stop() {
	go func() {
		d.Quit <- true
	}()
}

func (d *Dispatcher) Dispatch() {
	for {
		select {
		case job := <-JobQueue:
			go func(job Job) {
				jobChannel := <-d.WorkerPool
				jobChannel <- job
			}(job)
		case <-d.Quit:
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

func (msg *MqMsg) ExecTask() error {
	ok, packet := MessageHandler(msg.Topic, msg.Payload)
	if ok && packet.JsonData != nil {
		PrintInfo("==========", msg.Topic, " time:", time.Now().Format(SystemConf().Timeformat), "=========", GetGoroutineID(), len(JobQueue))
		PrintInfo(packet.JsonData.(Protocol).Print())
	} else {
		PrintInfo("analysis failed ->Topic:%s Payload:%s\n", msg.Topic, msg.Payload)
	}
	return nil
}

func init() {
	runtime.GOMAXPROCS(MaxWorker)
	JobQueue = make(chan Job, MaxQueue)
	dispatcher := NewDispatcher(MaxWorker)
	dispatcher.Run()
}

var MessageCb M.MessageHandler = func(client M.Client, msg M.Message) {
	var work Job = &MqMsg{Topic: msg.Topic(), Payload: msg.Payload()}
	JobQueue <- work
}

func StartMqttService() error {
	opts := M.NewClientOptions().AddBroker(MqttConf().Host)
	opts.SetClientID(MqttConf().Token)
	opts.SetUsername(MqttConf().Name)
	opts.SetPassword(MqttConf().Pwsd)
	opts.SetDefaultPublishHandler(MessageCb)

	client = M.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	if token := client.Subscribe("/#", 0, nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func StopMqttService() {
	if client != nil && client.IsConnected()  {

		client.Disconnect(250)
		t1 := time.NewTicker(time.Millisecond * 250)

		for exit := false; exit != true;  {
			select {
			case <-t1.C:
				exit = true
				t1.Stop()
				break
			}
		}
	}
}