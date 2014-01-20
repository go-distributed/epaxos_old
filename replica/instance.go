package replica

import (
	cmd "github.com/go-epaxos/epaxos/command"
)

const (
	preAccepted int8 = iota
	accepted
	committed
	executed
)

// a bookkeeping for infos like maxBallot, # of nack, # of ok, etc
type InstanceInfo struct {
	preaccCnt     int
	haveDiffReply bool

	acceptNackCnt int
	acceptOkCnt   int

	prepareCnt int
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

func NewInstance(cmds []cmd.Command, deps dependencies, status int8) *Instance {
	return &Instance{
		cmds:   cmds,
		deps:   deps,
		status: status,
		info:   NewInstanceInfo(),
	}
}

func (i *Instance) allReplyTheSame() bool {
	return i.info.haveDiffReply == false
}

func (i *Instance) afterStatus(status int8) bool {
	return i.status > status
}
