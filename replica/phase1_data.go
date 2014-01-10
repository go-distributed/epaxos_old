package replica

import (
	cmd "github.com/go-epaxos/epaxos/command"
)

const (
	ProposeType = iota
	PreAcceptType
)

type Propose struct {
	Cmds []cmd.Command
}

type PreAccept struct {
	Cmds  []cmd.Command
	Seq   int
	Deps  []InstanceIdType
	repId int
	InsId InstanceIdType
}

func (*PreAccept) getType() uint8 {
	return PreAcceptType
}
