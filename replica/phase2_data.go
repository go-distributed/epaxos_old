package replica

import (
	cmd "github.com/go-epaxos/epaxos/command"
)

type Accept struct {
	cmds   []cmd.Command
	seq    int
	deps   []InstanceIdType
	repId  int
	insId  InstanceIdType
	ballot uint64
}

type AcceptReply struct {
	ok     bool
	repId  int
	insId  InstanceIdType
	ballot uint64
	status uint64
}

type Commit struct {
	cmds []cmds.Command
	seq int
	deps []InstanceIdType
	repId int
	insId InstanceIdType
	ballot uint64
}

func (a *Accept) getType() uint8 {
	return acceptType
}

func (a *AcceptReply) getType() uint8 {
	return acceptReplyType
}

func (c *Commit) getType() uint8 {
	return commitType
}
