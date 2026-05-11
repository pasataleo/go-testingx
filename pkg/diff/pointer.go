package diff

import (
	"reflect"

	"github.com/pasataleo/go-testingx/pkg/render"
)

var (
	_ Result = (*pointerResult)(nil)
)

func PointerResult(got, want reflect.Value, inner Result) Result {
	return &pointerResult{got: got, want: want, inner: inner}
}

type pointerResult struct {
	got   reflect.Value
	want  reflect.Value
	inner Result
}

func (p *pointerResult) Status() Status {
	return p.inner.Status()
}

func (p *pointerResult) RenderGot(opts *render.Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	return render.Render(p.got, opts)
}

func (p *pointerResult) RenderWant(opts *render.Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	return render.Render(p.want, opts)
}

func (p *pointerResult) RenderDiff(opts *Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	if !opts.RenderOpts().ShowPointers {
		return p.inner.RenderDiff(opts)
	}

	switch p.inner.Status() {
	case StatusUnchanged:
		return "  " + opts.RenderOpts().Printf("{grey}&{reset}") + p.inner.RenderGot(opts.RenderOpts())
	case StatusChanged:
		if _, ok := p.inner.(CompositeResult); !ok {
			minusPrefix := opts.RenderOpts().Printf("{red}-{reset} ")
			plusPrefix := opts.RenderOpts().Printf("{green}+{reset} ")
			wantRendered := opts.RenderOpts().Printf("{grey}&{reset}") + p.inner.RenderWant(opts.RenderOpts())
			gotRendered := opts.RenderOpts().Printf("{grey}&{reset}") + p.inner.RenderGot(opts.RenderOpts())
			return minusPrefix + prefixLines(wantRendered, minusPrefix) + "\n" +
				plusPrefix + prefixLines(gotRendered, plusPrefix)
		}
		return opts.RenderOpts().Printf("{grey}~{reset} {grey}&{reset}%s", p.inner.RenderDiff(opts.AsNested()))
	case StatusMissing:
		prefix := opts.RenderOpts().Printf("{red}-{reset} ")
		rendered := p.RenderWant(opts.RenderOpts())
		return prefix + prefixLines(rendered, prefix)
	case StatusExtra:
		prefix := opts.RenderOpts().Printf("{green}+{reset} ")
		rendered := p.RenderGot(opts.RenderOpts())
		return prefix + prefixLines(rendered, prefix)
	default:
		panic("should not have reached here; bug in framework")
	}
}
