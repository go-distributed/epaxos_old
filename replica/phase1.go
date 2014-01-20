package replica

import (
	"fmt"

	cmd "github.com/go-epaxos/epaxos/command"
)

var _ = fmt.Printf

func (r *Replica) recvPropose(propose *Propose, messageChan chan Message) {
	deps := r.findDependencies(propose.cmds)

	// increment instance number
	r.MaxInstanceNum[r.Id]++
	instNo := r.MaxInstanceNum[r.Id]

	// set cmds
	r.InstanceMatrix[r.Id][instNo] = &Instance{
		cmds:   propose.cmds,
		deps:   deps,
		status: preAccepted,
		ballot: r.makeInitialBallot(),
		info:   NewInstanceInfo(),
	}

	// TODO: before we send the message, we need to record and sync it in disk/persistent.

	// send PreAccept
	preAccept := &PreAccept{
		cmds:       propose.cmds,
		deps:       deps,
		replicaId:  r.Id,
		instanceId: instNo,
	}

	// fast quorum
	go func() {
		for i := 0; i < r.fastQuorumSize(); i++ {
			messageChan <- preAccept
		}
	}()

}

func (r *Replica) recvPreAccept(preAccept *PreAccept, messageChan chan Message) {
	// TODO: we need to check ballot for that coming from prepare phase
	// update
	deps, changed := r.updateDependencies(preAccept.cmds, preAccept.deps, preAccept.replicaId)
	// set cmd
	r.InstanceMatrix[preAccept.replicaId][preAccept.instanceId] = &Instance{
		cmds:   preAccept.cmds,
		deps:   deps,
		status: preAccepted,
		info:   NewInstanceInfo(),
	}
	if preAccept.instanceId >= r.MaxInstanceNum[preAccept.replicaId] {
		r.MaxInstanceNum[preAccept.replicaId] = preAccept.instanceId + 1
	}
	// reply
	go func() {
		if !changed {
			paOK := &PreAcceptOK{
				instanceId: preAccept.instanceId,
			}
			messageChan <- paOK
		} else {
			paReply := &PreAcceptReply{
				deps:       deps,
				replicaId:  preAccept.replicaId,
				instanceId: preAccept.instanceId,
			}
			messageChan <- paReply
		}
	}()
}

func (r *Replica) recvPreAcceptOK(paOK *PreAcceptOK) {
	// recvpreacceptok() is subset of recvpreacceptreply()
	// It doens't need to do union
	// We need some refactoring between these two functions.
}

// handle preAcceptReply handles the reply for preAccept
// ignore the reply if the instance of the replica has passed the preAccept phase
// update the bookkeeping of the instance according to the reply
// if we receives at least floor(N/2) replies an
func (r *Replica) recvPreAcceptReply(par *PreAcceptReply) {
	instance := r.InstanceMatrix[par.replicaId][par.instanceId]

	if instance == nil {
		// TODO: should not happen
		panic("handlePreAcceptReply: receive nil instance")
		return
	}

	// when we receive a reply which is not what we expect (sometimes the status
	// is later than we thought), it's probably because the instance has been delayed/partitioned
	// for a period of time. And after it comes back, things have changed. It's been accepted, or
	// committed, or even executed. Since it's accepted by majority already, we ignore it here and hope
	// the instance would be fixed later (by asking dependencies when executing commands).
	if instance.isAfterStatus(preAccepted) {
		// TODO: log here.
		return
	}

	switch instance.processPreAcceptReply(par, r.QuorumSize()-1, r.fastQuorumSize()) {
	case committed:
		//fast path
	case accepted:
		// slow path
	}

}

// findDependencies finds the most recent interference instance from each instance space
// of this replica.
// It returns the ids of these instances as an array.
func (r *Replica) findDependencies(cmds []cmd.Command) []InstanceId {
	deps := make([]InstanceId, r.Size)

	for i := range r.InstanceMatrix {
		instances := r.InstanceMatrix[i]
		start := r.MaxInstanceNum[i]

		if conflict, ok := r.scanConflicts(instances, cmds, start, 0); ok {
			deps[i] = conflict
		}
	}

	return deps
}

// updateDependencies updates the passed in dependencies from replica[from].
// return updated dependencies and whether the dependencies has changed.
func (r *Replica) updateDependencies(cmds []cmd.Command, deps []InstanceId, from int) ([]InstanceId, bool) {
	changed := false

	for curr := range r.InstanceMatrix {
		// short cut here, the sender knows the latest dependencies for itself
		if curr == from {
			continue
		}

		instances := r.InstanceMatrix[curr]
		start, end := r.MaxInstanceNum[curr], deps[curr]

		if conflict, ok := r.scanConflicts(instances, cmds, start, end); ok {
			changed = true
			deps[curr] = conflict
		}
	}

	return deps, changed
}

// scanConflicts scans the instances from start to end (high to low).
// return the highest instance that has conflicts with passed in cmds.
func (r *Replica) scanConflicts(instances []*Instance, cmds []cmd.Command, start InstanceId, end InstanceId) (InstanceId, bool) {
	for i := start; i > end; i-- {
		if instances[i] == nil {
			continue
		}
		// we only need to find the highest instance in conflict
		if r.StateMac.HaveConflicts(cmds, instances[i].cmds) {
			return i, true
		}
	}

	return conflictNotFound, false
}
