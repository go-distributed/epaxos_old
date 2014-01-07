package replica

const (
        NONE int8 = iota
        PREACCEPTED
        ACCEPTED
        COMMITTED
        EXECUTED
)

type Instance struct {
	Cmds   interface{}
	Seq    int32
	Deps   []int32
	Status int8
}
