package replica

import (
	"fmt"
	"reflect"
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

// Receivers have no info about the Instance in Prepare
// Test if they will accept the Prepare message
// PrepareReply should be ok with no-op
// Success: Accept the Prepares, Fail: Reject the Prepares
func TestRecvPrepareNoInstance(t *testing.T) {
	g, r, messageChan := recoveryTestSetup(5)
	r.sendPrepare(0, conflictNotFound+1, messageChan)
	var pp *Prepare
	for i := 1; i < r.N; i++ {
		pp = (<-messageChan).(*Prepare)
		g[i].recvPrepare(pp, messageChan)
	}
	// receivers have no info about the instance, they should reply ok
	for i := 1; i < r.N; i++ {
		pr := (<-messageChan).(*PrepareReply)
		if !reflect.DeepEqual(pr, &PrepareReply{
			ok:     true,
			ballot: pp.ballot,
			status: -1,
			cmds:   nil,
			deps:   make([]InstanceIdType, r.N), // TODO: makeInitialDeps
			repId:  0,
			insId:  conflictNotFound + 1,
		}) {
			t.Fatal("PrepareReply message error")
		}
	}
	testNoMessagesLeft(messageChan, t)
}

// Receivers have larger ballot
// Test if they will reject the Prepare message
// PrepareReply should be not ok, with their version of instances
// Success: Reject the Prepares, Fail: Accept the Prepares
func TestRecvPrepareReject(t *testing.T) {
	g, r, messageChan := recoveryTestSetup(5)
	r.sendPrepare(0, conflictNotFound+1, messageChan)

	for i := 1; i < r.N; i++ {
		pp := (<-messageChan).(*Prepare)
		// create instance in receivers, and make larger ballots
		g[i].InstanceMatrix[0][conflictNotFound+1] = &Instance{
			status: accepted,
			cmds: []cmd.Command{
				cmd.Command("paxos"),
			},
			deps: []InstanceIdType{1, 0, 0, 0, 0},
			// ballot num == 1
			ballot: g[i].makeInitialBallot().getIncNumCopy(),
		}
		g[i].recvPrepare(pp, messageChan)
	}

	for i := 1; i < r.N; i++ {
		pr := (<-messageChan).(*PrepareReply)
		if !reflect.DeepEqual(pr, &PrepareReply{
			ok: false,
			// info of the receivers' instance
			status: accepted,
			cmds: []cmd.Command{
				cmd.Command("paxos"),
			},
			deps: []InstanceIdType{1, 0, 0, 0, 0},
			// ballot num == 1
			ballot: g[i].makeInitialBallot().getIncNumCopy(), // receiver's ballot
			repId:  0,
			insId:  conflictNotFound + 1,
		}) {
			t.Fatal("PrepareReply message error")
		}
	}
	testNoMessagesLeft(messageChan, t)
}

// Receivers have smaller ballot
// Test if they will accepte the Prepare message
// PrepareReply should be ok, with their version of instance
// Success: Accept the Prepares, Fail: Reject the Prepares
func TestRecvPrepareAccept(t *testing.T) {
	g, r, messageChan := recoveryTestSetup(5)
	r.sendPrepare(0, conflictNotFound+1, messageChan)

	for i := 1; i < r.N; i++ {
		pp := (<-messageChan).(*Prepare)
		// create instance in receivers, and make larger ballots
		g[i].InstanceMatrix[0][conflictNotFound+1] = &Instance{
			status: accepted,
			cmds: []cmd.Command{
				cmd.Command("paxos"),
			},
			deps:   []InstanceIdType{1, 0, 0, 0, 0},
			ballot: g[i].makeInitialBallot(),
		}
		g[i].recvPrepare(pp, messageChan)
	}

	for i := 1; i < r.N; i++ {
		pr := (<-messageChan).(*PrepareReply)
		if !reflect.DeepEqual(pr, &PrepareReply{
			ok: true,
			// info of the receivers' instance
			status: accepted,
			cmds: []cmd.Command{
				cmd.Command("paxos"),
			},
			deps:   []InstanceIdType{1, 0, 0, 0, 0},
			ballot: r.makeInitialBallot().getIncNumCopy(), // sender's ballot
			repId:  0,
			insId:  conflictNotFound + 1,
		}) {
			t.Fatal("PrepareReply message error")
		}
	}
	testNoMessagesLeft(messageChan, t)
}

func TestRecvPrepareReply(t *testing.T) {
	// A lot of work...
}

// helpers
func recoveryTestSetup(size int) ([]*Replica, *Replica, chan Message) {
	g := testMakeRepGroup(size)
	messageChan := make(chan Message, 100)
	return g, g[0], messageChan
}
