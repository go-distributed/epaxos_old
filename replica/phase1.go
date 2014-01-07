package replica

import (
	"fmt"
)

var _ = fmt.Printf

type Propose struct {
}

type PreAccept struct {
}

func (r *Replica) recvPropose(propose *Propose) {
	// increment instance number
        instNo := r.InstanceNo[r.Id]
        instNo = instNo
	r.InstanceNo[r.Id]++
	// update seq, deps

	// set cmds
        r.InstanceMatrix[r.Id][instNo] = &Instance {
        }
	// send PreAccept
}

func (r *Replica) recvPreAccept(preAccept *PreAccept) {
	return
}
