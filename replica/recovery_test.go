package replica

import (
	"fmt"
	"testing"

	cmd "github.com/go-epaxos/epaxos/command"
)

var _ = fmt.Printf

func TestSendPrepare(t *testing.T) {
	messageChan := make(chan Message)
	// send prepare
	r := startNewReplica(0, 5)
	r.sendPrepare(0, conflictNotFound+1, messageChan)
	// verify the prepare message for:
	// - incremented ballot number
	// - L.i instance
	for i := 0; i < r.N-1; i++ {
		message := <-messageChan
		prepare := message.(*Prepare)
		if prepare.ballot.getNumber() != 1 {
			t.Fatal("expected ballot number to be incremented to 1")
		}
	}

	r.InstanceMatrix[0][conflictNotFound+1].cmds = []cmd.Command{
		cmd.Command("hello"),
	}

	r.sendPrepare(0, conflictNotFound+1, messageChan)
	if r.InstanceMatrix[0][conflictNotFound+1].cmds[0].Compare(cmd.Command("hello")) != 0 {
		t.Fatal("The cmds shouldn't be changed")
	}
	for i := 0; i < r.N-1; i++ {
		message := <-messageChan
		prepare := message.(*Prepare)
		if prepare.ballot.getNumber() != 2 {
			t.Fatal("expected ballot number to be incremented to 2")
		}
	}

}

func TestRecvPrepare(t *testing.T) {
	// make prepare message to test for:
	// - ballot is larger than most recent one
	// - ballot isn't larger
	// - there's no instance!
}

func TestRecvPrepareReply(t *testing.T) {
	// A lot of work...
}
