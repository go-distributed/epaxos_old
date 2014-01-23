package replica

import (
	"fmt"
)

var _ = fmt.Printf

func (r *Replica) sendPrepare(replicaId int, instanceId InstanceId, messageChan chan Message) {
	if r.InstanceMatrix[replicaId][instanceId] == nil {
		// TODO: we need to prepare an instance that doesn't exist.
		r.InstanceMatrix[replicaId][instanceId] = &Instance{
			// TODO:
			// Assumed no-op to be nil here.
			// we need to do more since state machine needs to know how to interpret it.
			cmds:   nil,
			deps:   make([]InstanceId, r.Size), // TODO: makeInitialDeps
			status: -1,                         // 'none' might be a conflicting name. We currenctly pick '-1' for it
			ballot: r.makeInitialBallot(),
			info:   NewInstanceInfo(),
		}
	}

	inst := r.InstanceMatrix[replicaId][instanceId]

	// clean preAccepCount
	inst.recoveryInfo = NewRecoveryInfo()
	inst.recoveryInfo.status = preAccepted
	inst.ballot.incNumber()
	inst.ballot.setReplicaId(r.Id)

	prepare := &Prepare{
		ballot:     inst.ballot,
		replicaId:  replicaId,
		instanceId: instanceId,
	}

	go func() {
		for i := 0; i < r.Size-1; i++ {
			messageChan <- prepare
		}
	}()
}

func (r *Replica) recvPrepare(pp *Prepare, messageChan chan Message) {
	inst := r.InstanceMatrix[pp.replicaId][pp.instanceId]
	if inst == nil {
		// reply PrepareReply with no-op and invalid status
		pr := &PrepareReply{
			ok:         true,
			ballot:     r.makeInitialBallot(),
			status:     -1,                         // TODO: hardcode, not a best approach
			deps:       make([]InstanceId, r.Size), // TODO: makeInitialDeps
			replicaId:  pp.replicaId,
			instanceId: pp.instanceId,
		}
		r.sendPrepareReply(pr, messageChan)
		return
	}

	// we have some info about the instance
	pr := &PrepareReply{
		status:     inst.status,
		replicaId:  pp.replicaId,
		instanceId: pp.instanceId,
		cmds:       inst.cmds,
		deps:       inst.deps,
		ballot:     inst.ballot,
	}

	// we won't have the same ballot
	if pp.ballot.Compare(inst.ballot) > 0 {
		pr.ok = true
		inst.ballot = pp.ballot
	} else {
		pr.ok = false
	}

	r.sendPrepareReply(pr, messageChan)
}

func (r *Replica) sendPrepareReply(pr *PrepareReply, messageChan chan Message) {
	go func() {
		messageChan <- pr
	}()
}

func (r *Replica) recvPrepareReply(p *PrepareReply, m chan Message) {
	inst := r.InstanceMatrix[p.replicaId][p.instanceId]

	if inst == nil {
		// it shouldn't happen
		return
	}
	if !inst.isAtStatus(preparing) {
		// this is a delayed message. ignore it
		return
	}
	// once we receive a "commited" reply,
	// then we can leave preparing state and send commits
	if p.status == committed {
		inst.cmds, inst.deps, inst.status = p.cmds, p.deps, committed
		r.sendCommit(p.replicaId, p.instanceId, m)
		return
	}
	if !p.ok {
		return
	}

	rInfo := inst.recoveryInfo
	rInfo.replyCount++

	// handle differnt replies
	// once we receive an "accepted" reply,
	// then we can send Accepts, but we need to send the most recent one
	if p.status == accepted {
		rInfo.status = accepted

		// only record the most recent accepted instance
		if p.ballot.Compare(rInfo.maxAcceptBallot) > 0 {
			rInfo.maxAcceptBallot, rInfo.cmds, rInfo.deps = p.ballot, p.cmds, p.deps
		}
	}
	// if we receive a "preAccepted" reply, and we haven't receive any "accepted" replies
	// then we should check if we can(not) send Accept
	if p.status == preAccepted && rInfo.status < accepted {
		rInfo.status = preAccepted
		rInfo.preAcceptCount++

		// if violate the "N/2 identical replies" condition {
		//         union deps
		// }
	}

	// wait to receive enough replies
	if rInfo.replyCount < r.QuorumSize()-1 {
		return
	}

	// send the most recent Accept
	if rInfo.status == accepted {
		inst.cmds, inst.deps = rInfo.cmds, rInfo.deps
		r.sendAccept(p.replicaId, p.instanceId, m)
		return
	}

	// send Accepts or PreAccepts
	if rInfo.status == preAccepted {
		// if "N/2 identical" {
		//         sendAccept()
		//         return
		// }
		// sendPreAccept()

		return
	}

	// sendPreAccept(noop)
	// [*] I forgot why we can't send accept noop here...
}
