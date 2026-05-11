package diff

import (
	"reflect"

	"github.com/pasataleo/go-testingx/pkg/render"
)

var (
	_ Result = (*interfaceResult)(nil)
)

func InterfaceResult(got, want reflect.Value, inner Result) Result {
	return &interfaceResult{got: got, want: want, inner: inner}
}

type interfaceResult struct {
	got   reflect.Value
	want  reflect.Value
	inner Result
}

func (i *interfaceResult) Status() Status {
	return i.inner.Status()
}

func (i *interfaceResult) RenderGot(opts *render.Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	return render.Render(i.got, opts)
}

func (i *interfaceResult) RenderWant(opts *render.Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	return render.Render(i.want, opts)
}

func (i *interfaceResult) renderInternal(content string, opts *render.Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	if !opts.ShowTypes {
		return content
	}
	return opts.Printf("{grey}%s({reset}%s{grey}){reset}", render.RenderType(i.got.Type(), opts), content)
}

func (i *interfaceResult) RenderDiff(opts *Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	switch i.inner.Status() {
	case StatusUnchanged:
		return "  " + i.renderInternal(i.inner.RenderGot(opts.RenderOpts()), opts.RenderOpts())
	case StatusChanged:
		if _, ok := i.inner.(CompositeResult); !ok {
			minusPrefix := opts.RenderOpts().Printf("{red}-{reset} ")
			plusPrefix := opts.RenderOpts().Printf("{green}+{reset} ")
			wantRendered := i.renderInternal(i.inner.RenderWant(opts.RenderOpts()), opts.RenderOpts())
			gotRendered := i.renderInternal(i.inner.RenderGot(opts.RenderOpts()), opts.RenderOpts())
			return minusPrefix + prefixLines(wantRendered, minusPrefix) + "\n" +
				plusPrefix + prefixLines(gotRendered, plusPrefix)
		}
		return opts.RenderOpts().Printf("{grey}~{reset} %s", i.renderInternal(i.inner.RenderDiff(opts.AsNested()), opts.RenderOpts()))
	case StatusMissing:
		prefix := opts.RenderOpts().Printf("{red}-{reset} ")
		rendered := i.RenderWant(opts.RenderOpts())
		return prefix + prefixLines(rendered, prefix)
	case StatusExtra:
		prefix := opts.RenderOpts().Printf("{green}+{reset} ")
		rendered := i.RenderGot(opts.RenderOpts())
		return prefix + prefixLines(rendered, prefix)
	default:
		panic("should not have reached here; bug in framework")
	}
}
