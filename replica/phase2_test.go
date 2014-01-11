package replica

import (
	"testing"

	cmd "github.com/go-epaxos/epaxos/command"
)

func TestSendAccept(t *testing.T) {
	r := startNewReplica(0, 5)
	messageChan := make(chan Message)
	propose := &Propose{
		cmds: []cmd.Command{
			cmd.Command("hello"),
			cmd.Command("world"),
		},
	}
	r.recvPropose(propose, messageChan)
	r.recvPropose(propose, messageChan)

	for i := 0; i < r.N-1; i++ {
		<-messageChan
		<-messageChan
	}
	r.sendAccept(r.Id, 1, messageChan)

	for i := 0; i < r.N/2; i++ {
		a := (<-messageChan).(*Accept)
		if a.cmds[0].Compare(propose.cmds[0]) != 0 ||
			a.cmds[1].Compare(propose.cmds[1]) != 0 {
			t.Fatal("command isn't equal")
		}
		if a.insId != 1 {
			t.Fatal("instance id should be 1")
		}
		if a.seq != 1 {
			t.Fatal("seq should be 1")
		}
		if a.deps[0] != 0 {
			t.Fatal("deps[0] should be 0")
		}
	}

}
