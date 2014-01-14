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
	status int8
}

type Commit struct {
	cmds  []cmd.Command
	seq   int
	deps  []InstanceIdType
	repId int
	insId InstanceIdType
	//ballot uint64 // TODO: there should no need for ballot
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
