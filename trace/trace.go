package trace

import (
	"fmt"
	"io"
)

// Tracer used to trace events
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

type niltracer struct{}

func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))

}

func (t *niltracer) Trace(a ...interface{}) {
	// no output
}

// New return a new tracer
func New(w io.Writer) Tracer {

	return &tracer{out: w}
}

// Off no trace anymore
func Off(w io.Writer) Tracer {
	return &niltracer{}
}
