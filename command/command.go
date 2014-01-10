package command

type Command []byte

type StateMachine interface {
	Execute(c []Command) ([]interface{}, error)
	HaveConflicts(c1 []Command, c2 []Command) bool
}
