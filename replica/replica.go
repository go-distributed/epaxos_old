package replica

import (
	"fmt"
	"github.com/go-epaxos/epaxos/command"
	"github.com/go-epaxos/epaxos/command/dummySM"
)

var _ = fmt.Printf

type InstanceIdType int32

type Replica struct {
	Id             int
	N              int
	MaxInstanceNum []InstanceIdType // the highes instance number seen for each replica
	InstanceMatrix [][]*Instance
	StateMac       command.StateMachine
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
	}

	return r
}

func (r *Replica) run() {
}
