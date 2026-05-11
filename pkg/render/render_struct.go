package render

import (
	"reflect"
	"strings"
)

func renderStruct(value reflect.Value, opts *Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	type field struct {
		name, value string
	}

	structName := value.Type().Name()
	fields := make([]field, 0, value.NumField())
	for i := 0; i < value.NumField(); i++ {
		f := value.Type().Field(i)
		if !f.IsExported() {
			if !opts.SkipUnexported {
				opts.Fatalf("unexported field %s.%s: use WithSkipUnexported to skip", structName, f.Name)
			}
			continue
		}
		if opts.SkipField(structName, f.Name) {
			continue
		}
		fields = append(fields, field{
			name:  f.Name,
			value: render(value.Field(i), opts),
		})
	}

	if len(fields) == 0 {
		return renderInternal(value.Type(), opts.Printf("{grey}%s{reset}", "{}"), opts)
	}

	if !opts.ShowContents {
		return renderInternal(value.Type(), opts.Printf("{grey}%s{reset}", "{...}"), opts)
	}

	maxNameLen := 0
	for _, field := range fields {
		if len(field.name) > maxNameLen {
			maxNameLen = len(field.name)
		}
	}

	var builder strings.Builder
	if opts.ShowTypes {
		builder.WriteString(opts.Printf("{grey}%s({{reset}\n", RenderType(value.Type(), opts)))
	} else {
		builder.WriteString(opts.Print("{grey}{{reset}\n"))
	}
	for _, field := range fields {
		value := strings.ReplaceAll(field.value, "\n", "\n  ")
		builder.WriteString(opts.Printf("  %-*s{grey}:{reset} %s{grey},{reset}\n", maxNameLen, field.name, value))
	}
	if opts.ShowTypes {
		builder.WriteString(opts.Print("{grey}}){reset}"))

	} else {
		builder.WriteString(opts.Print("{grey}}{reset}"))

	}
	return builder.String()
}
