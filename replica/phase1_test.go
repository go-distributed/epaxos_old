package replica

import (
	"fmt"
	"testing"

	cmd "github.com/go-epaxos/epaxos/command"
)

var _ = fmt.Printf

func TestRecvPropose(t *testing.T) {
	r := startNewReplica(0, 5)
	messageChan := make(chan Message)
	propose := &Propose{
		cmds: []cmd.Command{
			cmd.Command("hello"),
			cmd.Command("world"),
		},
	}
	r.recvPropose(propose, messageChan)

	for i := 0; i < r.fastQuorumSize(); i++ {
		message := <-messageChan
		preAccept := message.(*PreAccept)

		if preAccept.cmds[0].Compare(propose.cmds[0]) != 0 ||
			preAccept.cmds[1].Compare(propose.cmds[1]) != 0 {
			t.Fatal("command isn't equal")
		}
	}

	// check deps, and instance id
	propose = &Propose{
		cmds: []cmd.Command{
			cmd.Command("hello"),
			cmd.Command("world"),
		},
	}
	r.recvPropose(propose, messageChan)

	for i := 0; i < r.fastQuorumSize(); i++ {
		message := <-messageChan
		preAccept := message.(*PreAccept)

		if preAccept.cmds[0].Compare(propose.cmds[0]) != 0 ||
			preAccept.cmds[1].Compare(propose.cmds[1]) != 0 {
			t.Fatal("command isn't equal")
		}
		if preAccept.insId != conflictNotFound+2 {
			t.Fatal("instance id is wrong")
		}
		if preAccept.deps[0] != conflictNotFound+1 {
			t.Fatal("deps[0] is wrong")
		}
	}
}

func TestRecvPreAccept(t *testing.T) {
	r := startNewReplica(0, 5)
	messageChan := make(chan Message)

	preAccept1 := &PreAccept{
		cmds:  []cmd.Command{cmd.Command("hello")},
		deps:  make([]InstanceIdType, 5),
		repId: 1,
		insId: conflictNotFound + 1,
	}

	r.recvPreAccept(preAccept1, messageChan)

	message := <-messageChan
	if message.getType() != preAcceptOKType {
		t.Fatal("return type should be preAcceptOK")
	}

	preAccept2 := &PreAccept{
		cmds:  []cmd.Command{cmd.Command("hello")},
		deps:  make([]InstanceIdType, 5),
		repId: 2,
		insId: conflictNotFound + 1,
	}

	r.recvPreAccept(preAccept2, messageChan)

	message = <-messageChan
	if message.getType() != preAcceptReplyType {
		t.Fatal("return type should be preAcceptReply")
	}

	paReply := message.(*PreAcceptReply)
	if paReply.deps[1] != conflictNotFound+1 {
		t.Fatal("deps is wrong")
	}

	preAccept3 := &PreAccept{
		cmds:  []cmd.Command{cmd.Command("world")},
		deps:  make([]InstanceIdType, 5),
		repId: 2,
		insId: conflictNotFound + 2,
	}

	r.recvPreAccept(preAccept3, messageChan)

	message = <-messageChan
	if message.getType() != preAcceptOKType {
		t.Fatal("return type should be preAcceptOK")
	}

}

func TestUnion(t *testing.T) {
	r := startNewReplica(0, 3)
	deps1 := []InstanceIdType{0, 0, 0}
	deps2 := []InstanceIdType{0, 0, 0}

	_, same := r.union(deps1, deps2)
	if !same {
		t.Fatal("two deps should be the same")
	}

	deps1[0] = 1
	deps1, same = r.union(deps1, deps2)
	if same || deps1[0] != 1 {
		t.Fatal("wrong deps after union")
	}

	deps2[0] = 2
	deps1, same = r.union(deps1, deps2)
	if same || deps1[0] != 2 {
		t.Fatal("wrong deps after union")
	}

}
