package replica

import (
	cmd "github.com/go-epaxos/epaxos/command"
)

const (
	proposeType = iota
	preAcceptType
	preAcceptOKType
	preAcceptReplyType
	acceptType
	acceptReplyType
)

const (
	// Ballot has a format like:
	// Epoch   | Ballot  | ReplicaId
	// 20 bits | 36 bits | 8 bits
	ballotEpochWidth     uint64 = 20
	ballotBallotWidth    uint64 = 36
	ballotReplicaIdWidth uint64 = 8
	ballotEpochMask      uint64 = (1<<ballotEpochWidth - 1) << (ballotBallotWidth + ballotReplicaIdWidth)
	ballotBallotMask     uint64 = (^(1<<ballotReplicaIdWidth - 1)) & (1<<(ballotBallotWidth+ballotReplicaIdWidth) - 1)
	ballotReplicaIdMask  uint64 = 1<<ballotReplicaIdWidth - 1
)

type Propose struct {
	cmds []cmd.Command
}

type PreAccept struct {
	cmds  []cmd.Command
	deps  []InstanceIdType
	repId int
	insId InstanceIdType
}

type PreAcceptOK struct {
	insId InstanceIdType
}

type PreAcceptReply struct {
	deps  []InstanceIdType
	repId int
	insId InstanceIdType
}

// a bookkeeping for infos like maxBallot, # of nack, # of ok, etc
type InstanceInfo struct {
	acceptNackCnt int
	acceptOkCnt   int
}

func (*PreAccept) getType() uint8 {
	return preAcceptType
}
func (*PreAcceptOK) getType() uint8 {
	return preAcceptOKType
}
func (*PreAcceptReply) getType() uint8 {
	return preAcceptReplyType
}

func makeLargerBallot(b uint64) uint64 {
	remain := b & ballotEpochMask & ballotReplicaIdMask
	ballot := b & ballotBallotMask
	ballot = (ballot + (1 << ballotReplicaIdWidth)) & ballotBallotMask
	return remain | ballot
}

func getEpoch(b uint64) uint64 {
	return b & ballotEpochMask
}

func makeBallot(epoch, replicaId uint64) uint64 {
	return (epoch << (64 - ballotEpochWidth - ballotBallotWidth)) | replicaId
}
