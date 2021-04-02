package mr

import (
	"log"
	"strconv"
	"sync"
)
import "net"
import "os"
import "net/rpc"
import "net/http"

type Master struct {
	// Your definitions here.
	FilesName         map[string]int //key: input file name, val: file status
	ReduceTasksStatus map[int]int    //key: reduce task number, val: reduce task status
	CurrMap           int
	NReduce           int
	MapFinished       bool
	ReduceFinished    bool
	InterFileName     [][]string
	Mutex             *sync.RWMutex
}

// task status
const (
	UnAssigned = 0
	Assigned   = 1
	Finished   = 2
)

var mapTasks chan string
var reduceTasks chan int

// Your code here -- RPC handlers for the worker to call.
func (m *Master) HandlerRPC(args *MsgArgs, reply *MsgReply) error {
	jobType := args.JobTypeArgs
	switch jobType {
	case MsgAskTaskRPC:
		select {
		case fileName := <-mapTasks:
			reply.FileName = fileName
			reply.CurrMap = m.CurrMap
			reply.NReduce = m.NReduce
			reply.JobTypeReply = "map"
			m.FilesName[fileName] = Assigned
			m.CurrMap++
			go m.workerTime()
			return nil
		case reduceTask := <-reduceTasks:
			reply.JobTypeReply = "reduce"
			reply.NReduce = m.NReduce
			reply.CurrReduce = reduceTask
			reply.InterFileName = m.InterFileName[reduceTask]
			m.ReduceTasksStatus[reduceTask] = Assigned
			go m.workerTime()
			return nil
		}
	case MsgFinishMapRPC:
		m.FilesName[args.Msg] = Finished
	case MsgFinishReduceRPC:
		reduceIntex, _ := strconv.Atoi(args.Msg)
		m.ReduceTasksStatus[reduceIntex] = Finished
	}
	return nil
}

func (m *Master) workerTime() {

}

func(m *Master) HandlerInterFileRPC(args *MsgArgs, reply *MsgReply) {
	reduceNum := args.TaskID
	fileName := args.Msg
	m.InterFileName[reduceNum] = append(m.InterFileName[reduceNum], fileName)
}

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//
func (m *Master) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

//
// start a thread that listens for RPCs from worker.go
//
func (m *Master) server() {
	mapTasks = make(chan string)
	reduceTasks = make(chan int)

	rpc.Register(m)
	rpc.HandleHTTP()

	go m.createTasks()

	//l, e := net.Listen("tcp", ":1234")
	sockname := masterSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrmaster.go calls Done() periodically to find out
// if the entire job has finished.
//
func (m *Master) Done() bool {
	//ret := false

	// Your code here.
	ret := m.ReduceFinished

	return ret
}

//
// create a Master.
// main/mrmaster.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeMaster(files []string, nReduce int) *Master {
	m := Master{}

	// Your code here.
	m.FilesName = make(map[string]int)
	m.ReduceTasksStatus = make(map[int]int)
	m.CurrMap = 0
	m.NReduce = nReduce
	m.MapFinished = false
	m.ReduceFinished = false
	m.InterFileName = make([][]string, m.NReduce)
	m.Mutex = new(sync.RWMutex)

	for _, file := range files {
		m.FilesName[file] = UnAssigned
	}

	for i := 0; i < nReduce; i++ {
		m.ReduceTasksStatus[i] = UnAssigned
	}

	m.server()
	return &m
}

func (m *Master) createTasks() {
	for fileName, status := range m.FilesName {
		if status == UnAssigned {
			mapTasks <- fileName
		}
	}
	done := false
	for !done {
		done = m.checkMapFinished()
	}
	m.MapFinished = true

	for taskID, status := range m.ReduceTasksStatus {
		if status == UnAssigned {
			reduceTasks <- taskID
		}
	}
	done = false
	for !done {
		done = m.checkMReduceFinished()
	}
	m.ReduceFinished = true
}

func (m *Master) checkMapFinished() bool {
	m.Mutex.RLock()
	defer m.Mutex.RUnlock()
	for _, v := range m.FilesName {
		if v != Finished {
			return false
		}
	}
	return true
}

func (m *Master) checkMReduceFinished() bool {
	m.Mutex.RLock()
	defer m.Mutex.RUnlock()
	for _, v := range m.ReduceTasksStatus {
		if v != Finished {
			return false
		}
	}
	return true
}
