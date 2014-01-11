package replica

import (
	cmd "github.com/go-epaxos/epaxos/command"
)

type Accept struct {
	cmds  []cmd.Command
	seq   int
	deps  []InstanceIdType
	repId int
	insId InstanceIdType
	// Ballot
}

func (a *Accept) getType() uint8 {
	return acceptType
}
