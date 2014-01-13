// TODO persistent log
package replica

// instId is an id of the already updated instance
// messageChan is a toy channel for emulating broadcast
// TODO: return error
func (r *Replica) sendAccept(repId int, insId InstanceIdType, messageChan chan Message) {
	inst := r.InstanceMatrix[repId][insId]
	accept := &Accept{
		cmds: inst.cmds,
		//seq:   inst.seq,
		deps:   inst.deps,
		repId:  repId,
		insId:  insId,
		ballot: inst.ballot,
	}

	// TODO: handle timeout
	for i := 0; i < r.N/2; i++ {
		go func() {
			messageChan <- accept
		}()
	}
}

func (r *Replica) recvAccept(ac *Accept, messageChan chan Message) {
	inst := r.InstanceMatrix[ac.repId][ac.insId]

	if inst == nil {
		r.InstanceMatrix[ac.repId][ac.insId] = &Instance{
			cmds: ac.cmds,
			//seq: inst.seq,
			deps:   ac.deps,
			status: accepted,
			ballot: ac.ballot,
			info:   &InstanceInfo{},
		}
	} else {
		if ac.ballot < inst.ballot {
			// return nack
			ar := &AcceptReply{
				ok:     false,
				ballot: inst.ballot,
				repId:  ac.repId,
				insId:  ac.insId,
			}
			r.sendAcceptReply(ar, messageChan)
			return
		} else {
			inst.cmds = ac.cmds
			//inst.seq = ac.seq
			inst.deps = ac.deps
			inst.status = accepted
			inst.ballot = ac.ballot
		}
	}

	// reply OK
	ar := &AcceptReply{
		ok:     true,
		ballot: inst.ballot,
		repId:  ac.repId,
		insId:  ac.insId,
	}
	r.sendAcceptReply(ar, messageChan)

}

func (r *Replica) sendAcceptReply(ar *AcceptReply, messageChan chan Message) {
	messageChan <- ar
}

func (r *Replica) recvAcceptReply(ar *AcceptReply, messageChan chan Message) {
	inst := r.InstanceMatrix[ar.repId][ar.insId]
	if inst == nil {
		// TODO: should not get here
	}

	if !ar.ok {
		// TODO: what if receive new epoch?
		b := (ar.ballot & ballotBallotMask) | (inst.ballot & (^ballotBallotMask))
		inst.ballot = makeLargerBallot(b)
		// re-send prepare
	}

	if ar.ok {
		inst.info.acceptOkCnt++
		if inst.info.acceptOkCnt >= (r.N / 2) {
			// ok, let's try to send commit
			inst.status = committed
			// TODO persistent store
			// send commit
		}
	}
}

func (r *Replica) sendCommit(cm *Commit, messageChan chan Message) {

}
