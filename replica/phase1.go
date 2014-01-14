package replica

import (
	"fmt"
	cmd "github.com/go-epaxos/epaxos/command"
)

var _ = fmt.Printf

func (r *Replica) recvPropose(propose *Propose, messageChan chan Message) {
	// update deps
	deps := make([]InstanceIdType, r.N)
	// we need to initiate deps to conflictnotfound. Since it's 0, we keep that assumption.

	deps, _ = r.update(propose.cmds, deps, r.Id)

	// increment instance number
	instNo := r.MaxInstanceNum[r.Id]
	instNo = instNo
	r.MaxInstanceNum[r.Id]++

	// set cmds
	r.InstanceMatrix[r.Id][instNo] = &Instance{
		cmds:   propose.cmds,
		deps:   deps,
		status: preaccepted,
		ballot: makeBallot(uint64(r.epoch), uint64(r.Id)),
		info:   new(InstanceInfo),
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
	deps, changed := r.update(preAccept.cmds, preAccept.deps, preAccept.repId)
	// set cmd
	r.InstanceMatrix[preAccept.repId][preAccept.insId] = &Instance{
		cmds:   preAccept.cmds,
		deps:   deps,
		status: preaccepted,
		info:   &InstanceInfo{},
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
	// It's almost the same as recvpreacceptreply() function. But it doesn't need to
	// union the dependencies.
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
        // committed, or even executed. Since it's done by majority, we just ignore it here and hope
        // the instance would be fixed later (by asking dependencies).
	if inst.status != preaccepted {
		// TODO: slow reply
		return
	}

	inst.info.preaccCnt++

        // only diff with recvpreacceptok() {
	deps, same := r.union(inst.deps, paReply.deps)
	if !same {
		inst.info.haveDiff = true
		inst.deps = deps
	}
	// }

	if inst.info.preaccCnt >= r.N/2 && inst.info.haveDiff {
		// slow path

	} else if inst.info.preaccCnt == r.fastQuorumSize() && !inst.info.haveDiff {
		// fast path
	}
}

func (r *Replica) update(cmds []cmd.Command, deps []InstanceIdType,
	repId int) ([]InstanceIdType, bool) {
	changed := false

	for rep := 0; rep < r.N; rep++ {
		if r.Id != repId && rep == repId {
			continue
		}
		repInst := r.InstanceMatrix[rep]
		InstId := r.MaxInstanceNum[rep] - 1

		for ; InstId > deps[rep]; InstId-- {
			if repInst[InstId] == nil {
				continue
			}
			// we only need to find the highest instance in conflict
			if r.StateMac.HaveConflicts(cmds, repInst[InstId].cmds) {
				changed = true
				break
			}
		}
		// if InstId > original dep, we found newer conflicted instance;
		// if InstId = original dep, we didn't find any new conflict;
		if InstId != deps[rep] {
			deps[rep] = InstId
		}
	}

	return deps, changed
}

func (r *Replica) union(deps1, deps2 []InstanceIdType) ([]InstanceIdType, bool) {
	same := true
	for rep := 0; rep < r.N; rep++ {
		if deps1[rep] != deps2[rep] {
			same = false
			if deps1[rep] < deps2[rep] {
				deps1[rep] = deps2[rep]
			}
		}
	}
	return deps1, same
}
