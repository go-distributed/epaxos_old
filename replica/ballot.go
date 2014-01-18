package replica

const (
	// Ballot has a format like:
	// Epoch   | Ballot  | ReplicaId
	// 20 bits | 36 bits | 8 bits
	ballotEpochWidth     uint64 = 20
	ballotNumberWidth    uint64 = 36
	ballotReplicaIdWidth uint64 = 8
	ballotEpochMask      uint64 = ((1 << ballotEpochWidth) - 1) << (ballotNumberWidth + ballotReplicaIdWidth)
	ballotNumberMask     uint64 = (^((1 << ballotReplicaIdWidth) - 1)) & ((1 << (ballotNumberWidth + ballotReplicaIdWidth)) - 1)
	ballotReplicaIdMask  uint64 = (1 << ballotReplicaIdWidth) - 1

	ballotNumberUnit uint64 = (uint64(1) << ballotReplicaIdWidth)
)

func makeLargerBallot(b uint64) uint64 {
	return (b & ballotEpochMask) |
		(((b & ballotNumberMask) + ballotNumberUnit) & ballotNumberMask) |
		(b & ballotReplicaIdMask)
}

func getEpoch(b uint64) uint64 {
	return b & ballotEpochMask
}

func makeBallot(epoch, replicaId uint64) uint64 {
	return (epoch << (64 - ballotEpochWidth - ballotNumberWidth)) | replicaId
}

func getBallotNumber(b uint64) uint64 {
	return ((b & ballotNumberMask) >> ballotReplicaIdWidth)
}

func isInitialBallot(b uint64) bool {
	return (b & ballotNumberMask) == 0
}
