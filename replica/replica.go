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

type InstanceId uint64

type Replica struct {
	Id             int
	Size           int
	MaxInstanceNum []InstanceId // the highest instance number seen for each replica
	InstanceMatrix [][]*Instance
	StateMac       command.StateMachine
	Epoch          uint32
}

func startNewReplica(replicaId, size int) (r *Replica) {
	r = &Replica{
		Id:             replicaId,
		Size:           size,
		MaxInstanceNum: make([]InstanceId, size),
		InstanceMatrix: make([][]*Instance, size),
		StateMac:       new(dummySM.DummySM),
		Epoch:          0,
	}

	for i := 0; i < size; i++ {
		r.InstanceMatrix[i] = make([]*Instance, 1024)
		r.MaxInstanceNum[i] = conflictNotFound
	}

	return r
}

func (r *Replica) fastQuorumSize() int {
	return r.Size - 2
}

func (r *Replica) QuorumSize() int {
	return r.Size/2 + 1
}

func (r *Replica) updateMaxInstanceNum(replicaId int, instanceId InstanceId) bool {
	if r.MaxInstanceNum[replicaId] <= instanceId {
		r.MaxInstanceNum[replicaId] = instanceId + 1
	}
	return false
}

func (r *Replica) run() {
}
