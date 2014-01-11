package replica

// instId is an id of the already updated instance
// messageChan is a toy channel for emulating broadcast
// TODO: return error
func (r *Replica) sendAccept(repId int, insId InstanceIdType, messageChan chan Message) {
	inst := r.InstanceMatrix[repId][insId]
	accept := &Accept{
		cmds:  inst.cmds,
		seq:   inst.seq,
		deps:  inst.deps,
		repId: repId,
		insId: insId,
	}

	// TODO: handle timeout
	for i := 0; i < r.N/2; i++ {
		go func() {
			messageChan <- accept
		}()
	}
}
