package replica

import (
	cmd "github.com/go-epaxos/epaxos/command"
)

const (
	none int8 = iota
	preaccepted
	accepted
	committed
	executed
)

type Instance struct {
	Cmds   []cmd.Command
	Seq    int
	Deps   []InstanceIdType
	Status int8
}

