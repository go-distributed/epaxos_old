package replica

import (
	cmd "github.com/go-epaxos/epaxos/command"
)

type Accept struct {
	cmds       []cmd.Command
	seq        int
	deps       dependencies
	replicaId  int
	instanceId InstanceId
	ballot     *Ballot
}

func newAccept(instance *Instance, replicaId int, instanceId InstanceId) *Accept {
	return &Accept{
		cmds:       instance.cmds,
		deps:       instance.deps,
		replicaId:  replicaId,
		instanceId: instanceId,
		ballot:     instance.ballot,
	}
}

type AcceptReply struct {
	ok         bool
	replicaId  int
	instanceId InstanceId
	ballot     *Ballot
	status     int8
}

func newAcceptReply(instance *Instance, replicaId int, instanceId InstanceId, ok bool) *AcceptReply {
	return &AcceptReply{
		ok:         ok,
		replicaId:  replicaId,
		instanceId: instanceId,
		ballot:     instance.ballot,
		status:     instance.status,
	}
}

type Commit struct {
	cmds       []cmd.Command
	seq        int
	deps       []InstanceId
	replicaId  int
	instanceId InstanceId
}

func newCommit(instance *Instance, replicaId int, instanceId InstanceId) *Commit {
	return &Commit{
		cmds:       instance.cmds,
		deps:       instance.deps,
		replicaId:  replicaId,
		instanceId: instanceId,
	}
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
