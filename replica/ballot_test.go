package replica

import (
	"fmt"
	"testing"
)

var _ = fmt.Printf

func TestToUint64(t *testing.T) {
	b := &Ballot{0, 0, 1}
	if b.toUint64() != 1 {
		t.Fatal("expected 1 for &Ballot{0,0,1}")
	}

	b = &Ballot{0, 1, 0}
	if b.toUint64() != (1 << ballotReplicaIdWidth) {
		t.Fatalf("expected %v for &Ballot{0,1,0}\n", (1 << ballotReplicaIdWidth))
	}

	b = &Ballot{1, 0, 0}
	if b.toUint64() != (1 << (ballotReplicaIdWidth + ballotNumberWidth)) {
		t.Fatalf("expected %v for &Ballot{1,0,0}\n", (1 << ballotReplicaIdWidth))
	}
}

func TestFromUint64(t *testing.T) {
	b := &Ballot{}

	b.fromUint64(1)
	if b.epoch != 0 && b.number != 0 && b.replicaId != 1 {
		t.Fatal("expected ballot replicaId to be 1")
	}

	b.fromUint64((1 << ballotReplicaIdWidth))
	if b.epoch != 0 && b.number != 1 && b.replicaId != 0 {
		t.Fatal("expected ballot number to be 1")
	}
	b.fromUint64((1 << (ballotReplicaIdWidth + ballotNumberWidth)))
	if b.epoch != 1 && b.number != 0 && b.replicaId != 0 {
		t.Fatal("expected ballot epoch to be 1")
	}
}
