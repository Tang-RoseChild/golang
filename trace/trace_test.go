package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("New can not return nil")
	}
	tracer.Trace("helloworld")
	if buf.String() != "helloworld\n" {
		t.Error("trace content wrong")
	}
}

func TestOff(t *testing.T) {
	var buf bytes.Buffer
	off := Off(&buf)
	if off == nil {
		t.Error("Off can not return nil")
	}
	off.Trace("helloworld")
	if buf.String() != "" {
		t.Error("off should output nothing")
	}

	off1 := Off(nil)
	if off1 == nil {
		t.Error("Off can not return nil")
	}
}
