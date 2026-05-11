package testingx

import (
	"reflect"
	"testing"
)

var testingTBType = reflect.TypeOf((*testing.TB)(nil)).Elem()

func (value *Value[Want]) Validate(fn any) *Value[Want] {
	value.t.Helper()

	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()

	if fnType.Kind() != reflect.Func {
		value.t.Fatalf("expected function, got %s", fnType)
		return value.next
	}

	if fnType.NumIn() != 2 {
		value.t.Fatalf("expected function with 2 inputs, got %d", fnType.NumIn())
		return value.next
	}

	if fnType.NumOut() != 0 {
		value.t.Fatalf("expected function with 0 outputs, got %d", fnType.NumOut())
		return value.next
	}

	if !testingTBType.AssignableTo(fnType.In(0)) {
		value.t.Fatalf("expected first argument to be testing.TB, got %s", fnType.In(0))
		return value.next
	}

	if !value.current.Type().AssignableTo(fnType.In(1)) {
		value.t.Fatalf("expected second argument to accept %s, got %s", value.current.Type(), fnType.In(1))
		return value.next
	}

	fnValue.Call([]reflect.Value{reflect.ValueOf(value.t), value.current})
	return value.next
}
