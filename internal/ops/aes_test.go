package ops

import (
	"bytes"
	"testing"
)

func TestPcks7Padding(t *testing.T) {
	output := Pcks7Padding([]byte("YELLOW SUBMARINE"), 20)
	if !bytes.Equal(output, []byte("YELLOW SUBMARINE\x04\x04\x04\x04")) {
		t.Fail()
	}
}
