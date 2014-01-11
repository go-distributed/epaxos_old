package replica

import (
	cmd "github.com/go-epaxos/epaxos/command"
)

const (
	proposeType = iota
	preAcceptType
	acceptType
)

type Propose struct {
	cmds []cmd.Command
}

type PreAccept struct {
	cmds  []cmd.Command
	seq   int
	deps  []InstanceIdType
	repId int
	insId InstanceIdType
}

func (*PreAccept) getType() uint8 {
	return preAcceptType
}
