package diff

import (
	"reflect"

	"github.com/pasataleo/go-testingx/pkg/render"
)

var (
	_ Result = (*missing)(nil)
)

func MissingResult(want reflect.Value) Result {
	return &missing{
		want: want,
	}
}

type missing struct {
	want reflect.Value
}

func (m *missing) Status() Status {
	return StatusMissing
}

func (m *missing) RenderGot(opts *render.Opts) string {
	panic("should never be called; bug in framework")
}

func (m *missing) RenderWant(opts *render.Opts) string {
	return render.Render(m.want, opts)
}

func (m *missing) RenderDiff(opts *Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	prefix := opts.RenderOpts().Printf("{red}-{reset} ")
	rendered := render.Render(m.want, opts.RenderOpts())
	return prefix + prefixLines(rendered, prefix)
}
