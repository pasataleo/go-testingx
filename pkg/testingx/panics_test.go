package testingx

import (
	"testing"

	"github.com/pasataleo/go-testingx/pkg/mocks"
	render2 "github.com/pasataleo/go-testingx/pkg/render"
)

func TestPanics(t *testing.T) {
	t.Run("panics", func(t *testing.T) {
		mock := &mocks.T{T: t}
		value := Panics(mock, nil, func() { panic("boom") })
		value.Equal("boom")

		if mock.ErrorMessage != "" {
			t.Errorf("unexpected error: %s", mock.ErrorMessage)
		}
		if mock.FatalMessage != "" {
			t.Errorf("unexpected fatal: %s", mock.FatalMessage)
		}
	})

	t.Run("panics no panic (fatal)", func(t *testing.T) {
		mock := &mocks.T{T: t}
		Panics(mock, nil, func() {})

		if mock.FatalMessage == "" {
			t.Error("expected fatal message")
		}
		if want := "expected function to panic, but it did not"; mock.FatalMessage != want {
			t.Errorf("expected %q, got %q", want, mock.FatalMessage)
		}
	})
}

func TestNotPanics(t *testing.T) {
	t.Run("not panics", func(t *testing.T) {
		mock := &mocks.T{T: t}
		value := NotPanics(mock, nil, func() string { return "hello" })
		value.Equal("hello")

		if mock.ErrorMessage != "" {
			t.Errorf("unexpected error: %s", mock.ErrorMessage)
		}
		if mock.FatalMessage != "" {
			t.Errorf("unexpected fatal: %s", mock.FatalMessage)
		}
	})

	t.Run("not panics with panic (fatal)", func(t *testing.T) {
		mock := &mocks.T{T: t}
		opts := render2.NewOpts(mock, render2.DisableColour())
		NotPanics(mock, opts, func() { panic("boom") })

		if mock.FatalMessage == "" {
			t.Error("expected fatal message")
		}
		if want := `expected function not to panic, but it panicked with: string("boom")`; mock.FatalMessage != want {
			t.Errorf("expected %q, got %q", want, mock.FatalMessage)
		}
	})
}

func TestPanicsAs(t *testing.T) {
	t.Run("panics as string", func(t *testing.T) {
		mock := &mocks.T{T: t}
		value := PanicsAs[string](mock, nil, func() { panic("boom") })
		value.Equal("boom")

		if mock.ErrorMessage != "" {
			t.Errorf("unexpected error: %s", mock.ErrorMessage)
		}
		if mock.FatalMessage != "" {
			t.Errorf("unexpected fatal: %s", mock.FatalMessage)
		}
	})

	t.Run("panics as wrong type (fatal)", func(t *testing.T) {
		mock := &mocks.T{T: t}
		PanicsAs[int](mock, nil, func() { panic("boom") })

		if mock.FatalMessage == "" {
			t.Error("expected fatal message")
		}
	})

	t.Run("panics as no panic (fatal)", func(t *testing.T) {
		mock := &mocks.T{T: t}
		PanicsAs[string](mock, nil, func() {})

		if mock.FatalMessage == "" {
			t.Error("expected fatal message")
		}
		if want := "expected function to panic, but it did not"; mock.FatalMessage != want {
			t.Errorf("expected %q, got %q", want, mock.FatalMessage)
		}
	})
}

func TestNotPanicsAs(t *testing.T) {
	t.Run("not panics as", func(t *testing.T) {
		mock := &mocks.T{T: t}
		value := NotPanicsAs[string](mock, nil, func() string { return "hello" })
		value.Equal("hello")

		if mock.ErrorMessage != "" {
			t.Errorf("unexpected error: %s", mock.ErrorMessage)
		}
		if mock.FatalMessage != "" {
			t.Errorf("unexpected fatal: %s", mock.FatalMessage)
		}
	})

	t.Run("not panics as with panic (fatal)", func(t *testing.T) {
		mock := &mocks.T{T: t}
		opts := render2.NewOpts(mock, render2.DisableColour())
		NotPanicsAs[string](mock, opts, func() { panic("boom") })

		if mock.FatalMessage == "" {
			t.Error("expected fatal message")
		}
		if want := `expected function not to panic, but it panicked with: string("boom")`; mock.FatalMessage != want {
			t.Errorf("expected %q, got %q", want, mock.FatalMessage)
		}
	})
}
