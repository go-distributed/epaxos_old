package replica

import (
	cmd "github.com/go-epaxos/epaxos/command"
)

const (
	none int8 = iota
	preaccepted
	accepted
	committed
	executed
)

const (
	// Ballot has a format like:
	// Epoch   | Ballot  | ReplicaId
	// 20 bits | 36 bits | 8 bits
	ballotEpochWidth     uint64 = 20
	ballotBallotWidth    uint64 = 36
	ballotReplicaIdWidth uint64 = 8
	ballotEpochMask      uint64 = ((1 << ballotEpochWidth) - 1) << (ballotBallotWidth + ballotReplicaIdWidth)
	ballotBallotMask     uint64 = (^((1 << ballotReplicaIdWidth) - 1)) & ((1 << (ballotBallotWidth + ballotReplicaIdWidth)) - 1)
	ballotReplicaIdMask  uint64 = (1 << ballotReplicaIdWidth) - 1
)

// a bookkeeping for infos like maxBallot, # of nack, # of ok, etc
type InstanceInfo struct {
	preaccCnt     int
	preaccOkCnt   int

	acceptNackCnt int
	acceptOkCnt   int
}

type Instance struct {
	cmds   []cmd.Command
	deps   []InstanceIdType
	status int8
	ballot uint64
	info   *InstanceInfo
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
