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

func (r *Replica) recvPrepare(pp *Prepare, messageChan chan Message) {
	inst := r.InstanceMatrix[pp.repId][pp.insId]
	if inst == nil {
		// reply PrepareReply with no-op and invalid status
		pr := &PrepareReply{
			ok:     true,
			ballot: pp.ballot,
			status: -1,  // hardcode, not a best approach
			cmds:   nil, // No-op, should be agreed by state machine
			deps:   nil,
			repId:  pp.repId,
			insId:  pp.insId,
		}
		r.sendPrepareReply(pr, messageChan)
		return
	}
	// we have some info about the instance
	pr := &PrepareReply{
		status: inst.status,
		cmds:   inst.cmds,
		deps:   inst.deps,
		repId:  pp.repId,
		insId:  pp.insId,
	}
	if pp.ballot >= inst.ballot {
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

func (r *Replica) recvPrepareReply(ppReply *PrepareReply, messageChan chan Message) {
	inst := r.InstanceMatrix[ppReply.repId][ppReply.insId]

	if inst == nil {
		// it shouldn't happen
		return
	}

	inst.info.prepareCnt++

	// for all replies, we only need to keep the ones with highest ballot number
	// inst.ppreplies = r.updateMaxBallot()

	// majority replies
	if inst.info.prepareCnt >= r.N/2 {
		// if inst.ppreplies.find( committed )
		// else if inst.ppreplies.find( accepted )
		// else if inst.ppreplies ( >= r.N/2, including itself) preaccepted for default balllot
		// else if inst.ppreplies.find( preaccepted )
		// else default: no-op
	}
}
