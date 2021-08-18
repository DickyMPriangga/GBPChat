package trace

import "io"

type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {

}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}
