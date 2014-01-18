package replica

import (
	cmd "github.com/go-epaxos/epaxos/command"
)

type Prepare struct {
	ballot uint64
	repId  int
	insId  InstanceIdType
}

type PrepareReply struct {
	ok     bool
	ballot uint64
	status int8
	cmds   []cmd.Command
	deps   []InstanceIdType
	repId  int
	insId  InstanceIdType
}

func (*Prepare) getType() uint8 {
	return prepareType
}

func (*PrepareReply) getType() uint8 {
	return prepareReplyType
}
