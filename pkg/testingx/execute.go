package testingx

import (
	"reflect"
	"testing"
)

func prepareCall(t testing.TB, fn any, args ...any) (reflect.Value, []reflect.Value) {
	t.Helper()

	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()

	// Ensure fn is a function
	if fnType.Kind() != reflect.Func {
		t.Fatalf("expected function, got %s", fnType)
		return reflect.Value{}, nil
	}

	isVariadic := fnType.IsVariadic()
	requiredArguments := fnType.NumIn()
	if isVariadic {
		requiredArguments = requiredArguments - 1
	}

	if !isVariadic && len(args) != requiredArguments {
		t.Fatalf("expected %d arguments, got %d", requiredArguments, len(args))
		return reflect.Value{}, nil
	}
	if isVariadic && len(args) < requiredArguments {
		t.Fatalf("expected at least %d arguments, got %d", requiredArguments, len(args))
		return reflect.Value{}, nil
	}

	var reflectedArgs []reflect.Value
	for ix := 0; ix < requiredArguments; ix++ {
		wantType := fnType.In(ix)
		argValue := reflectArg(args[ix], wantType)

		if !argValue.Type().AssignableTo(wantType) {
			t.Fatalf("argument %d: expected type %s, got %s", ix, wantType, argValue.Type())
			return reflect.Value{}, nil
		}

		reflectedArgs = append(reflectedArgs, argValue)
	}

	if isVariadic {
		wantType := fnType.In(requiredArguments).Elem()
		for ix := requiredArguments; ix < len(args); ix++ {
			argValue := reflectArg(args[ix], wantType)

			if !argValue.Type().AssignableTo(wantType) {
				t.Fatalf("argument %d: expected type %s, got %s", ix, wantType, argValue.Type())
				return reflect.Value{}, nil
			}

			reflectedArgs = append(reflectedArgs, argValue)
		}
	}

	return fnValue, reflectedArgs
}

func buildValues[T any](t testing.TB, results []reflect.Value) *Value[T] {
	current := &Value[T]{
		t:       t,
		current: reflect.Value{},
		next:    nil,
		fatal:   false,
	}

	for i := 0; i < len(results); i++ {
		previous := current
		current = &Value[T]{
			t:       t,
			current: results[i],
			next:    previous,
			fatal:   false,
		}
	}
	return current
}

func execute[T any](t testing.TB, fn any, args ...any) *Value[T] {
	t.Helper()
	fnValue, reflectedArgs := prepareCall(t, fn, args...)
	if !fnValue.IsValid() {
		return nil
	}
	results := fnValue.Call(reflectedArgs)
	return buildValues[T](t, results)
}

func reflectArg(arg any, wantType reflect.Type) reflect.Value {
	if arg == nil {
		return reflect.Zero(wantType)
	}
	return reflect.ValueOf(arg)
}

func Call(t testing.TB, fn any, args ...any) *Value[interface{}] {
	t.Helper()
	return execute[interface{}](t, fn, args...)
}

func CallAs[T any](t testing.TB, fn any, args ...any) *Value[T] {
	t.Helper()
	return execute[T](t, fn, args...)
}

func Capture(t testing.TB, values ...any) *Value[interface{}] {
	t.Helper()

	current := &Value[interface{}]{
		t:       t,
		current: reflect.Value{},
		next:    nil,
		fatal:   false,
	}

	for _, value := range values {
		var rv reflect.Value
		if value == nil {
			rv = reflect.ValueOf(&value).Elem()
		} else {
			rv = reflect.ValueOf(value)
		}
		current = &Value[interface{}]{
			t:       t,
			current: rv,
			next:    current,
			fatal:   false,
		}
	}
	return current
}
