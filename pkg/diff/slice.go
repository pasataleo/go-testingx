package diff

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pasataleo/go-testingx/pkg/render"
)

var (
	_ CompositeResult = (*sliceResult)(nil)
)

func ofSlice(got, want reflect.Value, opts *Opts) Result {
	if opts.TB != nil {
		opts.TB.Helper()
	}

	var entries []SliceEntry

	minLen := got.Len()
	if want.Len() < minLen {
		minLen = want.Len()
	}

	for i := 0; i < minLen; i++ {
		entries = append(entries, SliceEntry{
			Index:  i,
			Result: of(got.Index(i), want.Index(i), opts),
		})
	}
	for i := minLen; i < got.Len(); i++ {
		entries = append(entries, SliceEntry{
			Index:  i,
			Result: ExtraResult(got.Index(i)),
		})
	}
	for i := minLen; i < want.Len(); i++ {
		entries = append(entries, SliceEntry{
			Index:  i,
			Result: MissingResult(want.Index(i)),
		})
	}

	return &sliceResult{
		gotType:  got.Type(),
		wantType: want.Type(),
		entries:  entries,
	}
}

// SliceEntry pairs an index with its diff Result.
type SliceEntry struct {
	Index  int
	Result Result
}

// SliceResult creates a CompositeResult for a same-type slice or array diff.
func SliceResult(t reflect.Type, entries []SliceEntry) Result {
	return &sliceResult{
		gotType:  t,
		wantType: t,
		entries:  entries,
	}
}

type sliceResult struct {
	gotType  reflect.Type
	wantType reflect.Type
	entries  []SliceEntry
}

func (s *sliceResult) Composite() {}

func (s *sliceResult) Status() Status {
	for _, e := range s.entries {
		switch e.Result.Status() {
		case StatusChanged, StatusMissing, StatusExtra:
			return StatusChanged
		default:
			continue
		}
	}
	return StatusUnchanged
}

func renderSliceInternal(t reflect.Type, contents string, opts *render.Opts) string {
	if !opts.ShowTypes {
		return contents
	}
	return opts.Printf("{grey}%s({reset}%s{grey}){reset}", render.RenderType(t, opts), contents)
}

func (s *sliceResult) RenderGot(opts *render.Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	return s.renderSide(s.gotType, opts, func(r Result) string { return r.RenderGot(opts) }, func(r Result) bool {
		return r.Status() != StatusMissing
	})
}

func (s *sliceResult) RenderWant(opts *render.Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	return s.renderSide(s.wantType, opts, func(r Result) string { return r.RenderWant(opts) }, func(r Result) bool {
		return r.Status() != StatusExtra
	})
}

func (s *sliceResult) renderSide(t reflect.Type, opts *render.Opts, renderValue func(Result) string, include func(Result) bool) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}

	type element struct {
		index int
		value string
	}

	elements := make([]element, 0, len(s.entries))
	for _, e := range s.entries {
		if !include(e.Result) {
			continue
		}
		elements = append(elements, element{
			index: e.Index,
			value: renderValue(e.Result),
		})
	}

	if len(elements) == 0 {
		return renderSliceInternal(t, opts.Printf("{grey}%s{reset}", "[]"), opts)
	}

	if !opts.ShowContents {
		return renderSliceInternal(t, opts.Printf("{grey}%s{reset}", "[...]"), opts)
	}

	var builder strings.Builder
	if opts.ShowTypes {
		builder.WriteString(opts.Printf("{grey}%s([{reset}\n", render.RenderType(t, opts)))
	} else {
		builder.WriteString(opts.Print("{grey}[{reset}\n"))
	}
	if opts.ShowIndices {
		indexWidth := len(fmt.Sprintf("%d", len(elements)-1))
		for _, e := range elements {
			index := opts.Printf("{grey}[%*d]:{reset}", indexWidth, e.index)
			value := strings.ReplaceAll(e.value, "\n", "\n  "+strings.Repeat(" ", indexWidth+4))
			builder.WriteString(opts.Printf("  %s %s{grey},{reset}\n", index, value))
		}
	} else {
		for _, e := range elements {
			value := strings.ReplaceAll(e.value, "\n", "\n  ")
			builder.WriteString(opts.Printf("  %s{grey},{reset}\n", value))
		}
	}
	if opts.ShowTypes {
		builder.WriteString(opts.Print("{grey}]){reset}"))
	} else {
		builder.WriteString(opts.Print("{grey}]{reset}"))
	}
	return builder.String()
}

func (s *sliceResult) RenderDiff(opts *Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}

	type element struct {
		index   int
		content string
	}

	// Compute index width for alignment.
	maxIndex := 0
	for _, e := range s.entries {
		if e.Index > maxIndex {
			maxIndex = e.Index
		}
	}
	indexWidth := len(fmt.Sprintf("%d", maxIndex))

	renderIndex := func(i int) string {
		if !opts.RenderOpts().ShowIndices {
			return ""
		}
		return opts.RenderOpts().Printf("{grey}[%*d]:{reset} ", indexWidth, i)
	}

	elements := make([]element, 0, len(s.entries))
	for _, e := range s.entries {
		idx := renderIndex(e.Index)
		switch e.Result.Status() {
		case StatusUnchanged:
			if !opts.ShowUnchanged {
				continue
			}
			elements = append(elements, element{
				index:   e.Index,
				content: opts.RenderOpts().Printf("  %s%s{grey},{reset}", idx, e.Result.RenderGot(opts.RenderOpts())),
			})
		case StatusChanged:
			if _, ok := e.Result.(CompositeResult); !ok {
				elements = append(elements, element{
					index: e.Index,
					content: opts.RenderOpts().Printf("{red}-{reset} %s%s{grey},{reset}\n", idx, e.Result.RenderWant(opts.RenderOpts())) +
						opts.RenderOpts().Printf("{green}+{reset} %s%s{grey},{reset}", idx, e.Result.RenderGot(opts.RenderOpts())),
				})
			} else {
				elements = append(elements, element{
					index:   e.Index,
					content: opts.RenderOpts().Printf("{grey}~{reset} %s%s{grey},{reset}", idx, e.Result.RenderDiff(opts.AsNested())),
				})
			}
		case StatusMissing:
			prefix := opts.RenderOpts().Printf("{red}-{reset} ")
			rendered := prefixLines(e.Result.RenderWant(opts.RenderOpts()), prefix)
			elements = append(elements, element{
				index:   e.Index,
				content: opts.RenderOpts().Printf("{red}-{reset} %s%s{grey},{reset}", idx, rendered),
			})
		case StatusExtra:
			prefix := opts.RenderOpts().Printf("{green}+{reset} ")
			rendered := prefixLines(e.Result.RenderGot(opts.RenderOpts()), prefix)
			elements = append(elements, element{
				index:   e.Index,
				content: opts.RenderOpts().Printf("{green}+{reset} %s%s{grey},{reset}", idx, rendered),
			})
		}
	}

	sameType := s.gotType == s.wantType

	var prefix string
	if !opts.Nested() {
		if s.Status() == StatusChanged && sameType {
			prefix = opts.RenderOpts().Printf("{grey}~{reset} ")
		} else if s.Status() != StatusChanged {
			prefix = "  "
		}
	}

	if !opts.RenderOpts().ShowContents || len(elements) == 0 {
		if sameType {
			return prefix + renderSliceInternal(s.gotType, opts.RenderOpts().Printf("{grey}%s{reset}", "[...]"), opts.RenderOpts())
		}
		return prefix + opts.RenderOpts().Printf("{grey}%s{reset}", "[...]")
	}

	var builder strings.Builder
	if opts.RenderOpts().ShowTypes && sameType {
		builder.WriteString(opts.RenderOpts().Printf("%s{grey}%s([{reset}\n", prefix, render.RenderType(s.gotType, opts.RenderOpts())))
	} else if opts.RenderOpts().ShowTypes {
		builder.WriteString(opts.RenderOpts().Printf("%s{red}-{reset} {grey}%s([{reset}\n", prefix, render.RenderType(s.wantType, opts.RenderOpts())))
		builder.WriteString(opts.RenderOpts().Printf("{green}+{reset} {grey}%s([{reset}\n", render.RenderType(s.gotType, opts.RenderOpts())))
	} else {
		builder.WriteString(opts.RenderOpts().Printf("%s{grey}[{reset}\n", prefix))
	}
	for _, e := range elements {
		content := strings.ReplaceAll(e.content, "\n", "\n  ")
		builder.WriteString("  " + content + "\n")
	}
	if opts.RenderOpts().ShowTypes {
		builder.WriteString(opts.RenderOpts().Print("  {grey}]){reset}"))
	} else {
		builder.WriteString(opts.RenderOpts().Print("  {grey}]{reset}"))
	}
	return builder.String()
}
