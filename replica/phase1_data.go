package replica

import (
	cmd "github.com/go-epaxos/epaxos/command"
)

const (
	proposeType = iota
	preAcceptType
	preAcceptOKType
	preAcceptReplyType
)

type Propose struct {
	cmds []cmd.Command
}

type PreAccept struct {
	cmds  []cmd.Command
	deps  []InstanceIdType
	repId int
	insId InstanceIdType
}

type PreAcceptOK struct {
	insId InstanceIdType
}

type PreAcceptReply struct {
	deps  []InstanceIdType
	repId int
	insId InstanceIdType
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
