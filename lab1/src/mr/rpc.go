package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import "os"
import "strconv"

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

type JobTypeArgs int

// Add your RPC definitions here.
const (
	MsgAskTaskRPC      JobTypeArgs = 0
	MsgFinishMapRPC    JobTypeArgs = 2
	MsgFinishReduceRPC JobTypeArgs = 3
)

type MsgArgs struct {
	JobTypeArgs JobTypeArgs
	Msg         string
	TaskID      int
}

type MsgReply struct {
	FileName     string
	NMap         int
	NReduce      int
	JobTypeReply string
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the master.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func masterSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
