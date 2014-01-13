package replica

import (
	"fmt"
	"github.com/go-epaxos/epaxos/command"
	"github.com/go-epaxos/epaxos/command/dummySM"
)

var _ = fmt.Printf

const (
	conflictNotFound = 0
)

type InstanceIdType uint64

type Replica struct {
	Id             int
	N              int
	MaxInstanceNum []InstanceIdType // the highes instance number seen for each replica
	InstanceMatrix [][]*Instance
	StateMac       command.StateMachine
	epoch          int
}

func startNewReplica(repId, N int) (r *Replica) {
	r = &Replica{
		Id:             repId,
		N:              N,
		MaxInstanceNum: make([]InstanceIdType, N),
		InstanceMatrix: make([][]*Instance, N),
		StateMac:       new(dummySM.DummySM),
	}

	for i := 0; i < N; i++ {
		r.InstanceMatrix[i] = make([]*Instance, 1024)
		r.MaxInstanceNum[i] = conflictNotFound + 1
	}

	return r
}

func (r *Replica) run() {
}
