package replica

import (
	"fmt"

	cmd "github.com/go-epaxos/epaxos/command"
)

var _ = fmt.Printf

func (r *Replica) sendPrepare(L int, instanceId InstanceId, messageChan chan Message) {
	if r.InstanceMatrix[L][instanceId] == nil {
		// TODO: we need to commit an instance that doesn't exist.
		r.InstanceMatrix[L][instanceId] = &Instance{
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

	inst := r.InstanceMatrix[L][instanceId]

	inst.ballot.incNumber()
	inst.ballot.setRId(r.Id)

	prepare := &Prepare{
		ballot:     inst.ballot,
		replicaId:  L,
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
			ballot:     pp.ballot,
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
		cmds:       inst.cmds,
		deps:       inst.deps,
		replicaId:  pp.replicaId,
		instanceId: pp.instanceId,
	}
	if pp.ballot.Compare(inst.ballot) >= 0 {
		pr.ok = true
		inst.ballot = pp.ballot
	} else {
		pr.ok = false
	}
	pr.ballot = inst.ballot
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

	inst.info.prepareCount++

	// for all replies, we only need to keep the ones with highest ballot number
	// inst.ppreplies = r.updateMaxBallot()

	// majority replies
	if inst.info.prepareCount >= r.Size/2 {
		// if inst.ppreplies.find( committed )
		// else if inst.ppreplies.find( accepted )
		// else if inst.ppreplies ( >= r.Size/2, including itself) preaccepted for default balllot
		// else if inst.ppreplies.find( preaccepted )
		// else default: no-op
	}
}

func (r *Replica) updateRecovery(cmds []cmd.Command, deps []InstanceId, status int) {

}
