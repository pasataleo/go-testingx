package diff

import (
	"reflect"

	"github.com/pasataleo/go-testingx/pkg/render"
)

var (
	_ Result = (*extra)(nil)
)

func ExtraResult(got reflect.Value) Result {
	return &extra{
		got: got,
	}
}

type extra struct {
	got reflect.Value
}

func (e *extra) Status() Status {
	return StatusExtra
}

func (e *extra) RenderGot(opts *render.Opts) string {
	return render.Render(e.got, opts)
}

func (e *extra) RenderWant(opts *render.Opts) string {
	panic("should never be called; bug in framework")
}

func (e *extra) RenderDiff(opts *Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	prefix := opts.RenderOpts().Printf("{green}+{reset} ")
	rendered := render.Render(e.got, opts.RenderOpts())
	return prefix + prefixLines(rendered, prefix)
}
