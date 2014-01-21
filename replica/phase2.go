package replica

// TODO persistent store

import (
	"strconv"
)

// instId is an id of the already updated instance
// messageChan is a toy channel for emulating broadcast
// TODO: return error
func (r *Replica) sendAccept(replicaId int, instanceId InstanceId, messageChan chan Message) {
	inst := r.InstanceMatrix[replicaId][instanceId]
	if inst == nil {
		// shouldn't get here
		panic("shouldn't get here, replicaId = " + strconv.Itoa(replicaId) + " instanceId = " + strconv.Itoa(int(instanceId)))
	}
	inst.status = accepted
	inst.info.acceptCount = 0
	inst.info.acceptNackCount = 0

	// TODO: persistent store the status
	accept := &Accept{
		cmds:       inst.cmds,
		deps:       inst.deps,
		replicaId:  replicaId,
		instanceId: instanceId,
		ballot:     inst.ballot,
	}

	// TODO: handle timeout
	for i := 0; i < r.Size/2; i++ {
		go func() {
			messageChan <- accept
		}()
	}
}

func (r *Replica) recvAccept(ac *Accept, messageChan chan Message) {
	// TODO: remember to discuss on MaxInstanceNum
	if r.MaxInstanceNum[ac.replicaId] <= ac.instanceId {
		r.MaxInstanceNum[ac.replicaId] = ac.instanceId + 1
	}

	inst := r.InstanceMatrix[ac.replicaId][ac.instanceId]

	if inst == nil {
		r.InstanceMatrix[ac.replicaId][ac.instanceId] = &Instance{
			cmds:   ac.cmds,
			deps:   ac.deps,
			ballot: ac.ballot,
			status: accepted,
			info:   NewInstanceInfo(),
		}
		inst = r.InstanceMatrix[ac.replicaId][ac.instanceId] // for the reference in below
	} else {
		if inst.status >= accepted || ac.ballot.Compare(inst.ballot) < 0 {
			// return nack with status
			ar := &AcceptReply{
				ok:         false,
				replicaId:  ac.replicaId,
				instanceId: ac.instanceId,
				ballot:     inst.ballot,
				status:     inst.status,
			}
			r.sendAcceptReply(ar, messageChan)
			return
		} else {
			inst.cmds = ac.cmds
			inst.deps = ac.deps
			inst.status = accepted
			inst.ballot = ac.ballot
		}
	}

	// reply with ok
	ar := &AcceptReply{
		ok:         true,
		ballot:     inst.ballot,
		replicaId:  ac.replicaId,
		instanceId: ac.instanceId,
	}
	r.sendAcceptReply(ar, messageChan) // should not block
}

func (r *Replica) sendAcceptReply(ar *AcceptReply, messageChan chan Message) {
	go func() {
		messageChan <- ar
	}()
}

func (r *Replica) recvAcceptReply(ar *AcceptReply, messageChan chan Message) {
	inst := r.InstanceMatrix[ar.replicaId][ar.instanceId]
	if inst == nil {
		// TODO: should not get here
		panic("shouldn't get here, replicaId = " + strconv.Itoa(ar.replicaId) + " instanceId = " + strconv.Itoa(int(ar.instanceId)))
	}

	if inst.status > accepted {
		// we've already moved on, this reply is a delayed one
		// so just ignore it
		// TODO: maybe we can have some warning message here
		return
	}

	if !ar.ok { // there must be another proposer, so let's keep quiet
		return
	}

	if ar.ok {
		inst.info.acceptCount++
		if inst.info.acceptCount >= (r.Size / 2) {
			r.sendCommit(ar.replicaId, ar.instanceId, messageChan)
		}
	}
}

func (r *Replica) sendCommit(replicaId int, instanceId InstanceId, messageChan chan Message) {
	inst := r.InstanceMatrix[replicaId][instanceId]
	if inst == nil {
		// shouldn't get here
		panic("shouldn't get here, replicaId = " + strconv.Itoa(replicaId) + " instanceId = " + strconv.Itoa(int(instanceId)))
	}

	inst.status = committed
	// TODO: persistent store

	// make Commit message and send it to all
	cm := &Commit{
		cmds:       inst.cmds,
		deps:       inst.deps,
		replicaId:  replicaId,
		instanceId: instanceId,
	}
	for i := 0; i < r.Size-1; i++ {
		go func() {
			messageChan <- cm
		}()
	}
}

func (r *Replica) recvCommit(cm *Commit) {
	// TODO: remember to discuss on MaxInstanceNum
	if r.MaxInstanceNum[cm.replicaId] <= cm.instanceId {
		r.MaxInstanceNum[cm.replicaId] = cm.instanceId + 1
	}

	inst := r.InstanceMatrix[cm.replicaId][cm.instanceId]
	if inst == nil {
		r.InstanceMatrix[cm.replicaId][cm.instanceId] = &Instance{
			cmds:   cm.cmds,
			deps:   cm.deps,
			status: committed,
			info:   NewInstanceInfo(),
		}
	} else {
		if inst.status >= committed { // ignore the message
			return
		}

		// record the info
		inst.cmds = cm.cmds
		inst.deps = cm.deps
		inst.status = committed
		// TODO: persistence store
	}
}
