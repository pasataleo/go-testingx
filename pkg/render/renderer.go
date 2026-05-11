package render

import "reflect"

// Renderer renders values of type T as human-readable strings. Register
// custom Renderers using WithRenderer.
type Renderer[T any] interface {
	Render(value T, opts *Opts) string
}

// RendererFunc is a function adapter that implements Renderer[T].
type RendererFunc[T any] func(value T, opts *Opts) string

func (f RendererFunc[T]) Render(value T, opts *Opts) string {
	return f(value, opts)
}

// renderer is the internal, type-erased version of Renderer that operates on
// reflect.Values. See rendererAdapter for the bridge between the two.
type renderer interface {
	Render(value reflect.Value, opts *Opts) string
}

// rendererAdapter wraps a typed Renderer[T] to implement the internal renderer
// interface, converting reflect.Values back to their concrete type T.
type rendererAdapter[T any] struct {
	inner Renderer[T]
}

func (r *rendererAdapter[T]) Render(value reflect.Value, opts *Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	return r.inner.Render(value.Interface().(T), opts)
}
