package replica

import (
	cmd "github.com/go-epaxos/epaxos/command"
)

type Prepare struct {
	ballot     *Ballot
	replicaId  int
	instanceId InstanceId
}

type PrepareReply struct {
	ok         bool
	ballot     *Ballot
	replicaId  int
	instanceId InstanceId
	status     int8
	cmds       []cmd.Command
	deps       []InstanceId
}

func (*Prepare) getType() uint8 {
	return prepareType
}

func (*PrepareReply) getType() uint8 {
	return prepareReplyType
}
