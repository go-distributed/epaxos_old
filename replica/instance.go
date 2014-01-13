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
	cmds   []cmd.Command
	deps   []InstanceIdType
	status int8
	ballot uint64
	info   *InstanceInfo
}
