package testingx

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/pasataleo/go-errorsx/pkg/errorsx"
	"github.com/pasataleo/go-testingx/pkg/mocks"
	"github.com/pasataleo/go-testingx/pkg/render"
)

type hasEqual struct{ V int }

func (h hasEqual) Equal(other hasEqual) bool { return h.V == other.V }

type hasContains struct{ items []int }

func (h hasContains) Contains(item int) bool {
	for _, i := range h.items {
		if i == item {
			return true
		}
	}
	return false
}

func TestCall(t *testing.T) {
	tcs := map[string]struct {
		fn        interface{}
		args      []interface{}
		validate  func(value *Value[interface{}])
		wantError string
		wantFatal string
	}{
		// Equal
		"equal": {
			fn: func() string { return "hello" },
			validate: func(value *Value[interface{}]) {
				value.Equal("hello")
			},
		},
		"equal (error)": {
			fn: func() string { return "hello" },
			validate: func(value *Value[interface{}]) {
				value.Equal("world")
			},
			wantError: "expected values to be equal",
		},
		"equal (fatal)": {
			fn: func() string { return "hello" },
			validate: func(value *Value[interface{}]) {
				value.Fatal().Equal("world")
			},
			wantFatal: "expected values to be equal",
		},
		"equal typed nil with Equal method": {
			fn: func() hasEqual { return hasEqual{V: 1} },
			validate: func(value *Value[interface{}]) {
				value.Equal((*hasEqual)(nil))
			},
			wantError: "expected values to be equal",
		},
		"equal nil does not panic": {
			fn: func() hasEqual { return hasEqual{V: 1} },
			validate: func(value *Value[interface{}]) {
				value.Equal(nil)
			},
			wantError: "expected values to be equal",
		},

		// NotEqual
		"not equal": {
			fn: func() string { return "hello" },
			validate: func(value *Value[interface{}]) {
				value.NotEqual("world")
			},
		},
		"not equal (error)": {
			fn: func() string { return "hello" },
			validate: func(value *Value[interface{}]) {
				value.NotEqual("hello")
			},
			wantError: "expected values to not be equal",
		},
		"not equal (fatal)": {
			fn: func() string { return "hello" },
			validate: func(value *Value[interface{}]) {
				value.Fatal().NotEqual("hello")
			},
			wantFatal: "expected values to not be equal",
		},

		// True
		"true": {
			fn: func() bool { return true },
			validate: func(value *Value[interface{}]) {
				value.True()
			},
		},
		"true (error)": {
			fn: func() bool { return false },
			validate: func(value *Value[interface{}]) {
				value.True()
			},
			wantError: "expected true, got false",
		},
		"true (fatal)": {
			fn: func() bool { return false },
			validate: func(value *Value[interface{}]) {
				value.Fatal().True()
			},
			wantFatal: "expected true, got false",
		},

		// False
		"false": {
			fn: func() bool { return false },
			validate: func(value *Value[interface{}]) {
				value.False()
			},
		},
		"false (error)": {
			fn: func() bool { return true },
			validate: func(value *Value[interface{}]) {
				value.False()
			},
			wantError: "expected false, got true",
		},
		"false (fatal)": {
			fn: func() bool { return true },
			validate: func(value *Value[interface{}]) {
				value.Fatal().False()
			},
			wantFatal: "expected false, got true",
		},

		// Nil
		"nil": {
			fn: func() error { return nil },
			validate: func(value *Value[interface{}]) {
				value.Nil()
			},
		},
		"nil (error)": {
			fn: func() *string {
				s := "hello"
				return &s
			},
			validate: func(value *Value[interface{}]) {
				value.Nil()
			},
			wantError: "expected nil value",
		},
		"nil (fatal)": {
			fn: func() *string {
				s := "hello"
				return &s
			},
			validate: func(value *Value[interface{}]) {
				value.Fatal().Nil()
			},
			wantFatal: "expected nil value",
		},

		// NotNil
		"not nil": {
			fn: func() *string {
				s := "hello"
				return &s
			},
			validate: func(value *Value[interface{}]) {
				value.NotNil()
			},
		},
		"not nil (error)": {
			fn: func() error { return nil },
			validate: func(value *Value[interface{}]) {
				value.NotNil()
			},
			wantError: "expected non-nil value, got nil",
		},
		"not nil (fatal)": {
			fn: func() error { return nil },
			validate: func(value *Value[interface{}]) {
				value.Fatal().NotNil()
			},
			wantFatal: "expected non-nil value, got nil",
		},

		// Len
		"len": {
			fn: func() []int { return []int{1, 2, 3} },
			validate: func(value *Value[interface{}]) {
				value.Len(3)
			},
		},
		"len (error)": {
			fn: func() []int { return []int{1, 2, 3} },
			validate: func(value *Value[interface{}]) {
				value.Len(5)
			},
			wantError: "expected length 5, got 3",
		},
		"len (fatal)": {
			fn: func() []int { return []int{1, 2, 3} },
			validate: func(value *Value[interface{}]) {
				value.Fatal().Len(5)
			},
			wantFatal: "expected length 5, got 3",
		},

		// Empty
		"empty": {
			fn: func() []int { return []int{} },
			validate: func(value *Value[interface{}]) {
				value.Empty()
			},
		},
		"empty (error)": {
			fn: func() []int { return []int{1, 2} },
			validate: func(value *Value[interface{}]) {
				value.Empty()
			},
			wantError: "expected empty, got length 2",
		},
		"empty (fatal)": {
			fn: func() []int { return []int{1, 2} },
			validate: func(value *Value[interface{}]) {
				value.Fatal().Empty()
			},
			wantFatal: "expected empty, got length 2",
		},

		// NotEmpty
		"not empty": {
			fn: func() []int { return []int{1} },
			validate: func(value *Value[interface{}]) {
				value.NotEmpty()
			},
		},
		"not empty (error)": {
			fn: func() []int { return []int{} },
			validate: func(value *Value[interface{}]) {
				value.NotEmpty()
			},
			wantError: "expected non-empty, got length 0",
		},
		"not empty (fatal)": {
			fn: func() []int { return []int{} },
			validate: func(value *Value[interface{}]) {
				value.Fatal().NotEmpty()
			},
			wantFatal: "expected non-empty, got length 0",
		},

		// Contains (string)
		"contains string": {
			fn: func() string { return "hello world" },
			validate: func(value *Value[interface{}]) {
				value.Contains("world")
			},
		},
		"contains string (error)": {
			fn: func() string { return "hello world" },
			validate: func(value *Value[interface{}]) {
				value.Contains("foo", render.DisableColour())
			},
			wantError: `expected string("hello world") to contain string("foo")`,
		},

		// Contains (slice)
		"contains slice": {
			fn: func() []int { return []int{1, 2, 3} },
			validate: func(value *Value[interface{}]) {
				value.Contains(2)
			},
		},
		"contains slice (error)": {
			fn: func() []int { return []int{1, 2, 3} },
			validate: func(value *Value[interface{}]) {
				value.Contains(5, render.DisableColour())
			},
			wantError: "expected collection to contain int(5)",
		},

		// Contains (map)
		"contains map": {
			fn: func() map[string]int { return map[string]int{"a": 1} },
			validate: func(value *Value[interface{}]) {
				value.Contains("a")
			},
		},
		"contains map (error)": {
			fn: func() map[string]int { return map[string]int{"a": 1} },
			validate: func(value *Value[interface{}]) {
				value.Contains("b", render.DisableColour())
			},
			wantError: `expected map to contain key string("b")`,
		},
		"contains nil with Contains method": {
			fn: func() hasContains { return hasContains{items: []int{1}} },
			validate: func(value *Value[interface{}]) {
				value.Contains(nil)
			},
			wantFatal: "cannot check contains on value of kind struct",
		},

		// NotContains (string)
		"not contains string": {
			fn: func() string { return "hello world" },
			validate: func(value *Value[interface{}]) {
				value.NotContains("foo")
			},
		},
		"not contains string (error)": {
			fn: func() string { return "hello world" },
			validate: func(value *Value[interface{}]) {
				value.NotContains("world", render.DisableColour())
			},
			wantError: `expected string("hello world") to not contain string("world")`,
		},

		// NotContains (slice)
		"not contains slice": {
			fn: func() []int { return []int{1, 2, 3} },
			validate: func(value *Value[interface{}]) {
				value.NotContains(5)
			},
		},
		"not contains slice (error)": {
			fn: func() []int { return []int{1, 2, 3} },
			validate: func(value *Value[interface{}]) {
				value.NotContains(2, render.DisableColour())
			},
			wantError: "expected collection to not contain int(2)",
		},

		// NotContains (map)
		"not contains map": {
			fn: func() map[string]int { return map[string]int{"a": 1} },
			validate: func(value *Value[interface{}]) {
				value.NotContains("b")
			},
		},
		"not contains map (error)": {
			fn: func() map[string]int { return map[string]int{"a": 1} },
			validate: func(value *Value[interface{}]) {
				value.NotContains("a", render.DisableColour())
			},
			wantError: `expected map to not contain key string("a")`,
		},

		// NoError
		"no error": {
			fn: func() error { return nil },
			validate: func(value *Value[interface{}]) {
				value.NoError()
			},
		},
		"no error (error)": {
			fn: func() error { return errors.New("something broke") },
			validate: func(value *Value[interface{}]) {
				value.NoError()
			},
			wantError: `expected no error, got "something broke"`,
		},
		"no error (fatal)": {
			fn: func() error { return errors.New("something broke") },
			validate: func(value *Value[interface{}]) {
				value.Fatal().NoError()
			},
			wantFatal: `expected no error, got "something broke"`,
		},

		// Error
		"error": {
			fn: func() error { return errors.New("something broke") },
			validate: func(value *Value[interface{}]) {
				value.Error()
			},
		},
		"error (error)": {
			fn: func() error { return nil },
			validate: func(value *Value[interface{}]) {
				value.Error()
			},
			wantError: "expected error, got nil",
		},
		"error (fatal)": {
			fn: func() error { return nil },
			validate: func(value *Value[interface{}]) {
				value.Fatal().Error()
			},
			wantFatal: "expected error, got nil",
		},

		// MatchesError
		"matches error": {
			fn: func() error { return errors.New("specific error") },
			validate: func(value *Value[interface{}]) {
				value.MatchesError("specific error")
			},
		},
		"matches error (error)": {
			fn: func() error { return errors.New("specific error") },
			validate: func(value *Value[interface{}]) {
				value.MatchesError("different error")
			},
			wantError: `expected error "different error", got "specific error"`,
		},
		"matches error nil (error)": {
			fn: func() error { return nil },
			validate: func(value *Value[interface{}]) {
				value.MatchesError("specific error")
			},
			wantError: `expected error "specific error", got nil`,
		},

		// MatchesErrorf
		"matches errorf": {
			fn: func() error { return errors.New("error 42") },
			validate: func(value *Value[interface{}]) {
				value.MatchesErrorf("error %d", 42)
			},
		},

		// ErrorCode
		"error code": {
			fn: func() error { return errorsx.New(errorsx.NotFound, nil, "not found") },
			validate: func(value *Value[interface{}]) {
				value.ErrorCode(errorsx.NotFound)
			},
		},
		"error code (error)": {
			fn: func() error { return errorsx.New(errorsx.NotFound, nil, "not found") },
			validate: func(value *Value[interface{}]) {
				value.ErrorCode(errorsx.Internal)
			},
			wantError: `expected error code "internal", got "not_found"`,
		},
		"error code nil (error)": {
			fn: func() error { return nil },
			validate: func(value *Value[interface{}]) {
				value.ErrorCode(errorsx.NotFound)
			},
			wantError: `expected error with code "not_found", got nil`,
		},

		// HasError
		"contains error": {
			fn: func() error {
				return errorsx.Append(errors.New("first"), errors.New("second"))
			},
			validate: func(value *Value[interface{}]) {
				value.HasError("second")
			},
		},
		"contains error (error)": {
			fn: func() error {
				return errorsx.Append(errors.New("first"), errors.New("second"))
			},
			validate: func(value *Value[interface{}]) {
				value.HasError("third")
			},
			wantError: `expected error containing "third"`,
		},
		"contains error nil (error)": {
			fn: func() error { return nil },
			validate: func(value *Value[interface{}]) {
				value.HasError("something")
			},
			wantError: `expected aggregated error containing "something", got nil`,
		},

		// HasErrorf
		"contains errorf": {
			fn: func() error {
				return errorsx.Append(errors.New("error 42"))
			},
			validate: func(value *Value[interface{}]) {
				value.HasErrorf("error %d", 42)
			},
		},

		// Multiple return values
		"multiple returns": {
			fn: func() (string, error) { return "hello", nil },
			validate: func(value *Value[interface{}]) {
				value.NoError().Equal("hello")
			},
		},
		"multiple returns (error)": {
			fn: func() (string, error) { return "", fmt.Errorf("failed") },
			validate: func(value *Value[interface{}]) {
				value.NoError()
			},
			wantError: `expected no error, got "failed"`,
		},

		// Capture
		"capture": {
			fn: func() string { return "hello" },
			validate: func(value *Value[interface{}]) {
				got := value.Capture()
				if got != "hello" {
					panic("unexpected value")
				}
			},
		},
	}
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			mock := &mocks.T{T: t}
			value := Call(mock, tc.fn, tc.args...)

			tc.validate(value)

			if len(tc.wantError) > 0 {
				if !strings.Contains(mock.ErrorMessage, tc.wantError) {
					t.Errorf("\nwant: %s\ngot: %s", tc.wantError, mock.ErrorMessage)
				}
			} else if len(mock.ErrorMessage) > 0 {
				t.Errorf("\nwant: %s\ngot: %s", tc.wantError, mock.ErrorMessage)
			}

			if len(tc.wantFatal) > 0 {
				if !strings.Contains(mock.FatalMessage, tc.wantFatal) {
					t.Errorf("\nwant: %s\ngot:  %s", tc.wantFatal, mock.FatalMessage)
				}
			} else if len(mock.FatalMessage) > 0 {
				t.Errorf("\n got: %v\nwant: %v", mock.FatalMessage, tc.wantFatal)
			}
		})
	}
}

func TestCallAs(t *testing.T) {
	mock := &mocks.T{T: t}
	value := CallAs[string](mock, func() string { return "hello" })

	got := value.Capture()
	if got != "hello" {
		t.Errorf("expected %q, got %q", "hello", got)
	}
}

func TestCapture(t *testing.T) {
	mock := &mocks.T{T: t}
	value := Capture(mock, "hello", 42)

	value.Equal(42).Equal("hello")

	if mock.ErrorMessage != "" {
		t.Errorf("unexpected error: %s", mock.ErrorMessage)
	}
}
