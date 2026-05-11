package testingx

import (
	"reflect"
	"testing"

	render3 "github.com/pasataleo/go-testingx/pkg/render"
)

func Panics(t testing.TB, opts *render3.Opts, fn any, args ...any) *Value[interface{}] {
	t.Helper()

	fnValue, reflectedArgs := prepareCall(t, fn, args...)
	if !fnValue.IsValid() {
		return nil
	}

	var recovered any
	func() {
		defer func() {
			recovered = recover()
		}()
		fnValue.Call(reflectedArgs)
	}()

	if recovered == nil {
		t.Fatalf("expected function to panic, but it did not")
		return nil
	}

	return &Value[interface{}]{
		t:       t,
		current: reflect.ValueOf(recovered),
		next: &Value[interface{}]{
			t:       t,
			current: reflect.Value{},
			fatal:   false,
		},
		fatal: false,
	}
}

func NotPanics(t testing.TB, opts *render3.Opts, fn any, args ...any) *Value[interface{}] {
	t.Helper()

	fnValue, reflectedArgs := prepareCall(t, fn, args...)
	if !fnValue.IsValid() {
		return nil
	}

	var results []reflect.Value
	var recovered any
	func() {
		defer func() {
			recovered = recover()
		}()
		results = fnValue.Call(reflectedArgs)
	}()

	if recovered != nil {
		if opts == nil {
			opts = render3.NewOpts(t)
		}
		t.Fatalf("expected function not to panic, but it panicked with: %s", render3.Render(recovered, opts))
		return nil
	}

	return buildValues[interface{}](t, results)
}

func PanicsAs[T any](t testing.TB, opts *render3.Opts, fn any, args ...any) *Value[T] {
	t.Helper()

	fnValue, reflectedArgs := prepareCall(t, fn, args...)
	if !fnValue.IsValid() {
		return nil
	}

	var recovered any
	func() {
		defer func() {
			recovered = recover()
		}()
		fnValue.Call(reflectedArgs)
	}()

	if recovered == nil {
		t.Fatalf("expected function to panic, but it did not")
		return nil
	}

	val, ok := recovered.(T)
	if !ok {
		if opts == nil {
			opts = render3.NewOpts(t)
		}
		t.Fatalf("expected panic value of type %T, got %T: %s", *new(T), recovered, render3.Render(recovered, opts))
		return nil
	}

	return &Value[T]{
		t:       t,
		current: reflect.ValueOf(val),
		next: &Value[T]{
			t:       t,
			current: reflect.Value{},
			fatal:   false,
		},
		fatal: false,
	}
}

func NotPanicsAs[T any](t testing.TB, opts *render3.Opts, fn any, args ...any) *Value[T] {
	t.Helper()

	fnValue, reflectedArgs := prepareCall(t, fn, args...)
	if !fnValue.IsValid() {
		return nil
	}

	var results []reflect.Value
	var recovered any
	func() {
		defer func() {
			recovered = recover()
		}()
		results = fnValue.Call(reflectedArgs)
	}()

	if recovered != nil {
		if opts == nil {
			opts = render3.NewOpts(t)
		}
		t.Fatalf("expected function not to panic, but it panicked with: %s", render3.Render(recovered, opts))
		return nil
	}

	return buildValues[T](t, results)
}
