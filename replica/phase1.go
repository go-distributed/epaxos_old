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
		info:   &InstanceInfo{},
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
	// update
	deps, changed := r.update(preAccept.cmds, preAccept.deps, preAccept.repId)
	// set cmd
	r.InstanceMatrix[preAccept.repId][preAccept.insId] = &Instance{
		cmds:   preAccept.cmds,
		deps:   deps,
		status: preaccepted,
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
	inst := r.InstanceMatrix[r.Id][paOK.insId]
	inst.info.preaccOkCnt++
	inst.info.preaccCnt++
}

func (r *Replica) recvPreAcceptReply(paReply *PreAcceptReply) {
	inst := r.InstanceMatrix[r.Id][paReply.insId]
	inst.info.preaccCnt++

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
