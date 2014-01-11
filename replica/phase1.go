package replica

import (
	"fmt"
)

var _ = fmt.Printf

func (r *Replica) recvPropose(propose *Propose, messageChan chan Message) {
	// update seq, deps
	deps := make([]InstanceIdType, r.N)
	seq := 0
	for rep := 0; rep < r.N; rep++ {
		repInst := r.InstanceMatrix[rep]

		InstId := r.MaxInstanceNum[rep] - 1
		for ; InstId >= 0; InstId-- {
			if r.StateMac.HaveConflicts(propose.cmds, repInst[InstId].cmds) {
				if seq <= repInst[InstId].seq {
					seq = repInst[InstId].seq + 1
				}
				break
			}
		}
		// if InstId >= 0, we found conflicted instance;
		// if InstId == -1, we didn't have any interference;
		deps[rep] = InstId
	}

	// increment instance number
	instNo := r.MaxInstanceNum[r.Id]
	instNo = instNo
	r.MaxInstanceNum[r.Id]++

	// set cmds
	r.InstanceMatrix[r.Id][instNo] = &Instance{
		cmds:   propose.cmds,
		seq:    seq,
		deps:   deps,
		status: preaccepted,
	}

	// TODO: before we send messages, we need to record and sync it in disk/persistent.

	// send PreAccept
	preAccept := &PreAccept{
		cmds:  propose.cmds,
		seq:   seq,
		deps:  deps,
		repId: r.Id,
		insId: instNo,
	}

	// fast quorum
	go func() {
		for i := 0; i < r.N-1; i++ {
			messageChan <- preAccept
		}
	}()

}

func (r *Replica) recvPreAccept(preAccept *PreAccept) {
	return
}
