package replica

type Message interface {
	getType() uint8
}
