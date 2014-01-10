package replica

import (
	"bytes"
	"fmt"
	"testing"

	cmd "github.com/go-epaxos/epaxos/command"
)

var _ = fmt.Printf

func TestRecvPropose(t *testing.T) {
	r := startNewReplica(0, 5)
	messageChan := make(chan Message)
	propose := &Propose{
		Cmds: []cmd.Command{
			cmd.Command("hello"),
			cmd.Command("world"),
		},
	}
	r.recvPropose(propose, messageChan)

	for i := 0; i < r.N-1; i++ {
		message := <-messageChan
		preAccept := message.(*PreAccept)

		if bytes.Compare(preAccept.Cmds[0], propose.Cmds[0]) != 0 ||
			bytes.Compare(preAccept.Cmds[1], propose.Cmds[1]) != 0 {
			t.Fatal("command isn't equal")
		}
	}
}
