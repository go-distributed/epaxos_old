package command

import (
	"bytes"
	"errors"
	"testing"
)

type DummySM bool

func (d *DummySM) Execute(c []Command) ([]interface{}, error) {
	result := make([]interface{}, 0)
	for i := range c {
		if bytes.Compare(c[i], Command("error")) == 0 {
			return nil, errors.New("error")
		}
		result = append(result, "ok")
	}
	return result, nil
}

func (d *DummySM) HaveConflicts(c1 []Command, c2 []Command) bool {
	for i := range c1 {
		for j := range c2 {
			if bytes.Compare(c1[i], c2[j]) == 0 {
				return true
			}
		}
	}
	return false
}

func TestCommand(t *testing.T) {
	dm := new(DummySM)
	a := Command("hello")
	b := Command("world")
	c := Command("hello")
	d := Command("error")
	e := Command("part-time")
	f := Command("parlement")

	q1 := []Command{
		a,
		b,
	}
	q2 := []Command{
		c,
		d,
	}
	q3 := []Command{
		e,
		f,
	}

	sm := StateMachine(dm)
	r, err := sm.Execute(q1)
	if err != nil {
		t.Fatal("should not have error: ", err)
	}

	for i := range r {
		if r[i] != "ok" {
			t.Fatal("result should be 'ok'")
		}
	}

	r, err = sm.Execute(q2)
	if err == nil {
		t.Fatal("should have error")
	}

	if !sm.HaveConflicts(q1, q2) {
		t.Fatal("should have conflict")
	}

	if sm.HaveConflicts(q1, q3) {
		t.Fatal("should not have conflict")
	}
}
