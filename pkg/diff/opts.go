package diff

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	render2 "github.com/pasataleo/go-testingx/pkg/render"
)

type (
	// OptsFn is a functional option for configuring diff operations.
	OptsFn func(opts *Opts)

	// Opts holds configuration for diff operations, including the test
	// context, colour settings, and any registered custom Differs.
	Opts struct {
		TB testing.TB

		// ShowUnchanged controls whether unchanged lines are included
		// in diff output. Defaults to false.
		ShowUnchanged bool

		differs    map[reflect.Type]differ
		renderOpts *render2.Opts
		nested     bool
	}
)

// WithRenderOpts sets the render options used when rendering diff results.
func WithRenderOpts(opts *render2.Opts) OptsFn {
	return func(o *Opts) {
		o.renderOpts = opts
	}
}

// WithShowUnchanged controls whether unchanged lines are included in diff output.
func WithShowUnchanged(show bool) OptsFn {
	return func(opts *Opts) {
		opts.ShowUnchanged = show
	}
}

// WithDiffer registers a custom Differ for type T. If T is an interface type,
// it will be used for any concrete types that implement T.
func WithDiffer[T any](differ Differ[T]) OptsFn {
	return func(opts *Opts) {
		t := reflect.TypeOf((*T)(nil)).Elem()
		opts.differs[t] = &differAdapter[T]{inner: differ}
	}
}

// NewOpts creates a new Opts with the given testing.TB and functional options.
// If tb is nil, failures will panic instead of calling tb.Fatal.
func NewOpts(tb testing.TB, fns ...OptsFn) *Opts {
	opts := &Opts{
		TB:         tb,
		differs:    make(map[reflect.Type]differ),
		renderOpts: render2.NewOpts(tb),
	}
	for _, fn := range fns {
		fn(opts)
	}
	return opts
}

// RenderOpts returns the render options for rendering diff results.
func (opts *Opts) RenderOpts() *render2.Opts {
	return opts.renderOpts
}

// Nested reports whether the opts are for a nested diff context, where the
// parent handles rendering the status prefix.
func (opts *Opts) Nested() bool {
	return opts.nested
}

// AsNested returns a shallow copy with nested set to true. This should only
// be called by composite diff types (maps, structs, slices) when rendering
// inner diff results.
func (opts *Opts) AsNested() *Opts {
	cp := *opts
	cp.nested = true
	return &cp
}

// prefixLines inserts prefix after the leading whitespace of each continuation
// line (lines after the first) in s.
func prefixLines(s, prefix string) string {
	if !strings.Contains(s, "\n") {
		return s
	}
	lines := strings.Split(s, "\n")
	for i := 1; i < len(lines); i++ {
		trimmed := strings.TrimLeft(lines[i], " ")
		indent := lines[i][:len(lines[i])-len(trimmed)]
		if len(indent) == 0 {
			continue
		}
		lines[i] = indent + prefix + trimmed
	}
	return strings.Join(lines, "\n")
}

// Fatal reports a fatal error via TB if available, otherwise panics.
func (opts *Opts) Fatal(msg string) {
	if opts.TB != nil {
		opts.TB.Helper()
		opts.TB.Fatal(msg)
	}
	panic(msg)
}

// Fatalf reports a formatted fatal error via TB if available, otherwise panics.
func (opts *Opts) Fatalf(format string, args ...any) {
	if opts.TB != nil {
		opts.TB.Helper()
		opts.TB.Fatalf(format, args...)
	}
	panic(fmt.Sprintf(format, args...))
}

// Differ looks up a registered custom Differ for the types of got and want.
// It first checks for an exact type match, then checks for interface matches.
// Returns the Result and true if a custom Differ was found, or nil and false otherwise.
func (opts *Opts) Differ(got, want reflect.Value) (Result, bool) {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	if differ, ok := opts.differs[got.Type()]; ok {
		if got.Type() == want.Type() {
			return differ.Diff(got, want, opts), true
		}
	}
	for t, differ := range opts.differs {
		if t.Kind() == reflect.Interface && got.Type().Implements(t) && want.Type().Implements(t) {
			return differ.Diff(got, want, opts), true
		}
	}
	return nil, false
}
