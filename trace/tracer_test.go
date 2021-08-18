package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Return function New should not be nil")
	} else {
		tracer.Trace("Testing trace package")
		if buf.String() != "Testing trace package\n" {
			t.Errorf("Trace should not write '%s'", buf.String())
		}
	}
}
