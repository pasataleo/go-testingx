package render

import (
	"fmt"
	"reflect"
	"strings"
)

func renderSlice(value reflect.Value, opts *Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}

	if value.Kind() == reflect.Slice && value.IsNil() {
		return renderNil(value.Type(), opts)
	}

	if value.Len() == 0 {
		return renderInternal(value.Type(), opts.Printf("{grey}%s{reset}", "[]"), opts)
	}

	if !opts.ShowContents {
		return renderInternal(value.Type(), opts.Printf("{grey}%s{reset}", "[...]"), opts)
	}

	elements := make([]string, 0, value.Len())
	for i := 0; i < value.Len(); i++ {
		elements = append(elements, render(value.Index(i), opts))
	}

	var builder strings.Builder
	if opts.ShowTypes {
		builder.WriteString(opts.Printf("{grey}%s([{reset}\n", RenderType(value.Type(), opts)))
	} else {
		builder.WriteString(opts.Print("{grey}[{reset}\n"))
	}
	if opts.ShowIndices {
		// Compute index width for alignment.
		indexWidth := len(fmt.Sprintf("%d", len(elements)-1))

		for i, elem := range elements {
			index := opts.Printf("{grey}[%*d]:{reset}", indexWidth, i)
			elem = strings.ReplaceAll(elem, "\n", "\n  "+strings.Repeat(" ", indexWidth+4))
			builder.WriteString(opts.Printf("  %s %s{grey},{reset}\n", index, elem))
		}
	} else {
		for _, elem := range elements {
			elem = strings.ReplaceAll(elem, "\n", "\n  ")
			builder.WriteString(opts.Printf("  %s{grey},{reset}\n", elem))
		}
	}
	if opts.ShowTypes {
		builder.WriteString(opts.Print("{grey}]){reset}"))

	} else {
		builder.WriteString(opts.Print("{grey}]{reset}"))

	}
	return builder.String()
}
