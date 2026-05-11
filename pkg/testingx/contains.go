package testingx

import (
	"reflect"
	"strings"

	render3 "github.com/pasataleo/go-testingx/pkg/render"
)

func (value *Value[Want]) Contains(want any, opts ...render3.OptsFn) *Value[Want] {
	value.t.Helper()

	if !value.current.IsValid() {
		value.t.Fatalf("no value to check: more assertions chained than values returned")
		return value.next
	}

	wantValue := reflect.ValueOf(want)
	o := render3.NewOpts(value.t, opts...)

	if contains, ok := tryContains(value.current, wantValue); ok {
		if !contains {
			value.Fail("expected value to contain %s", render3.Render(want, o))
		}
		return value.next
	}

	if value.current.Kind() == reflect.String {
		wantStr, ok := want.(string)
		if !ok {
			value.t.Fatalf("expected string argument for string contains, got %T", want)
			return value.next
		}
		if !strings.Contains(value.current.String(), wantStr) {
			value.Fail("expected %s to contain %s", render3.Render(value.current, o), render3.Render(wantStr, o))
		}
		return value.next
	}

	switch value.current.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.current.Len(); i++ {
			elem := value.current.Index(i)
			if reflect.DeepEqual(elem.Interface(), want) {
				return value.next
			}
		}
		value.Fail("expected collection to contain %s", render3.Render(want, o))
		return value.next
	case reflect.Map:
		mapValue := value.current.MapIndex(wantValue)
		if !mapValue.IsValid() {
			value.Fail("expected map to contain key %s", render3.Render(want, o))
		}
		return value.next
	default:
		value.t.Fatalf("cannot check contains on value of kind %s", value.current.Kind())
		return value.next
	}
}

func (value *Value[Want]) NotContains(want any, opts ...render3.OptsFn) *Value[Want] {
	value.t.Helper()

	if !value.current.IsValid() {
		value.t.Fatalf("no value to check: more assertions chained than values returned")
		return value.next
	}

	wantValue := reflect.ValueOf(want)
	o := render3.NewOpts(value.t, opts...)

	if contains, ok := tryContains(value.current, wantValue); ok {
		if contains {
			value.Fail("expected value to not contain %s", render3.Render(want, o))
		}
		return value.next
	}

	if value.current.Kind() == reflect.String {
		wantStr, ok := want.(string)
		if !ok {
			value.t.Fatalf("expected string argument for string contains, got %T", want)
			return value.next
		}
		if strings.Contains(value.current.String(), wantStr) {
			value.Fail("expected %s to not contain %s", render3.Render(value.current, o), render3.Render(wantStr, o))
		}
		return value.next
	}

	switch value.current.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.current.Len(); i++ {
			elem := value.current.Index(i)
			if reflect.DeepEqual(elem.Interface(), want) {
				value.Fail("expected collection to not contain %s", render3.Render(want, o))
				return value.next
			}
		}
		return value.next
	case reflect.Map:
		mapValue := value.current.MapIndex(wantValue)
		if mapValue.IsValid() {
			value.Fail("expected map to not contain key %s", render3.Render(want, o))
		}
		return value.next
	default:
		value.t.Fatalf("cannot check contains on value of kind %s", value.current.Kind())
		return value.next
	}
}

func tryContains(got, want reflect.Value) (bool, bool) {
	if !want.IsValid() {
		return false, false
	}
	containsMethod := got.MethodByName("Contains")
	if !containsMethod.IsValid() && got.CanAddr() {
		containsMethod = got.Addr().MethodByName("Contains")
	}
	if !containsMethod.IsValid() {
		return false, false
	}

	mt := containsMethod.Type()
	if mt.NumIn() != 1 || mt.NumOut() != 1 {
		return false, false
	}
	if mt.Out(0) != reflect.TypeOf(true) {
		return false, false
	}
	if !want.Type().AssignableTo(mt.In(0)) {
		return false, false
	}

	results := containsMethod.Call([]reflect.Value{want})
	return results[0].Bool(), true
}
