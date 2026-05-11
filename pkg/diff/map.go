package diff

import (
	"reflect"
	"sort"
	"strings"

	"github.com/pasataleo/go-testingx/pkg/render"
)

var (
	_ CompositeResult = (*mapResult)(nil)
)

func renderMapInternal(t reflect.Type, contents string, opts *render.Opts) string {
	if !opts.ShowTypes {
		return contents
	}
	return opts.Printf("{grey}%s({reset}%s{grey}){reset}", render.RenderType(t, opts), contents)
}

func ofMap(got, want reflect.Value, opts *Opts) Result {
	if opts.TB != nil {
		opts.TB.Helper()
	}

	entries := make(map[interface{}]Result)

	for _, key := range want.MapKeys() {
		k := key.Interface()
		gotVal := got.MapIndex(key)
		wantVal := want.MapIndex(key)
		if !gotVal.IsValid() {
			entries[k] = MissingResult(wantVal)
		} else {
			entries[k] = of(gotVal, wantVal, opts)
		}
	}

	for _, key := range got.MapKeys() {
		k := key.Interface()
		if _, ok := entries[k]; ok {
			continue
		}
		entries[k] = ExtraResult(got.MapIndex(key))
	}

	return &mapResult{
		gotType:  got.Type(),
		wantType: want.Type(),
		entries:  entries,
	}
}

func MapResult(t reflect.Type, entries map[interface{}]Result) Result {
	return &mapResult{
		gotType:  t,
		wantType: t,
		entries:  entries,
	}
}

type mapResult struct {
	gotType  reflect.Type
	wantType reflect.Type
	entries  map[interface{}]Result
}

func (m *mapResult) Composite() {}

func (m *mapResult) Status() Status {
	for _, result := range m.entries {
		switch result.Status() {
		case StatusChanged, StatusMissing, StatusExtra:
			return StatusChanged
		default:
			continue
		}
	}
	return StatusUnchanged
}

func (m *mapResult) RenderGot(opts *render.Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	return m.renderSide(m.gotType, opts, func(r Result) string { return r.RenderGot(opts) }, func(r Result) bool {
		return r.Status() != StatusMissing
	})
}

func (m *mapResult) RenderWant(opts *render.Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}
	return m.renderSide(m.wantType, opts, func(r Result) string { return r.RenderWant(opts) }, func(r Result) bool {
		return r.Status() != StatusExtra
	})
}

func (m *mapResult) renderSide(t reflect.Type, opts *render.Opts, renderValue func(Result) string, include func(Result) bool) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}

	type entry struct {
		key, value string
	}

	entries := make([]entry, 0, len(m.entries))
	for key, result := range m.entries {
		if !include(result) {
			continue
		}
		entries = append(entries, entry{
			key:   render.Render(key, opts),
			value: renderValue(result),
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].key < entries[j].key
	})

	if len(entries) == 0 {
		return renderMapInternal(t, opts.Printf("{grey}%s{reset}", "{}"), opts)
	}

	if !opts.ShowContents {
		return renderMapInternal(t, opts.Printf("{grey}%s{reset}", "{...}"), opts)
	}

	maxKeyLen := 0
	for _, e := range entries {
		if !strings.Contains(e.key, "\n") && len(e.key) > maxKeyLen {
			maxKeyLen = len(e.key)
		}
	}

	var builder strings.Builder
	if opts.ShowTypes {
		builder.WriteString(opts.Printf("{grey}%s({{reset}\n", render.RenderType(t, opts)))
	} else {
		builder.WriteString(opts.Print("{grey}{{reset}\n"))
	}
	for _, e := range entries {
		key := strings.ReplaceAll(e.key, "\n", "\n  ")
		value := strings.ReplaceAll(e.value, "\n", "\n  ")
		builder.WriteString(opts.Printf("  %-*s{grey}:{reset} %s{grey},{reset}\n", maxKeyLen, key, value))
	}
	if opts.ShowTypes {
		builder.WriteString(opts.Print("{grey}}){reset}"))
	} else {
		builder.WriteString(opts.Print("{grey}}{reset}"))
	}
	return builder.String()
}

func (m *mapResult) RenderDiff(opts *Opts) string {
	if opts.TB != nil {
		opts.TB.Helper()
	}

	type entry struct {
		key     string
		content string
	}

	entries := make([]entry, 0, len(m.entries))
	for key, result := range m.entries {
		renderedKey := render.Render(key, opts.RenderOpts())
		switch result.Status() {
		case StatusUnchanged:
			if !opts.ShowUnchanged {
				continue
			}
			entries = append(entries, entry{
				key:     renderedKey,
				content: opts.RenderOpts().Printf("  %-*s{grey}:{reset} %s{grey},{reset}", 0, renderedKey, result.RenderGot(opts.RenderOpts())),
			})
		case StatusChanged:
			if _, ok := result.(CompositeResult); !ok {
				entries = append(entries, entry{
					key: renderedKey,
					content: opts.RenderOpts().Printf("{red}-{reset} %-*s{grey}:{reset} %s{grey},{reset}\n", 0, renderedKey, result.RenderWant(opts.RenderOpts())) +
						opts.RenderOpts().Printf("{green}+{reset} %-*s{grey}:{reset} %s{grey},{reset}", 0, renderedKey, result.RenderGot(opts.RenderOpts())),
				})
			} else {
				entries = append(entries, entry{
					key:     renderedKey,
					content: opts.RenderOpts().Printf("{grey}~{reset} %-*s{grey}:{reset} %s{grey},{reset}", 0, renderedKey, result.RenderDiff(opts.AsNested())),
				})
			}
		case StatusMissing:
			prefix := opts.RenderOpts().Printf("{red}-{reset} ")
			rendered := prefixLines(result.RenderWant(opts.RenderOpts()), prefix)
			entries = append(entries, entry{
				key:     renderedKey,
				content: opts.RenderOpts().Printf("{red}-{reset} %-*s{grey}:{reset} %s{grey},{reset}", 0, renderedKey, rendered),
			})
		case StatusExtra:
			prefix := opts.RenderOpts().Printf("{green}+{reset} ")
			rendered := prefixLines(result.RenderGot(opts.RenderOpts()), prefix)
			entries = append(entries, entry{
				key:     renderedKey,
				content: opts.RenderOpts().Printf("{green}+{reset} %-*s{grey}:{reset} %s{grey},{reset}", 0, renderedKey, rendered),
			})
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].key < entries[j].key
	})

	sameType := m.gotType == m.wantType

	var prefix string
	if !opts.Nested() {
		if m.Status() == StatusChanged && sameType {
			prefix = opts.RenderOpts().Printf("{grey}~{reset} ")
		} else if m.Status() != StatusChanged {
			prefix = "  "
		}
	}

	if !opts.RenderOpts().ShowContents || len(entries) == 0 {
		if sameType {
			return prefix + renderMapInternal(m.gotType, opts.RenderOpts().Printf("{grey}%s{reset}", "{...}"), opts.RenderOpts())
		}
		return prefix + opts.RenderOpts().Printf("{grey}%s{reset}", "{...}")
	}

	var builder strings.Builder
	if opts.RenderOpts().ShowTypes && sameType {
		builder.WriteString(opts.RenderOpts().Printf("%s{grey}%s({{reset}\n", prefix, render.RenderType(m.gotType, opts.RenderOpts())))
	} else if opts.RenderOpts().ShowTypes {
		builder.WriteString(opts.RenderOpts().Printf("%s{red}-{reset} {grey}%s({{reset}\n", prefix, render.RenderType(m.wantType, opts.RenderOpts())))
		builder.WriteString(opts.RenderOpts().Printf("{green}+{reset} {grey}%s({{reset}\n", render.RenderType(m.gotType, opts.RenderOpts())))
	} else {
		builder.WriteString(opts.RenderOpts().Printf("%s{grey}{{reset}\n", prefix))
	}
	for _, e := range entries {
		content := strings.ReplaceAll(e.content, "\n", "\n  ")
		builder.WriteString("  " + content + "\n")
	}
	if opts.RenderOpts().ShowTypes {
		builder.WriteString(opts.RenderOpts().Print("  {grey}}){reset}"))
	} else {
		builder.WriteString(opts.RenderOpts().Print("  {grey}}{reset}"))
	}
	return builder.String()
}
