package replica

import (
	"fmt"
	cmd "github.com/go-epaxos/epaxos/command"
)

var _ = fmt.Printf

func (r *Replica) recvPropose(propose *Propose, messageChan chan Message) {
	deps := r.findDependencies(propose.cmds)

	// increment instance number
	instNo := r.MaxInstanceNum[r.Id]
	instNo = instNo
	r.MaxInstanceNum[r.Id]++

	// set cmds
	r.InstanceMatrix[r.Id][instNo] = &Instance{
		cmds:   propose.cmds,
		deps:   deps,
		status: preaccepted,
		ballot: r.makeInitialBallot(),
		info:   NewInstanceInfo(),
	}

	// TODO: before we send the message, we need to record and sync it in disk/persistent.

	// send PreAccept
	preAccept := &PreAccept{
		cmds:  propose.cmds,
		deps:  deps,
		repId: r.Id,
		insId: instNo,
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
	deps, changed := r.updateDependencies(preAccept.cmds, preAccept.deps, preAccept.repId)
	// set cmd
	r.InstanceMatrix[preAccept.repId][preAccept.insId] = &Instance{
		cmds:   preAccept.cmds,
		deps:   deps,
		status: preaccepted,
		info:   NewInstanceInfo(),
	}
	if preAccept.insId >= r.MaxInstanceNum[preAccept.repId] {
		r.MaxInstanceNum[preAccept.repId] = preAccept.insId + 1
	}
	// reply
	go func() {
		if !changed {
			paOK := &PreAcceptOK{
				insId: preAccept.insId,
			}
			messageChan <- paOK
		} else {
			paReply := &PreAcceptReply{
				deps:  deps,
				repId: preAccept.repId,
				insId: preAccept.insId,
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

func (r *Replica) recvPreAcceptReply(paReply *PreAcceptReply) {
	inst := r.InstanceMatrix[paReply.repId][paReply.insId]

	if inst == nil {
		// TODO: should not happen
		return
	}

	// when we receive a reply which is not what we expect (sometimes the status
	// is later than we thought), it's probably because the instance has been delayed/partitioned
	// for a period of time. And after it comes back, things have changed. It's been accepted, or
	// committed, or even executed. Since it's accepted by majority already, we ignore it here and hope
	// the instance would be fixed later (by asking dependencies when executing commands).
	if inst.status > preaccepted {
		// TODO: slow reply
		return
	}

	inst.info.preaccCnt++

	// recvpreacceptok doesn't need this {
	deps, same := r.union(inst.deps, paReply.deps)
	if !same {
		if inst.info.preaccCnt > 1 {
			inst.info.haveDiffReply = true
		}
		inst.deps = deps
	}
	// }

	if inst.info.preaccCnt >= r.Size/2 && !inst.allReplyTheSame() {
		// slow path

	} else if inst.info.preaccCnt == r.fastQuorumSize() && inst.allReplyTheSame() {
		// fast path
	}
}

// findDependencies finds the most recent interference instance from each instance space
// of this replica.
// It returns the ids of these instances as an array.
func (r *Replica) findDependencies(cmds []cmd.Command) []InstanceIdType {
	deps := make([]InstanceIdType, r.Size)

	for i := range r.InstanceMatrix {
		instances := r.InstanceMatrix[i]
		start := r.MaxInstanceNum[i] - 1

		if conflict, ok := r.scanConflicts(instances, cmds, start, 0); ok {
			deps[i] = conflict
		}
	}

	return deps
}

func (r *Replica) updateDependencies(cmds []cmd.Command, deps []InstanceIdType, from int) ([]InstanceIdType, bool) {
	changed := false

	for curr := range r.InstanceMatrix {
		// short cut here, the sender knows the latest dependencies for itself
		if curr == from {
			continue
		}

		instances := r.InstanceMatrix[curr]
		start, end := r.MaxInstanceNum[curr]-1, deps[curr]

		if conflict, ok := r.scanConflicts(instances, cmds, start, end); ok {
			changed = true
			deps[curr] = conflict
		}
	}

	return deps, changed
}

func (r *Replica) scanConflicts(instances []*Instance, cmds []cmd.Command, start InstanceIdType, end InstanceIdType) (InstanceIdType, bool) {
	for id := start; id > end; id-- {
		if instances[id] == nil {
			continue
		}
		// we only need to find the highest instance in conflict
		if r.StateMac.HaveConflicts(cmds, instances[id].cmds) {
			return InstanceIdType(id), true
		}
	}

	return 0, false
}

func (r *Replica) union(deps1, deps2 []InstanceIdType) ([]InstanceIdType, bool) {
	same := true
	for rep := 0; rep < r.Size; rep++ {
		if deps1[rep] != deps2[rep] {
			same = false
			if deps1[rep] < deps2[rep] {
				deps1[rep] = deps2[rep]
			}
		}
	}
	return deps1, same
}
