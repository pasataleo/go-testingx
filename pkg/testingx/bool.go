package testingx

import "reflect"

func (value *Value[Want]) True() *Value[Want] {
	value.t.Helper()

	if !value.current.IsValid() {
		value.t.Fatalf("no value to check: more assertions chained than values returned")
		return value.next
	}

	if value.current.Kind() != reflect.Bool {
		value.t.Fatalf("expected bool, got %s", value.current.Type())
		return value.next
	}

	if !value.current.Bool() {
		value.Fail("expected true, got false")
	}
	return value.next
}

func (value *Value[Want]) False() *Value[Want] {
	value.t.Helper()

	if !value.current.IsValid() {
		value.t.Fatalf("no value to check: more assertions chained than values returned")
		return value.next
	}

	if value.current.Kind() != reflect.Bool {
		value.t.Fatalf("expected bool, got %s", value.current.Type())
		return value.next
	}

	if value.current.Bool() {
		value.Fail("expected false, got true")
	}
	return value.next
}
