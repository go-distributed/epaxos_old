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
	preaccCnt     int
	haveDiffReply bool

	acceptNackCnt int
	acceptOkCnt   int
}

type Instance struct {
	cmds   []cmd.Command
	deps   []InstanceIdType
	status int8
	ballot uint64
	info   *InstanceInfo
}

func (Inst *Instance) allReplyTheSame() bool {
	return Inst.info.haveDiffReply == false
}
