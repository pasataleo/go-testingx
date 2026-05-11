package diff

import (
	"reflect"

	"github.com/pasataleo/go-testingx/pkg/render"
)

var (
	_ Result = (*value)(nil)
)

func ValueResult(got, want reflect.Value) Result {
	return &value{
		got:  got,
		want: want,
	}
}

type value struct {
	got  reflect.Value
	want reflect.Value
}

func (v *value) Status() Status {
	if reflect.DeepEqual(v.got.Interface(), v.want.Interface()) {
		return StatusUnchanged
	}
	return StatusChanged
}

func (v *value) RenderGot(opts *render.Opts) string {
	return render.Render(v.got, opts)
}

func (v *value) RenderWant(opts *render.Opts) string {
	return render.Render(v.want, opts)
}

func (v *value) RenderDiff(opts *Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	if v.Status() == StatusUnchanged {
		return "  " + render.Render(v.got, opts.RenderOpts())
	}
	minusPrefix := opts.RenderOpts().Printf("{red}-{reset} ")
	plusPrefix := opts.RenderOpts().Printf("{green}+{reset} ")
	wantRendered := render.Render(v.want, opts.RenderOpts())
	gotRendered := render.Render(v.got, opts.RenderOpts())
	return minusPrefix + prefixLines(wantRendered, minusPrefix) + "\n" +
		plusPrefix + prefixLines(gotRendered, plusPrefix)
}
