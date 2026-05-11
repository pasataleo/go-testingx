package testingx

import (
	"testing"

	"github.com/pasataleo/go-testingx/pkg/mocks"
)

func TestValidate(t *testing.T) {
	t.Run("validate", func(t *testing.T) {
		mock := &mocks.T{T: t}
		Call(mock, func() string { return "hello" }).Validate(func(t testing.TB, got string) {
			if got != "hello" {
				t.Errorf("expected %q, got %q", "hello", got)
			}
		})

		if mock.ErrorMessage != "" {
			t.Errorf("unexpected error: %s", mock.ErrorMessage)
		}
		if mock.FatalMessage != "" {
			t.Errorf("unexpected fatal: %s", mock.FatalMessage)
		}
	})

	t.Run("validate interface", func(t *testing.T) {
		mock := &mocks.T{T: t}
		Call(mock, func() string { return "hello" }).Validate(func(t testing.TB, got interface{}) {
			if got != "hello" {
				t.Errorf("expected %q, got %q", "hello", got)
			}
		})

		if mock.ErrorMessage != "" {
			t.Errorf("unexpected error: %s", mock.ErrorMessage)
		}
		if mock.FatalMessage != "" {
			t.Errorf("unexpected fatal: %s", mock.FatalMessage)
		}
	})

	t.Run("validate (error)", func(t *testing.T) {
		mock := &mocks.T{T: t}
		Call(mock, func() string { return "hello" }).Validate(func(t testing.TB, got string) {
			t.Errorf("validation failed: got %q", got)
		})

		if mock.ErrorMessage == "" {
			t.Error("expected error message")
		}
	})
}
