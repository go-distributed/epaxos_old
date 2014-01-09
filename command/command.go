package command

import (
	"errors"
)

type Command []byte

type StateMachine interface {
	Execute(c Command) (interface{}, error)
	IsConflict(c1 Command, c2 Command) bool
}

func ExecuteBatch(s StateMachine, q []Command) ([]interface{}, error) {
	result := make([]interface{}, 0)
	for i := range q {
		iResult, err := s.Execute(q[i])
		if err != nil {
			return nil, errors.New("command:error in execute batched commands[" + string(i) + "]")
		}
		result = append(result, iResult)
	}
	return result, nil
}

func HaveConflicts(s StateMachine, q1 []Command, q2 []Command) bool {
	for i := range q1 {
		for j := range q2 {
			if s.IsConflict(q1[i], q2[j]) {
				return true
			}
		}
	}
	return false
}
