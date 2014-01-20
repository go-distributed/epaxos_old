package replica

import (
	cmd "github.com/go-epaxos/epaxos/command"
)

const (
	proposeType = iota
	preAcceptType
	preAcceptOKType
	preAcceptReplyType
	acceptType
	acceptReplyType
	commitType
	prepareType
	prepareReplyType
)

type Propose struct {
	cmds []cmd.Command
}

type PreAccept struct {
	cmds       []cmd.Command
	deps       dependencies
	replicaId  int
	instanceId InstanceId
	ballot     *Ballot
}

type PreAcceptOK struct {
	instanceId InstanceId
}

type PreAcceptReply struct {
	deps       []InstanceId
	replicaId  int
	instanceId InstanceId
	ok         bool
}

func (*PreAccept) getType() uint8 {
	return preAcceptType
}
func (*PreAcceptOK) getType() uint8 {
	return preAcceptOKType
}
func (*PreAcceptReply) getType() uint8 {
	return preAcceptReplyType
}
