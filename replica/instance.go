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
	preAcceptCount int
	isFastPath     bool

	acceptNackCount int
	acceptCount     int

	prepareCount int

	recovery *RecoveryInfo
}

type RecoveryInfo struct {
	preAcceptedCount int
}

type Instance struct {
	cmds   []cmd.Command
	deps   dependencies
	status int8
	ballot *Ballot
	info   *InstanceInfo
}

func NewRecoveryInfo() *RecoveryInfo {
	return &RecoveryInfo{}
}

func NewInstanceInfo() *InstanceInfo {
	return &InstanceInfo{
		isFastPath: true,
	}
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
			i.setFastPath(false)
		}
	}

	if i.preAcceptCount() >= quorum-1 && !i.isFastPath() {
		i.status = accepted
	}

	if i.preAcceptCount() == fastQuorum && i.isFastPath() {
		i.status = committed
	}

	return i.status
}

func (i *Instance) processAccept(a *Accept) bool {
	if i.isEqualOrAfterStatus(accepted) || a.ballot.Compare(i.ballot) < 0 {
		return false
	}

	i.status = accepted
	i.cmds = a.cmds
	i.deps = a.deps
	i.ballot = a.ballot
	return true
}

func (i *Instance) processAcceptReply(ar *AcceptReply, quorum int) bool {
	if i.isAfterStatus(accepted) {
		// we've already moved on, this reply is a delayed one
		// so just ignore it
		// TODO: maybe we can have some warning message here
		return false
	}
	i.info.acceptCount++
	if i.info.acceptCount >= quorum-1 {
		i.status = committed
		return true
	}
	return false
}

func (i *Instance) processCommit(c *Commit) bool {
	if i.status >= committed { // ignore the message
		return false
	}

	i.cmds = c.cmds
	i.deps = c.deps
	i.status = committed
	return true
}

func (i *Instance) preAcceptCount() int {
	return i.info.preAcceptCount
}

func (i *Instance) isFastPath() bool {
	return i.info.isFastPath
}

func (i *Instance) setFastPath(ok bool) {
	i.info.isFastPath = ok
}

func (i *Instance) unionDeps(deps dependencies) bool {
	return i.deps.union(deps)
}

func (i *Instance) isAtStatus(status int8) bool {
	return i.status == status
}

func (i *Instance) isAfterStatus(status int8) bool {
	return i.status > status
}

func (i *Instance) isEqualOrAfterStatus(status int8) bool {
	return i.status >= status
}
