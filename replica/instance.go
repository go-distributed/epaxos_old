package replica

import (
	cmd "github.com/go-epaxos/epaxos/command"
)

const (
	preaccepted int8 = iota
	accepted
	committed
	executed
)

// a bookkeeping for infos like maxBallot, # of nack, # of ok, etc
type InstanceInfo struct {
	preAcceptCount int
	haveDiffReply  bool

	acceptNackCount int
	acceptCount     int

	prepareCount int

	recovery *RecoveryInfo
}

type RecoveryInfo struct {
	preAcceptCount int
}

type Instance struct {
	cmds   []cmd.Command
	deps   dependencies
	status int8
	ballot *Ballot
	info   *InstanceInfo
}

func NewInstanceInfo() *InstanceInfo {
	return &InstanceInfo{}
}

func (Inst *Instance) allReplyTheSame() bool {
	return Inst.info.haveDiffReply == false
}
