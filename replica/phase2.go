package replica

// TODO persistent store

import (
	"strconv"
)

// instId is an id of the already updated instance
// messageChan is a toy channel for emulating broadcast
// TODO: return error
func (r *Replica) sendAccept(repId int, insId InstanceIdType, messageChan chan Message) {
	inst := r.InstanceMatrix[repId][insId]
	if inst == nil {
		// shouldn't get here
		panic("shouldn't get here, repId = " + strconv.Itoa(repId) + " insId = " + strconv.Itoa(int(insId)))
	}
	inst.status = accepted

	// TODO: persistent store the status
	accept := &Accept{
		cmds:   inst.cmds,
		deps:   inst.deps,
		repId:  repId,
		insId:  insId,
		ballot: inst.ballot,
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
	if r.MaxInstanceNum[ac.repId] <= ac.insId {
		r.MaxInstanceNum[ac.repId] = ac.insId + 1
	}

	inst := r.InstanceMatrix[ac.repId][ac.insId]

	if inst == nil {
		r.InstanceMatrix[ac.repId][ac.insId] = &Instance{
			cmds:   ac.cmds,
			deps:   ac.deps,
			ballot: ac.ballot,
			status: accepted,
			info:   NewInstanceInfo(),
		}
		inst = r.InstanceMatrix[ac.repId][ac.insId] // for the reference in below
	} else {
		if inst.status >= accepted || ac.ballot.Compare(inst.ballot) < 0 {
			// return nack with status
			ar := &AcceptReply{
				ok:     false,
				repId:  ac.repId,
				insId:  ac.insId,
				ballot: inst.ballot,
				status: inst.status,
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
		ok:     true,
		ballot: inst.ballot,
		repId:  ac.repId,
		insId:  ac.insId,
	}
	r.sendAcceptReply(ar, messageChan) // should not block
}

func (r *Replica) sendAcceptReply(ar *AcceptReply, messageChan chan Message) {
	go func() {
		messageChan <- ar
	}()
}

func (r *Replica) recvAcceptReply(ar *AcceptReply, messageChan chan Message) {
	inst := r.InstanceMatrix[ar.repId][ar.insId]
	if inst == nil {
		// TODO: should not get here
		panic("shouldn't get here, repId = " + strconv.Itoa(ar.repId) + " insId = " + strconv.Itoa(int(ar.insId)))
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
		inst.info.acceptOkCnt++
		if inst.info.acceptOkCnt >= (r.Size / 2) {
			r.sendCommit(ar.repId, ar.insId, messageChan)
		}
	}
}

func (r *Replica) sendCommit(repId int, insId InstanceIdType, messageChan chan Message) {
	inst := r.InstanceMatrix[repId][insId]
	if inst == nil {
		// shouldn't get here
		panic("shouldn't get here, repId = " + strconv.Itoa(repId) + " insId = " + strconv.Itoa(int(insId)))
	}

	inst.status = committed
	// TODO: persistent store

	// make Commit message and send it to all
	cm := &Commit{
		cmds:  inst.cmds,
		deps:  inst.deps,
		repId: repId,
		insId: insId,
	}
	for i := 0; i < r.Size-1; i++ {
		go func() {
			messageChan <- cm
		}()
	}
}

func (r *Replica) recvCommit(cm *Commit) {
	// TODO: remember to discuss on MaxInstanceNum
	if r.MaxInstanceNum[cm.repId] <= cm.insId {
		r.MaxInstanceNum[cm.repId] = cm.insId + 1
	}

	inst := r.InstanceMatrix[cm.repId][cm.insId]
	if inst == nil {
		r.InstanceMatrix[cm.repId][cm.insId] = &Instance{
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
