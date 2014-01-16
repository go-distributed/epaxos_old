package replica

type Prepare struct {
	ballot uint64
	repId  int
	insId  InstanceIdType
}

type PrepareReply struct {
	repId  int
	insId  InstanceIdType
}

func (*Prepare) getType() uint8 {
	return prepareType
}

func (*PrepareReply) getType() uint8 {
	return prepareReplyType
}
