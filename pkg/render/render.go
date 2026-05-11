package render

import (
	"fmt"
	"reflect"
)

// Render returns a human-readable string representation of value. It accepts
// raw values or reflect.Values, and delegates to custom Renderers registered
// in opts.
func Render(value interface{}, opts *Opts) string {
	if opts == nil {
		opts = NewOpts(nil)
	}
	if opts.TB != nil {
		opts.TB.Helper()
	}

	if value, ok := value.(reflect.Value); ok {
		return render(value, opts)
	}

	return render(reflect.ValueOf(value), opts)
}

func RenderType(t reflect.Type, opts *Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}

	switch t.Kind() {
	case reflect.Invalid:
		panic("should not have reached here; bug in framework")
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return t.String()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return t.String()
	case reflect.Float32, reflect.Float64:
		return t.String()
	case reflect.Complex64, reflect.Complex128:
		return t.String()
	case reflect.Pointer:
		if !opts.ShowPointers {
			return RenderType(t.Elem(), opts)
		}
		return opts.Printf("{grey}*{reset}%s", RenderType(t.Elem(), opts))
	case reflect.Slice:
		return t.String()
	case reflect.Array:
		return t.String()
	case reflect.Struct:
		return t.String()
	case reflect.Map:
		return t.String()
	case reflect.Chan:
		opts.Fatal("channels are unsupported")
	case reflect.Func:
		opts.Fatal("functions are unsupported")
	case reflect.Interface:
		return t.String()
	default:
		panic("should not have reached here; bug in framework")
	}
	panic("should not have reached here; bug in framework")
}

func renderInternal(t reflect.Type, value string, opts *Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	if !opts.ShowTypes {
		return value
	}
	return opts.Printf("{grey}%s({reset}%s{grey}){reset}", t.String(), value)
}

func renderNil(t reflect.Type, opts *Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	if !opts.ShowTypes {
		return "nil"
	}
	return opts.Printf("{grey}nil({reset}%s{grey}){reset}", RenderType(t, opts))
}

func render(v reflect.Value, opts *Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}

	if rendered, ok := opts.Renderer(v); ok {
		return rendered
	}
	if v.CanInterface() {
		if renderable, ok := v.Interface().(Renderable); ok {
			return renderable.Render(opts)
		}
	}

	switch v.Kind() {
	case reflect.Invalid:
		panic("should not have reached here; bug in framework")
	case reflect.String:
		return renderInternal(v.Type(), fmt.Sprintf("%q", v.String()), opts)
	case reflect.Bool:
		return renderInternal(v.Type(), fmt.Sprintf("%t", v.Bool()), opts)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return renderInternal(v.Type(), fmt.Sprintf("%d", v.Int()), opts)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return renderInternal(v.Type(), fmt.Sprintf("%d", v.Uint()), opts)
	case reflect.Float32, reflect.Float64:
		return renderInternal(v.Type(), fmt.Sprintf("%g", v.Float()), opts)
	case reflect.Complex64, reflect.Complex128:
		return renderInternal(v.Type(), fmt.Sprintf("%v", v.Complex()), opts)
	case reflect.Slice, reflect.Array:
		return renderSlice(v, opts)
	case reflect.Struct:
		return renderStruct(v, opts)
	case reflect.Map:
		return renderMap(v, opts)
	case reflect.Chan:
		opts.Fatal("channels are unsupported")
	case reflect.Func:
		opts.Fatal("functions are unsupported")
	case reflect.Pointer:
		if v.IsNil() {
			return renderNil(v.Type().Elem(), opts)
		}
		if !opts.ShowPointers {
			return render(v.Elem(), opts)
		}
		return opts.Printf("{grey}&{reset}%s", render(v.Elem(), opts))
	case reflect.Interface:
		if v.IsNil() {
			return renderNil(v.Type(), opts)
		}
		return renderInternal(v.Type(), render(v.Elem(), opts), opts)
	default:
		panic("should not have reached here; bug in framework")
	}
	panic("should not have reached here; bug in framework")
}
