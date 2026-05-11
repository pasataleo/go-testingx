package diff

import "reflect"

// Differ compares two values of type T and returns a Result describing
// the differences between them. Register custom Differs using WithDiffer.
type Differ[T any] interface {
	Diff(got, want T, opts *Opts) Result
}

// differ is the internal, type-erased version of Differ that operates on
// reflect.Values. See differAdapter for the bridge between the two.
type differ interface {
	Diff(got, want reflect.Value, opts *Opts) Result
}

// differAdapter wraps a typed Differ[T] to implement the internal differ
// interface, converting reflect.Values back to their concrete type T.
type differAdapter[T any] struct {
	inner Differ[T]
}

func (d *differAdapter[T]) Diff(got, want reflect.Value, opts *Opts) Result {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	return d.inner.Diff(got.Interface().(T), want.Interface().(T), opts)
}
