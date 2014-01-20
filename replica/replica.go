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
	Size           int
	MaxInstanceNum []InstanceIdType // the highest instance number seen for each replica
	InstanceMatrix [][]*Instance
	StateMac       command.StateMachine
	Epoch          uint32
}

func startNewReplica(repId, size int) (r *Replica) {
	r = &Replica{
		Id:             repId,
		Size:           size,
		MaxInstanceNum: make([]InstanceIdType, size),
		InstanceMatrix: make([][]*Instance, size),
		StateMac:       new(dummySM.DummySM),
		Epoch:          0,
	}

	for i := 0; i < size; i++ {
		r.InstanceMatrix[i] = make([]*Instance, 1024)
		r.MaxInstanceNum[i] = conflictNotFound + 1
	}

	return r
}

func (r *Replica) fastQuorumSize() int {
	return r.Size - 2
}

func (r *Replica) QuorumSize() int {
	return r.Size/2 + 1
}

func (r *Replica) run() {
}
