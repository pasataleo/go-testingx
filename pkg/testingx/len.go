package testingx

import "reflect"

type hasLen interface {
	Len() int
}

func (value *Value[Want]) length() int {
	value.t.Helper()

	if !value.current.IsValid() {
		value.t.Fatalf("no value to check: more assertions chained than values returned")
		return 0
	}

	if l, ok := value.current.Interface().(hasLen); ok {
		return l.Len()
	}

	switch value.current.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return value.current.Len()
	default:
		value.t.Fatalf("cannot check length on value of kind %s", value.current.Kind())
		return 0
	}
}

func (value *Value[Want]) Len(want int) *Value[Want] {
	value.t.Helper()
	if got := value.length(); got != want {
		value.Fail("expected length %d, got %d", want, got)
	}
	return value.next
}

func (value *Value[Want]) Empty() *Value[Want] {
	value.t.Helper()
	if got := value.length(); got != 0 {
		value.Fail("expected empty, got length %d", got)
	}
	return value.next
}

func (value *Value[Want]) NotEmpty() *Value[Want] {
	value.t.Helper()
	if got := value.length(); got == 0 {
		value.Fail("expected non-empty, got length 0")
	}
	return value.next
}
