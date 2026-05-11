package render

import (
	"reflect"
	"sort"
	"strings"
)

func renderMap(value reflect.Value, opts *Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	type entry struct {
		key, value string
	}

	if value.IsNil() {
		return renderNil(value.Type(), opts)
	}

	if value.Len() == 0 {
		return renderInternal(value.Type(), opts.Printf("{grey}%s{reset}", "{}"), opts)
	}

	if !opts.ShowContents {
		return renderInternal(value.Type(), opts.Printf("{grey}%s{reset}", "{...}"), opts)
	}

	entries := make([]entry, 0, value.Len())
	for _, key := range value.MapKeys() {
		entries = append(entries, entry{
			key:   render(key, opts),
			value: render(value.MapIndex(key), opts),
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].key < entries[j].key
	})

	hasMultilineKey := false
	maxKeyLen := 0
	for _, entry := range entries {
		if strings.Contains(entry.key, "\n") {
			hasMultilineKey = true
		}
		if len(entry.key) > maxKeyLen {
			maxKeyLen = len(entry.key)
		}
	}
	if hasMultilineKey {
		maxKeyLen = 0
	}

	var builder strings.Builder
	if opts.ShowTypes {
		builder.WriteString(opts.Printf("{grey}%s({{reset}\n", RenderType(value.Type(), opts)))
	} else {
		builder.WriteString(opts.Print("{grey}{{reset}\n"))
	}
	for _, entry := range entries {
		key := strings.ReplaceAll(entry.key, "\n", "\n  ")
		value := strings.ReplaceAll(entry.value, "\n", "\n  ")
		builder.WriteString(opts.Printf("  %-*s{grey}:{reset} %s{grey},{reset}\n", maxKeyLen, key, value))
	}
	if opts.ShowTypes {
		builder.WriteString(opts.Print("{grey}}){reset}"))

	} else {
		builder.WriteString(opts.Print("{grey}}{reset}"))

	}
	return builder.String()
}
