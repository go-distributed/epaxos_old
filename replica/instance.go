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
	preAcceptCount  int
	needAcceptPhase bool

	acceptNackCount int
	acceptCount     int

	prepareCount int

	recovery *RecoveryInfo
}

type RecoveryInfo struct {
	hasCommandCount int
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

func (i *Instance) processPreAcceptReply(par *PreAcceptReply, quorum, fastQuorum int) int8 {
	i.info.preAcceptCount++

	if update := i.unionDeps(par.deps); !update {
		if i.preAcceptCount() > 1 {
			i.setNeedAcceptPhase(true)
		}
	}

	if i.preAcceptCount() >= quorum-1 && i.needAcceptPhase() {
		i.status = accepted
	}

	if i.preAcceptCount() == fastQuorum && !i.needAcceptPhase() {
		i.status = committed
	}

	return i.status
}

func (i *Instance) preAcceptCount() int {
	return i.info.preAcceptCount
}

func (i *Instance) needAcceptPhase() bool {
	return i.info.needAcceptPhase
}

func (i *Instance) setNeedAcceptPhase(need bool) {
	i.info.needAcceptPhase = need
}

func (i *Instance) unionDeps(deps dependencies) bool {
	return i.deps.union(deps)
}

func (i *Instance) isAfterStatus(status int8) bool {
	return i.status > status
}
