package testingx

import (
	"reflect"

	"github.com/pasataleo/go-testingx/pkg/diff"
	"github.com/pasataleo/go-testingx/pkg/render"
)

func (value *Value[Want]) Equal(want any, opts ...diff.OptsFn) *Value[Want] {
	value.t.Helper()

	if !value.current.IsValid() {
		value.t.Fatalf("no value to check: more assertions chained than values returned")
		return value.next
	}

	wantValue := reflect.ValueOf(want)

	if !wantValue.IsValid() {
		if value.isNillable() {
			return value.Nil()
		}
		o := diff.NewOpts(value.t, opts...)
		value.Fail("expected values to be equal\n  got:  %s\n  want: <nil>", render.Render(value.current, o.RenderOpts()))
		return value.next
	}

	if equal, ok := tryEqual(value.current, wantValue); ok {
		if !equal {
			o := diff.NewOpts(value.t, opts...)
			value.Fail("expected values to be equal\n  got:  %s\n  want: %s", render.Render(value.current, o.RenderOpts()), render.Render(want, o.RenderOpts()))
		}
		return value.next
	}

	o := diff.NewOpts(value.t, opts...)
	result := diff.Of(value.current.Interface(), want, o)
	if result.Status() != diff.StatusUnchanged {
		value.Fail("expected values to be equal\n%s", result.RenderDiff(o))
	}
	return value.next
}

func (value *Value[Want]) NotEqual(want any, opts ...diff.OptsFn) *Value[Want] {
	value.t.Helper()

	if !value.current.IsValid() {
		value.t.Fatalf("no value to check: more assertions chained than values returned")
		return value.next
	}

	wantValue := reflect.ValueOf(want)

	if !wantValue.IsValid() {
		if value.isNillable() {
			return value.NotNil()
		}
		return value.next
	}

	o := diff.NewOpts(value.t, opts...)
	if equal, ok := tryEqual(value.current, wantValue); ok {
		if equal {
			value.Fail("expected values to not be equal\n  got: %v", render.Render(value.current, o.RenderOpts()))
		}
		return value.next
	}

	result := diff.Of(value.current.Interface(), want, o)
	if result.Status() == diff.StatusUnchanged {
		value.Fail("expected values to not be equal\n  got: %s", render.Render(value.current, o.RenderOpts()))
	}
	return value.next
}

func tryEqual(got, want reflect.Value) (bool, bool) {
	if !want.IsValid() {
		return false, false
	}
	equalMethod := got.MethodByName("Equal")
	if !equalMethod.IsValid() && got.CanAddr() {
		equalMethod = got.Addr().MethodByName("Equal")
	}
	if !equalMethod.IsValid() {
		return false, false
	}

	mt := equalMethod.Type()
	if mt.NumIn() != 1 || mt.NumOut() != 1 {
		return false, false
	}
	if mt.Out(0) != reflect.TypeOf(true) {
		return false, false
	}
	if !want.Type().AssignableTo(mt.In(0)) {
		return false, false
	}

	results := equalMethod.Call([]reflect.Value{want})
	return results[0].Bool(), true
}
