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
			if r.StateMac.HaveConflicts(propose.Cmds, repInst[InstId].Cmds) {
				if seq <= repInst[InstId].Seq {
					seq = repInst[InstId].Seq + 1
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
		Cmds:   propose.Cmds,
		Seq:    seq,
		Deps:   deps,
		Status: preaccepted,
	}

	// TODO: before we send messages, we need to record and sync it in disk/persistent.

	// send PreAccept
	preAccept := &PreAccept{
		Cmds:  propose.Cmds,
		Seq:   seq,
		Deps:  deps,
		repId: r.Id,
		InsId: instNo,
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
