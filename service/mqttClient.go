package Service

import (
	. "GoServer/packet"
	. "GoServer/utils"
	"bytes"
	"fmt"
	M "github.com/eclipse/paho.mqtt.golang"
	//	"reflect"
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
	Broker  string
	Topic   string
	Payload []byte
}

var (
	MaxWorker = runtime.NumCPU()
	MaxQueue  = 512
	JobQueue  chan Job
)

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
	ok, packet := MessageHandler(msg.Payload)
	if ok && packet.JsonData != nil {
		//fmt.Println("==========", msg.Topic, "time:", time.Now().Format(conf.GetConfig().GetSystem().Timeformat), "=========", GetGoroutineID(), len(JobQueue))
		//fmt.Println(packet.JsonData.(Protocol).Print())
		PrintInfo("[", msg.Broker, "] =========>>", msg.Topic, " time:", time.Now().Format(GetConfig().GetSystem().Timeformat), "=========", GetGoroutineID(), len(JobQueue))
		PrintInfo(packet.JsonData.(Protocol).Print())
	} else {
		fmt.Printf("analysis failed ->Topic:%s Payload:%s\n", msg.Topic, msg.Payload)
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
	rOps := client.OptionsReader()
	servers := rOps.Servers()
	broker := servers[0]

	//fmt.Println("M.MessageHandler--->", broker, reflect.TypeOf(broker).String())
	var work Job = &MqMsg{Broker: broker.Host, Topic: msg.Topic(), Payload: msg.Payload()}
	JobQueue <- work
}

func StartMqttService() error {
	opts1 := M.NewClientOptions().AddBroker(GetConfig().GetMqtt().Host)
	opts1.SetClientID(GetConfig().GetMqtt().Token)
	opts1.SetUsername(GetConfig().GetMqtt().Name)
	opts1.SetPassword(GetConfig().GetMqtt().Pwsd)
	opts1.SetAutoReconnect(true)
	opts1.SetDefaultPublishHandler(MessageCb)

	opts2 := M.NewClientOptions().AddBroker("tcp://47.106.235.93:1883")
	opts2.SetClientID(GetConfig().GetMqtt().Token)
	opts2.SetUsername(GetConfig().GetMqtt().Name)
	opts2.SetPassword(GetConfig().GetMqtt().Pwsd)
	opts2.SetAutoReconnect(true)
	opts2.SetDefaultPublishHandler(MessageCb)

	Client1 := M.NewClient(opts1)
	if token := Client1.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	if token := Client1.Subscribe("/#", 0, nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	Client2 := M.NewClient(opts2)
	if token := Client2.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	if token := Client2.Subscribe("/#", 0, nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}
