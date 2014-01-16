package replica

import (
	"fmt"
)

var _ = fmt.Printf

func (r *Replica) sendPrepare(L int, insId InstanceIdType, messageChan chan Message) {
	inst := r.InstanceMatrix[L][insId]
	if inst == nil {
		// TODO: we need to commit an instance that doesn't exist.
	}
	inst.ballot = makeLargerBallot(inst.ballot)

	prepare := &Prepare{}

	go func() {
		for i := 0; i < r.N; i++ {
			messageChan <- prepare
		}
	}()

}

func (r *Replica) recvPrepare(prepare *Prepare, messageChan chan Message) {
	inst := r.InstanceMatrix[prepare.repId][prepare.insId]
	if inst == nil {
		// reply prepareOK
	} else {
		if prepare.ballot > inst.ballot {
			// reply prepareOK
		} else {
			// reply NACK
		}
	}
}

func (r *Replica) recvPrepareReply(ppReply *PrepareReply, messageChan chan Message) {
	inst := r.InstanceMatrix[ppReply.repId][ppReply.insId]

	if inst == nil {
		// it shouldn't happen
		return
	}

	inst.info.prepareCnt++

	// for all replies, we only need to keep the ones with highest ballot number
	// inst.ppreplies = r.keepmaxballot()

	// majority replies
	if inst.info.prepareCnt >= r.N/2 {
          // if inst.ppreplies.find( committed )
          // else if inst.ppreplies.find( accepted )
          // else if inst.ppreplies ( >= r.N/2, including itself) preaccepted for default balllot
          // else if inst.ppreplies.find( preaccepted )
          // else default: no-op
	}
}
