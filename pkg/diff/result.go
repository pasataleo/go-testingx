package diff

import (
	render2 "github.com/pasataleo/go-testingx/pkg/render"
)

// Status describes the outcome of comparing a got value against a want value.
type Status int

const (
	StatusUnchanged Status = iota // got == want
	StatusChanged                 // -want, +got
	StatusMissing                 // in want but not got
	StatusExtra                   // in got but not want
)

// Result represents the outcome of diffing two values and provides methods
// to render the got value, want value, or a unified diff.
type Result interface {
	Status() Status
	RenderGot(opts *render2.Opts) string
	RenderWant(opts *render2.Opts) string
	RenderDiff(opts *Opts) string
}

// CompositeResult is implemented by Result types that represent composite
// structures (maps, structs, slices) whose diffs render as multi-line blocks
// with ~ prefix, as opposed to leaf types that render as - want / + got pairs.
type CompositeResult interface {
	Result
	Composite()
}
