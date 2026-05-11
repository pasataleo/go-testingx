package mocks

import (
	"fmt"
	"testing"
)

type T struct {
	*testing.T

	ErrorMessage string
	FatalMessage string
}

func (testing *T) Error(args ...any) {
	testing.T.Helper()
	testing.ErrorMessage = fmt.Sprint(args...)
}

func (testing *T) Errorf(format string, args ...any) {
	testing.T.Helper()
	testing.ErrorMessage = fmt.Sprintf(format, args...)
}

func (testing *T) Fatal(args ...any) {
	testing.T.Helper()
	testing.FatalMessage = fmt.Sprint(args...)
}

func (testing *T) Fatalf(format string, args ...any) {
	testing.T.Helper()
	testing.FatalMessage = fmt.Sprintf(format, args...)
}
