package dummySM

import (
	"bytes"
	"errors"

	"github.com/go-epaxos/epaxos/command"
)

type DummySM bool

func (d *DummySM) Execute(c []command.Command) ([]interface{}, error) {
	result := make([]interface{}, 0)
	for i := range c {
		if bytes.Compare(c[i], command.Command("error")) == 0 {
			return nil, errors.New("error")
		}
		result = append(result, "ok")
	}
	return result, nil
}

func (d *DummySM) HaveConflicts(c1 []command.Command, c2 []command.Command) bool {
	for i := range c1 {
		for j := range c2 {
			if bytes.Compare(c1[i], c2[j]) == 0 {
				return true
			}
		}
	}
	return false
}
