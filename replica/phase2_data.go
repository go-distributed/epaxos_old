package replica

import (
	cmd "github.com/go-epaxos/epaxos/command"
)

type Accept struct {
	cmds       []cmd.Command
	seq        int
	deps       []InstanceId
	replicaId  int
	instanceId InstanceId
	ballot     *Ballot
}

type AcceptReply struct {
	ok         bool
	replicaId  int
	instanceId InstanceId
	ballot     *Ballot
	status     int8
}

type Commit struct {
	cmds       []cmd.Command
	seq        int
	deps       []InstanceId
	replicaId  int
	instanceId InstanceId
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
