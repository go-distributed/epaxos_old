package replica

import (
	"testing"

	cmd "github.com/go-epaxos/epaxos/command"
	//"github.com/go-epaxos/epaxos/logger"
)

// Test if the Accept messages can be sent correctly
func TestSendAccept(t *testing.T) {
	g := makeReplicaGroup(5)
	r, messageChan, propose := setup(g)

	// done setup, now send Accepts
	r.sendAccept(r.Id, 2, messageChan)

	// test if the Accept messages are correct
	for i := 0; i < r.N/2; i++ {
		a := (<-messageChan).(*Accept)
		if a.cmds[0].Compare(propose.cmds[0]) != 0 ||
			a.cmds[1].Compare(propose.cmds[1]) != 0 {
			t.Fatal("command isn't equal")
		}
		if a.insId != 2 {
			t.Fatal("instance id should be 2")
		}
		if a.deps[0] != 1 {
			t.Fatal("deps[0] should be 1")
		}
	}
}

// Test if we can accpet the Accept messages correctly
func TestRecvAcceptOk(t *testing.T) {
	g := makeReplicaGroup(5)
	r, messageChan, _ := setup(g)

	// done setup, now send Accepts
	r.sendAccept(r.Id, 2, messageChan)
	for i := 0; i < r.N/2; i++ {
		ac := (<-messageChan).(*Accept)
		g[i+1].recvAccept(ac, messageChan)
	}

	// test if the Accepts's reply is ok (they should be)
	for i := 0; i < r.N/2; i++ {
		ar := (<-messageChan).(*AcceptReply)
		if !ar.ok {
			t.Fatal("should be ok")
		}
	}
	testNoMessageLeft(messageChan, t)
}

// Test if we reject the Accepts correctly
func TestRecvAcceptNackBallot(t *testing.T) {
	g := makeReplicaGroup(5)
	r, messageChan, _ := setup(g)

	// done setup, let's send Accepts
	r.sendAccept(r.Id, 2, messageChan)
	for i := 0; i < r.N/2; i++ {
		ac := (<-messageChan).(*Accept)
		g[i+1].InstanceMatrix[r.Id][2] = &Instance{
			// make local ballot larger, so the replica will reject Accepts
			ballot: makeLargerBallot(r.InstanceMatrix[r.Id][1].ballot),
		}
		g[i+1].recvAccept(ac, messageChan)
	}

	for i := 0; i < r.N/2; i++ {
		ar := (<-messageChan).(*AcceptReply)
		if ar.ok {
			t.Fatal("should not be ok")
		}
	}
	testNoMessageLeft(messageChan, t)
}

// Test if we reject the Accepts correctly
func TestRecvAcceptNackStatus(t *testing.T) {
	g := makeReplicaGroup(5)
	r, messageChan, _ := setup(g)

	// done setup, let send Accepts
	r.sendAccept(r.Id, 2, messageChan)

	// modify some replica's status to make them reject Accepts
	g[1].InstanceMatrix[r.Id][2] = &Instance{
		status: accepted,
	}
	g[2].InstanceMatrix[r.Id][2] = &Instance{
		status: committed,
	}

	for i := 0; i < r.N/2; i++ {
		ac := (<-messageChan).(*Accept)
		g[i+1].recvAccept(ac, messageChan)
	}

	for i := 0; i < r.N/2; i++ {
		ar := (<-messageChan).(*AcceptReply)
		if ar.ok {
			t.Fatal("should not be ok")
		}
	}
	testNoMessageLeft(messageChan, t)
}

// Test if the Commit messages are sent successfully
func TestSendCommit(t *testing.T) {
	g := makeReplicaGroup(5)
	r, messageChan, propose := setup(g)

	// done setup, let's send Commits
	r.sendCommit(r.Id, 2, messageChan)

	// test the Commit messages
	for i := 1; i < r.N; i++ {
		m := (<-messageChan).(*Commit)
		if m.cmds[0].Compare(propose.cmds[0]) != 0 ||
			m.cmds[1].Compare(propose.cmds[1]) != 0 {
			t.Fatal("command isn't equal")
		}
		if m.insId != 2 {
			t.Fatal("instance id should be 2")
		}
		if m.deps[0] != 1 {
			t.Fatal("deps[0] should be 1")
		}
	}
	testNoMessageLeft(messageChan, t)
}

// Receive the Commit messages and accept them
func TestRecvCommitOk(t *testing.T) {
	g := makeReplicaGroup(5)
	r, messageChan, propose := setup(g)

	// done setup, let's send Commits
	r.sendCommit(r.Id, 2, messageChan)

	for i := 1; i < r.N; i++ {
		m := (<-messageChan).(*Commit)
		g[i].recvCommit(m)
	}

	// test if the Commmit messages are correct
	for i := 1; i < r.N; i++ {
		inst := g[i].InstanceMatrix[r.Id][2]
		if inst.cmds[0].Compare(propose.cmds[0]) != 0 ||
			inst.cmds[1].Compare(propose.cmds[1]) != 0 {
			t.Fatal("command isn't equal")
		}
		if inst.deps[0] != 1 {
			t.Fatal("deps[0] should be 1")
		}
		if inst.status != committed {
			t.Fatal("status is not committed")
		}
	}
	testNoMessageLeft(messageChan, t)
}

// Receive the Commits, but ignore them
func TestRecvCommitIgnore(t *testing.T) {
	g := makeReplicaGroup(10)
	r, messageChan, propose := setup(g)

	// done setup, let's send Commits
	r.sendCommit(r.Id, 2, messageChan)

	// modify some replica's instances, make them ignore the coming Commits
	g[1].InstanceMatrix[r.Id][2] = &Instance{
		cmds: []cmd.Command{
			cmd.Command("paxos"),
		},
		status: committed,
	}
	g[2].InstanceMatrix[r.Id][2] = &Instance{
		cmds: []cmd.Command{
			cmd.Command("paxos"),
		},
		status: executed,
	}
	// g[3] should still accept the Commit, although it has a larger ballot
	g[3].InstanceMatrix[r.Id][2] = &Instance{
		cmds: []cmd.Command{
			cmd.Command("paxos"),
		},
		ballot: makeLargerBallot(r.InstanceMatrix[r.Id][2].ballot),
	}

	// recv Commits
	for i := 1; i < r.N; i++ {
		m := (<-messageChan).(*Commit)
		g[i].recvCommit(m)
	}

	// test if some replicas really ignore the Commits
	for i := 1; i < 3; i++ {
		if g[i].InstanceMatrix[r.Id][2].cmds[0].Compare(cmd.Command("paxos")) != 0 {
			t.Fatal("command not correct, should be 'paxos'")
		}
	}
	// test if other replicas accept the Commits
	for i := 3; i < r.N; i++ {
		inst := g[i].InstanceMatrix[r.Id][2]
		if inst.cmds[0].Compare(propose.cmds[0]) != 0 ||
			inst.cmds[1].Compare(propose.cmds[1]) != 0 {
			t.Fatal("command isn't equal")
		}
		if inst.deps[0] != 1 {
			t.Fatal("deps[0] should be 1")
		}
		if inst.status != committed {
			t.Fatal("status is not committed")
		}
	}
	testNoMessageLeft(messageChan, t)
}

// Test send Accept and then Commit messages
func TestAcceptAndCommit(t *testing.T) {
	g := makeReplicaGroup(5)
	r, messageChan, _ := setup(g)

	// done setup, let's send Accepts
	r.sendAccept(r.Id, 2, messageChan)

	for i := 0; i < r.N/2; i++ {
		ac := (<-messageChan).(*Accept)
		g[i+1].recvAccept(ac, messageChan)
	}

	for i := 0; i < r.N/2; i++ {
		ar := (<-messageChan).(*AcceptReply)
		r.recvAcceptReply(ar, messageChan)
	}

	// now r should have received enough AcceptReplies, and send out the Commits
	for i := 1; i < r.N; i++ {
		m := (<-messageChan).(*Commit)
		g[i].recvCommit(m)
	}
	testNoMessageLeft(messageChan, t)
}

// Test send Accept but no Commit messages
func TestAcceptAndAbortCommit(t *testing.T) {
	g := makeReplicaGroup(5)
	r, messageChan, _ := setup(g)

	// done setup, let's send Accepts
	r.sendAccept(r.Id, 2, messageChan)

	for i := 0; i < r.N/2; i++ {
		ac := (<-messageChan).(*Accept)
		g[i+1].InstanceMatrix[r.Id][2] = &Instance{
			// make local ballot larger, so the replica will reject Accepts
			ballot: makeLargerBallot(r.InstanceMatrix[r.Id][1].ballot),
		}
		g[i+1].recvAccept(ac, messageChan)
	}

	for i := 0; i < r.N/2; i++ {
		ar := (<-messageChan).(*AcceptReply)
		r.recvAcceptReply(ar, messageChan)
	}

	// now r should not send out any Commits
	testNoMessageLeft(messageChan, t)
}

// helpers
func makeReplicaGroup(size int) []*Replica {
	g := make([]*Replica, size)
	for i := range g {
		g[i] = startNewReplica(i, size)
	}

	return g
}

func testNoMessageLeft(messageChan chan Message, t *testing.T) {
	select {
	case <-messageChan:
		t.Fatal("should be no messages left")
	default:
		return
	}
}

func setup(g []*Replica) (*Replica, chan Message, *Propose) {
	//log = logger.ScreenErrLogger
	messageChan := make(chan Message, 100)
	propose := &Propose{
		cmds: []cmd.Command{
			cmd.Command("hello"),
			cmd.Command("world"),
		},
	}
	r := g[0]
	r.recvPropose(propose, messageChan)
	r.recvPropose(propose, messageChan) // to make the second one conflict with first one

	for i := 0; i < r.fastQuorumSize(); i++ {
		<-messageChan
		<-messageChan // clean out the PreAccepts
	}
	return r, messageChan, propose
}
