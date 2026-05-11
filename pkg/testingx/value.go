package testingx

import (
	"reflect"
	"testing"
)

type Value[Want any] struct {
	t       testing.TB
	current reflect.Value
	next    *Value[Want]
	fatal   bool
}

func (value *Value[Want]) Capture() Want {
	value.t.Helper()

	if !value.current.IsValid() {
		value.t.Fatalf("no value to check: more assertions chained than values returned")
		var zero Want
		return zero
	}

	want, ok := value.current.Interface().(Want)
	if !ok {
		value.t.Fatalf("expected type %T, got %T", *new(Want), value.current.Interface())
	}
	return want
}

func (value *Value[Want]) Fatal() *Value[Want] {
	value.t.Helper()
	value.fatal = true
	if value.next != nil {
		value.next.Fatal()
	}
	return value
}

func (value *Value[Want]) NonFatal() *Value[Want] {
	value.t.Helper()
	value.fatal = false
	if value.next != nil {
		value.next.NonFatal()
	}
	return value
}

func (value *Value[Want]) Fail(format string, args ...any) {
	value.t.Helper()
	if value.fatal {
		value.t.Fatalf(format, args...)
		return
	}
	value.t.Errorf(format, args...)
}
