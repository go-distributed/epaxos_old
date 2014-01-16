package replica

import (
	"fmt"
	"testing"
)

var _ = fmt.Printf

func TestSendPrepare(t *testing.T) {
	// send prepare
	// verify the prepare message for:
	// - incremented ballot number
	// - L.i instance
}

func TestRecvPrepare(t *testing.T) {
	// make prepare message to test for:
	// - ballot is larger than most recent one
	// - ballot isn't larger
	// - there's no instance!
}

func TestRecvPrepareReply(t *testing.T) {
	// A lot of work...
}
