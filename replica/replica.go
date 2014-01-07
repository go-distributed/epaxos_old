package replica

import (
	"fmt"
)

var _ = fmt.Printf

type Replica struct {
	Id             int
	N              int
	InstanceNo     []uint64 // the highes instance number seen for each replica
	InstanceMatrix [][]*Instance

	// channels for receiving messages
	ProposeChan   chan *Propose
	PreAcceptChan chan *PreAccept
}

func startNewReplica(newid, N int) (r *Replica) {
	r = &Replica{
		Id:             newid,
		N:              N,
		InstanceNo:     make([]uint64, N),
		InstanceMatrix: make([][]*Instance, N),
		// channels for receiving messages
		ProposeChan:   make(chan *Propose),
		PreAcceptChan: make(chan *PreAccept),
	}

	for i := 0; i < N; i++ {
		r.InstanceMatrix[i] = make([]*Instance, 1024)
	}

	go r.run()

	return r
}

func (r *Replica) run() {
	for {
		select {
		case propose := <-r.ProposeChan:
			r.recvPropose(propose)
		case preAccept := <-r.PreAcceptChan:
			r.recvPreAccept(preAccept)

		}
	}
}
