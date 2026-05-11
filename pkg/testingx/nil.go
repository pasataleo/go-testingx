package testingx

import (
	"reflect"

	render3 "github.com/pasataleo/go-testingx/pkg/render"
)

func (value *Value[Want]) isNillable() bool {
	switch value.current.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return true
	default:
		return false
	}
}

func (value *Value[Want]) NotNil() *Value[Want] {
	value.t.Helper()

	if !value.current.IsValid() {
		value.t.Fatalf("no value to check: more assertions chained than values returned")
		return value.next
	}

	switch value.current.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		if value.current.IsNil() {
			value.Fail("expected non-nil value, got nil")
		}
	default:
		value.t.Fatalf("cannot check nil on value of kind %s", value.current.Kind())
	}

	return value.next
}

func (value *Value[Want]) Nil(opts ...render3.OptsFn) *Value[Want] {
	value.t.Helper()

	if !value.current.IsValid() {
		value.t.Fatalf("no value to check: more assertions chained than values returned")
		return value.next
	}

	switch value.current.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		if !value.current.IsNil() {
			o := render3.NewOpts(value.t, opts...)
			value.Fail("expected nil value, got %s", render3.Render(value.current, o))
		}
	default:
		value.t.Fatalf("cannot check nil on value of kind %s", value.current.Kind())
	}

	return value.next
}
