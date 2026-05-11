package testingx

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pasataleo/go-errorsx/pkg/errorsx"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func (value *Value[Want]) asError() (error, bool) {
	value.t.Helper()

	if !value.current.IsValid() {
		value.t.Fatalf("no value to check: more assertions chained than values returned")
		return nil, false
	}

	if !value.current.Type().Implements(errorType) {
		value.t.Fatalf("expected error type, got %s", value.current.Type())
		return nil, false
	}

	if value.current.IsNil() {
		return nil, true
	}

	return value.current.Interface().(error), true
}

func (value *Value[Want]) NoError() *Value[Want] {
	value.t.Helper()

	err, ok := value.asError()
	if !ok {
		return value.next
	}

	if err != nil {
		value.Fail("expected no error, got %q", err.Error())
	}
	return value.next
}

func (value *Value[Want]) Error() *Value[Want] {
	value.t.Helper()

	err, ok := value.asError()
	if !ok {
		return value.next
	}

	if err == nil {
		value.Fail("expected error, got nil")
	}
	return value.next
}

func (value *Value[Want]) ErrorCode(code errorsx.Code) *Value[Want] {
	value.t.Helper()

	err, ok := value.asError()
	if !ok {
		return value.next
	}

	if err == nil {
		value.Fail("expected error with code %q, got nil", code)
		return value.next
	}

	got := errorsx.ErrorCode(err)
	if got != code {
		value.Fail("expected error code %q, got %q", code, got)
	}
	return value.next
}

func (value *Value[Want]) MatchesError(msg string) *Value[Want] {
	value.t.Helper()

	err, ok := value.asError()
	if !ok {
		return value.next
	}

	if err == nil {
		value.Fail("expected error %q, got nil", msg)
		return value.next
	}

	if err.Error() != msg {
		value.Fail("expected error %q, got %q", msg, err.Error())
	}
	return value.next
}

func (value *Value[Want]) MatchesErrorf(format string, args ...interface{}) *Value[Want] {
	value.t.Helper()
	return value.MatchesError(fmt.Sprintf(format, args...))
}

func (value *Value[Want]) ErrorContains(substring string) *Value[Want] {
	value.t.Helper()

	err, ok := value.asError()
	if !ok {
		return value.next
	}

	if err == nil {
		value.Fail("expected error containing %q, got nil", substring)
		return value.next
	}

	if !strings.Contains(err.Error(), substring) {
		value.Fail("expected error containing %q, got %q", substring, err.Error())
	}
	return value.next
}

func (value *Value[Want]) ErrorContainsf(format string, args ...interface{}) *Value[Want] {
	value.t.Helper()
	return value.ErrorContains(fmt.Sprintf(format, args...))
}

func (value *Value[Want]) HasError(msg string) *Value[Want] {
	value.t.Helper()

	err, ok := value.asError()
	if !ok {
		return value.next
	}

	if err == nil {
		value.Fail("expected aggregated error containing %q, got nil", msg)
		return value.next
	}

	for _, child := range errorsx.Errors(err) {
		if child.Error() == msg {
			return value.next
		}
	}

	value.Fail("expected error containing %q, got %q", msg, err.Error())
	return value.next
}

func (value *Value[Want]) HasErrorf(format string, args ...interface{}) *Value[Want] {
	value.t.Helper()
	return value.HasError(fmt.Sprintf(format, args...))
}
