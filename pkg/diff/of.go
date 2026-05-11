package diff

import "reflect"

func isNilable(k reflect.Kind) bool {
	switch k {
	case reflect.Map, reflect.Slice, reflect.Pointer, reflect.Interface, reflect.Chan, reflect.Func:
		return true
	default:
		return false
	}
}

// Of computes the diff between got and want. It accepts raw values or
// reflect.Values, and delegates to custom Differs registered in opts.
func Of(got, want interface{}, opts *Opts) Result {
	if opts == nil {
		opts = NewOpts(nil)
	}
	if opts.TB != nil {
		opts.TB.Helper()
	}

	var g, w reflect.Value
	if value, ok := got.(reflect.Value); ok {
		g = value
	} else {
		g = reflect.ValueOf(got)
	}
	if value, ok := want.(reflect.Value); ok {
		w = value
	} else {
		w = reflect.ValueOf(want)
	}

	return of(g, w, opts)
}

func of(got, want reflect.Value, opts *Opts) Result {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	if !got.IsValid() && !want.IsValid() {
		panic("can not diff two invalid values")
	}
	if !got.IsValid() {
		return MissingResult(want)
	}
	if !want.IsValid() {
		return ExtraResult(got)
	}

	if got.Kind() != want.Kind() {
		return ValueResult(got, want)
	}

	if isNilable(got.Kind()) && (got.IsNil() || want.IsNil()) {
		return ValueResult(got, want)
	}

	// Check for registered Differ (exact type, then interface match).
	if result, ok := opts.Differ(got, want); ok {
		return result
	}

	// Check for Diffable interface (got knows how to diff itself).
	if got.CanInterface() && got.Type() == want.Type() {
		if diffable, ok := tryDiffable(got, want, opts); ok {
			return diffable
		}
	}

	switch got.Kind() {
	case reflect.Map:
		return ofMap(got, want, opts)
	case reflect.Struct:
		return ofStruct(got, want, opts)
	case reflect.Slice, reflect.Array:
		return ofSlice(got, want, opts)
	case reflect.Pointer:
		return PointerResult(got, want, of(got.Elem(), want.Elem(), opts))
	case reflect.Interface:
		return InterfaceResult(got, want, of(got.Elem(), want.Elem(), opts))
	case reflect.Chan:
		opts.Fatal("channels are unsupported")
	case reflect.Func:
		opts.Fatal("functions are unsupported")
	default:
		return ValueResult(got, want)
	}
	panic("should not have reached here; bug in framework")
}
