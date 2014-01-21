package replica

// TODO persistent store

import (
	"fmt"
)

// instId is an id of the already updated instance
// messageChan is a toy channel for emulating broadcast
// TODO: return error
func (r *Replica) sendAccept(replicaId int, instanceId InstanceId, messageChan chan Message) {
	instance := r.InstanceMatrix[replicaId][instanceId]
	if instance == nil {
		// shouldn't get here
		msg := fmt.Sprintf("shouldn't get here, replicaId = %d, instanceId = %d",
			replicaId, instanceId)
		panic(msg)
	}

	if instance.status != accepted {
		msg := fmt.Sprintf("[sendAccept] unaccepted instance %v", instanceId)
		panic(msg)
	}

	// TODO: persistent store the status
	accept := newAccept(instance, replicaId, instanceId)

	// TODO: handle timeout
	for i := 0; i < r.QuorumSize()-1; i++ {
		go func() {
			messageChan <- accept
		}()
	}
}

func (r *Replica) recvAccept(accept *Accept, messageChan chan Message) {
	ok := true

	rId, iId := accept.replicaId, accept.instanceId
	// TODO: remember to discuss on MaxInstanceNum
	r.updateMaxInstanceNum(rId, iId)

	i := r.InstanceMatrix[rId][iId]
	if i == nil {
		i = &Instance{
			cmds:   accept.cmds,
			deps:   accept.deps,
			ballot: accept.ballot,
			status: accepted,
			info:   NewInstanceInfo(),
		}
		r.InstanceMatrix[rId][iId] = i
	} else {
		ok = i.processAccept(accept)
	}

	ar := newAcceptReply(i, rId, iId, ok)
	r.sendAcceptReply(ar, messageChan) // should not block
}

func (r *Replica) sendAcceptReply(ar *AcceptReply, messageChan chan Message) {
	go func() {
		messageChan <- ar
	}()
}

func (r *Replica) recvAcceptReply(ar *AcceptReply, messageChan chan Message) {
	i := r.InstanceMatrix[ar.replicaId][ar.instanceId]
	if i == nil {
		// TODO: should not get here
		msg := fmt.Sprintf("shouldn't get here, replicaId = %d, instanceId = %d",
			ar.replicaId, ar.instanceId)
		panic(msg)
	}

	if !ar.ok { // there must be another proposer, so let's keep quiet
		return
	}

	if ok := i.processAcceptReply(ar, r.QuorumSize()); ok {
		r.sendCommit(ar.replicaId, ar.instanceId, messageChan)
	}
}

func (r *Replica) sendCommit(replicaId int, instanceId InstanceId, messageChan chan Message) {
	instance := r.InstanceMatrix[replicaId][instanceId]
	if instance == nil {
		msg := fmt.Sprintf("shouldn't get here, replicaId = %d, instanceId = %d",
			replicaId, instanceId)
		panic(msg)
	}

	if instance.status != committed {
		msg := fmt.Sprintf("[sendCommit] uncommitted instance %v", instanceId)
		panic(msg)
	}

	// TODO: persistent store

	// make Commit message and send it to all
	cm := newCommit(instance, replicaId, instanceId)
	for i := 0; i < r.Size-1; i++ {
		go func() {
			messageChan <- cm
		}()
	}
}

func (r *Replica) recvCommit(commit *Commit) {
	rId, iId := commit.replicaId, commit.instanceId

	// TODO: remember to discuss on MaxInstanceNum
	r.updateMaxInstanceNum(rId, iId)

	i := r.InstanceMatrix[rId][iId]
	if i == nil {
		r.InstanceMatrix[rId][iId] = &Instance{
			cmds:   commit.cmds,
			deps:   commit.deps,
			status: committed,
			info:   NewInstanceInfo(),
		}
	} else {
		i.processCommit(commit)
		// TODO: persistence store
	}
}
