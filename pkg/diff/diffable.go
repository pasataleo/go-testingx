package diff

import "reflect"

// Diffable is implemented by types that know how to diff themselves against
// another value. The receiver is got and the argument is want.
type Diffable[T any] interface {
	Diff(want T, opts *Opts) Result
}

var (
	optsType   = reflect.TypeOf((*Opts)(nil))
	resultType = reflect.TypeOf((*Result)(nil)).Elem()
)

// tryDiffable checks whether got implements Diffable[T] where T is got's type,
// by looking for a Diff(T, *Opts) Result method via reflection.
func tryDiffable(got, want reflect.Value, opts *Opts) (Result, bool) {
	if opts.TB != nil {
		opts.TB.Helper()
	}

	method := got.MethodByName("Diff")
	if !method.IsValid() {
		return nil, false
	}

	mt := method.Type()
	if mt.NumIn() != 2 || mt.NumOut() != 1 {
		return nil, false
	}
	if mt.In(0) != got.Type() || mt.In(1) != optsType || mt.Out(0) != resultType {
		return nil, false
	}

	results := method.Call([]reflect.Value{want, reflect.ValueOf(opts)})
	return results[0].Interface().(Result), true
}
