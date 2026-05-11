package render

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pasataleo/go-colour/pkg/colour"
)

type (
	// OptsFn is a functional option for configuring render operations.
	OptsFn func(opt *Opts)

	// Opts holds configuration for render operations, including the test
	// context, colour settings, and any registered custom Renderers.
	Opts struct {
		// TB is the testing context used for error reporting. If nil,
		// failures will panic instead of calling TB.Fatal.
		TB testing.TB
		// Colour controls colour output settings.
		Colour colour.Colour

		// ShowTypes controls whether values are wrapped with type
		// information, e.g. string("hello") vs "hello".
		ShowTypes bool
		// ShowPointers controls whether pointer metadata (& prefix,
		// * in type names) is rendered.
		ShowPointers bool
		// ShowContents controls whether composite types (maps, slices,
		// arrays, structs) are expanded or collapsed to {...}.
		ShowContents bool
		// ShowIndices controls whether slice and array elements are
		// prefixed with their index, e.g. [0]: int(1).
		ShowIndices bool
		// SkipUnexported controls whether unexported struct fields are
		// silently skipped. If false, encountering an unexported field
		// will call Fatal.
		SkipUnexported bool

		renderers  map[reflect.Type]renderer
		skipFields map[string]struct{}
	}
)

// DisableColour returns an OptsFn that disables colour output.
func DisableColour() OptsFn {
	return func(opts *Opts) {
		opts.Colour.Disable = true
	}
}

// WithTypes returns an OptsFn that controls whether type information is
// included in rendered output.
func WithTypes(show bool) OptsFn {
	return func(opts *Opts) {
		opts.ShowTypes = show
	}
}

// WithPointers returns an OptsFn that controls whether pointer metadata
// (& prefix, nil annotations) is included in rendered output.
func WithPointers(show bool) OptsFn {
	return func(opts *Opts) {
		opts.ShowPointers = show
	}
}

// WithIndices returns an OptsFn that controls whether slice and array
// elements are prefixed with their index.
func WithIndices(show bool) OptsFn {
	return func(opts *Opts) {
		opts.ShowIndices = show
	}
}

// WithContents returns an OptsFn that controls whether the contents of
// composite types (maps, slices, arrays, structs) are expanded or collapsed.
func WithContents(show bool) OptsFn {
	return func(opts *Opts) {
		opts.ShowContents = show
	}
}

// WithSkipUnexported returns an OptsFn that controls whether unexported struct
// fields are silently skipped. If false (the default), encountering an
// unexported field will call Fatal.
func WithSkipUnexported(skip bool) OptsFn {
	return func(opts *Opts) {
		opts.SkipUnexported = skip
	}
}

// WithSkipField returns an OptsFn that registers a field to skip during struct
// rendering. The field argument can be just a field name (e.g. "Age") to skip
// that field in all structs, or "StructName.FieldName" (e.g. "Person.Age") to
// skip only in a specific struct type.
func WithSkipField(field string) OptsFn {
	return func(opts *Opts) {
		opts.skipFields[field] = struct{}{}
	}
}

// WithRenderer registers a custom Renderer for type T. If T is an interface
// type, it will be used for any concrete types that implement T.
func WithRenderer[T any](renderer Renderer[T]) OptsFn {
	return func(opts *Opts) {
		t := reflect.TypeOf((*T)(nil)).Elem()
		opts.renderers[t] = &rendererAdapter[T]{inner: renderer}
	}
}

// NewOpts creates a new Opts with the given testing.TB and functional options.
// If tb is nil, failures will panic instead of calling tb.Fatal.
func NewOpts(tb testing.TB, fns ...OptsFn) *Opts {
	opts := &Opts{
		TB:           tb,
		Colour:       colour.New(),
		ShowTypes:    true,
		ShowPointers: true,
		ShowContents: true,
		ShowIndices:  true,
		renderers:    make(map[reflect.Type]renderer),
		skipFields:   make(map[string]struct{}),
	}
	for _, fn := range fns {
		fn(opts)
	}
	return opts
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

// Print formats using fmt.Sprint and applies colour settings.
func (opts *Opts) Print(args ...interface{}) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	return opts.Colour.Colour(fmt.Sprint(args...))
}

// Printf formats using fmt.Sprintf and applies colour settings.
func (opts *Opts) Printf(format string, args ...interface{}) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	return opts.Colour.Colourf(format, args...)
}

// Println formats using fmt.Sprintln and applies colour settings.
func (opts *Opts) Println(args ...interface{}) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	return opts.Colour.Colour(fmt.Sprintln(args...))
}

// SkipField reports whether the given struct field should be skipped. It
// matches against both "fieldName" entries (which match in any struct) and
// "structName.fieldName" entries (which match only in the named struct).
func (opts *Opts) SkipField(structName, fieldName string) bool {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	if _, ok := opts.skipFields[fieldName]; ok {
		return true
	}
	if _, ok := opts.skipFields[structName+"."+fieldName]; ok {
		return true
	}
	return false
}

// Renderer looks up a registered custom Renderer for the type of value.
// It first checks for an exact type match, then checks for interface matches.
// Returns the rendered string and true if a custom Renderer was found, or ""
// and false otherwise.
func (opts *Opts) Renderer(value reflect.Value) (string, bool) {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	if renderer, ok := opts.renderers[value.Type()]; ok {
		return renderer.Render(value, opts), true
	}
	for t, differ := range opts.renderers {
		if t.Kind() == reflect.Interface && value.Type().Implements(t) {
			return differ.Render(value, opts), true
		}
	}
	return "", false
}
