package diff

import (
	"reflect"
	"strings"

	"github.com/pasataleo/go-testingx/pkg/render"
)

var (
	_ CompositeResult = (*structResult)(nil)
)

func ofStruct(got, want reflect.Value, opts *Opts) Result {
	if opts.TB != nil {
		opts.TB.Helper()
	}

	seen := make(map[string]bool)
	var fields []StructField

	for i := 0; i < got.NumField(); i++ {
		f := got.Type().Field(i)
		if !f.IsExported() {
			if !opts.RenderOpts().SkipUnexported {
				opts.Fatalf("unexported field %s.%s: use render.WithSkipUnexported to skip", got.Type().Name(), f.Name)
			}
			continue
		}
		if opts.RenderOpts().SkipField(got.Type().Name(), f.Name) {
			continue
		}
		seen[f.Name] = true
		if _, ok := want.Type().FieldByName(f.Name); ok {
			fields = append(fields, StructField{
				Name:   f.Name,
				Result: of(got.Field(i), want.FieldByName(f.Name), opts),
			})
		} else {
			fields = append(fields, StructField{
				Name:   f.Name,
				Result: ExtraResult(got.Field(i)),
			})
		}
	}

	// Append want-only fields.
	if got.Type() != want.Type() {
		wantName := want.Type().Name()
		for i := 0; i < want.NumField(); i++ {
			f := want.Type().Field(i)
			if !f.IsExported() || seen[f.Name] {
				continue
			}
			if opts.RenderOpts().SkipField(wantName, f.Name) {
				continue
			}
			fields = append(fields, StructField{
				Name:   f.Name,
				Result: MissingResult(want.Field(i)),
			})
		}
	}

	return &structResult{
		gotType:  got.Type(),
		wantType: want.Type(),
		fields:   fields,
	}
}

// StructField pairs a field name with its diff Result.
type StructField struct {
	Name   string
	Result Result
}

// StructResult creates a CompositeResult for a same-type struct diff.
func StructResult(t reflect.Type, fields []StructField) Result {
	return &structResult{
		gotType:  t,
		wantType: t,
		fields:   fields,
	}
}

type structResult struct {
	gotType  reflect.Type
	wantType reflect.Type
	fields   []StructField
}

func (s *structResult) Composite() {}

func (s *structResult) Status() Status {
	for _, f := range s.fields {
		switch f.Result.Status() {
		case StatusChanged, StatusMissing, StatusExtra:
			return StatusChanged
		default:
			continue
		}
	}
	return StatusUnchanged
}

func (s *structResult) RenderGot(opts *render.Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	return s.renderSide(s.gotType, opts, func(r Result) string { return r.RenderGot(opts) }, func(r Result) bool {
		return r.Status() != StatusMissing
	})
}

func (s *structResult) RenderWant(opts *render.Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	return s.renderSide(s.wantType, opts, func(r Result) string { return r.RenderWant(opts) }, func(r Result) bool {
		return r.Status() != StatusExtra
	})
}

func renderStructInternal(t reflect.Type, contents string, opts *render.Opts) string {
	if !opts.ShowTypes {
		return contents
	}
	return opts.Printf("{grey}%s({reset}%s{grey}){reset}", render.RenderType(t, opts), contents)
}

func (s *structResult) renderSide(t reflect.Type, opts *render.Opts, renderValue func(Result) string, include func(Result) bool) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}

	type field struct {
		name, value string
	}

	fields := make([]field, 0, len(s.fields))
	for _, f := range s.fields {
		if !include(f.Result) {
			continue
		}
		fields = append(fields, field{
			name:  f.Name,
			value: renderValue(f.Result),
		})
	}

	if len(fields) == 0 {
		return renderStructInternal(t, opts.Printf("{grey}%s{reset}", "{}"), opts)
	}

	if !opts.ShowContents {
		return renderStructInternal(t, opts.Printf("{grey}%s{reset}", "{...}"), opts)
	}

	maxNameLen := 0
	for _, f := range fields {
		if len(f.name) > maxNameLen {
			maxNameLen = len(f.name)
		}
	}

	var builder strings.Builder
	if opts.ShowTypes {
		builder.WriteString(opts.Printf("{grey}%s({{reset}\n", render.RenderType(t, opts)))
	} else {
		builder.WriteString(opts.Print("{grey}{{reset}\n"))
	}
	for _, f := range fields {
		value := strings.ReplaceAll(f.value, "\n", "\n  ")
		builder.WriteString(opts.Printf("  %-*s{grey}:{reset} %s{grey},{reset}\n", maxNameLen, f.name, value))
	}
	if opts.ShowTypes {
		builder.WriteString(opts.Print("{grey}}){reset}"))
	} else {
		builder.WriteString(opts.Print("{grey}}{reset}"))
	}
	return builder.String()
}

func (s *structResult) RenderDiff(opts *Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}

	type field struct {
		name    string
		content string
	}

	fields := make([]field, 0, len(s.fields))
	for _, f := range s.fields {
		switch f.Result.Status() {
		case StatusUnchanged:
			if !opts.ShowUnchanged {
				continue
			}
			fields = append(fields, field{
				name:    f.Name,
				content: opts.RenderOpts().Printf("  %-*s{grey}:{reset} %s{grey},{reset}", 0, f.Name, f.Result.RenderGot(opts.RenderOpts())),
			})
		case StatusChanged:
			if _, ok := f.Result.(CompositeResult); !ok {
				fields = append(fields, field{
					name: f.Name,
					content: opts.RenderOpts().Printf("{red}-{reset} %-*s{grey}:{reset} %s{grey},{reset}\n", 0, f.Name, f.Result.RenderWant(opts.RenderOpts())) +
						opts.RenderOpts().Printf("{green}+{reset} %-*s{grey}:{reset} %s{grey},{reset}", 0, f.Name, f.Result.RenderGot(opts.RenderOpts())),
				})
			} else {
				fields = append(fields, field{
					name:    f.Name,
					content: opts.RenderOpts().Printf("{grey}~{reset} %-*s{grey}:{reset} %s{grey},{reset}", 0, f.Name, f.Result.RenderDiff(opts.AsNested())),
				})
			}
		case StatusMissing:
			prefix := opts.RenderOpts().Printf("{red}-{reset} ")
			rendered := prefixLines(f.Result.RenderWant(opts.RenderOpts()), prefix)
			fields = append(fields, field{
				name:    f.Name,
				content: opts.RenderOpts().Printf("{red}-{reset} %-*s{grey}:{reset} %s{grey},{reset}", 0, f.Name, rendered),
			})
		case StatusExtra:
			prefix := opts.RenderOpts().Printf("{green}+{reset} ")
			rendered := prefixLines(f.Result.RenderGot(opts.RenderOpts()), prefix)
			fields = append(fields, field{
				name:    f.Name,
				content: opts.RenderOpts().Printf("{green}+{reset} %-*s{grey}:{reset} %s{grey},{reset}", 0, f.Name, rendered),
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

	if !opts.RenderOpts().ShowContents || len(fields) == 0 {
		if sameType {
			return prefix + renderStructInternal(s.gotType, opts.RenderOpts().Printf("{grey}%s{reset}", "{...}"), opts.RenderOpts())
		}
		return prefix + opts.RenderOpts().Printf("{grey}%s{reset}", "{...}")
	}

	var builder strings.Builder
	if opts.RenderOpts().ShowTypes && sameType {
		builder.WriteString(opts.RenderOpts().Printf("%s{grey}%s({{reset}\n", prefix, render.RenderType(s.gotType, opts.RenderOpts())))
	} else if opts.RenderOpts().ShowTypes {
		builder.WriteString(opts.RenderOpts().Printf("%s{red}-{reset} {grey}%s({{reset}\n", prefix, render.RenderType(s.wantType, opts.RenderOpts())))
		builder.WriteString(opts.RenderOpts().Printf("{green}+{reset} {grey}%s({{reset}\n", render.RenderType(s.gotType, opts.RenderOpts())))
	} else {
		builder.WriteString(opts.RenderOpts().Printf("%s{grey}{{reset}\n", prefix))
	}
	for _, f := range fields {
		content := strings.ReplaceAll(f.content, "\n", "\n  ")
		builder.WriteString("  " + content + "\n")
	}
	if opts.RenderOpts().ShowTypes {
		builder.WriteString(opts.RenderOpts().Print("  {grey}}){reset}"))
	} else {
		builder.WriteString(opts.RenderOpts().Print("  {grey}}{reset}"))
	}
	return builder.String()
}
