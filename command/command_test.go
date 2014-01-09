package command

import (
	"bytes"
	"errors"
	"testing"
)

type DummySM bool

func (d *DummySM) Execute(c Command) (interface{}, error) {
	if bytes.Compare(c, Command("error")) == 0 {
		return nil, errors.New("error")
	}
	return "ok", nil
}

func (d *DummySM) IsConflict(c1 Command, c2 Command) bool {
	if bytes.Compare(c1, c2) == 0 {
		return true
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

	r, err := ExecuteBatch(dm, q1)
	if err != nil {
		t.Fatal("should not have error: ", err)
	}

	for i := range r {
		if r[i] != "ok" {
			t.Fatal("result should be 'ok'")
		}
	}

	r, err = ExecuteBatch(dm, q2)
	if err == nil {
		t.Fatal("should have error")
	}

	if !HaveConflicts(dm, q1, q2) {
		t.Fatal("should have conflict")
	}

	if HaveConflicts(dm, q1, q3) {
		t.Fatal("should not have conflict")
	}
}
