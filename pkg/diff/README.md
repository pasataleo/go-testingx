# diff

Package `diff` computes structured diffs between arbitrary Go values using reflection. It is designed for use in test assertions, producing human-readable output that shows exactly what changed, what's missing, and what's extra.

## Usage

```go
opts := diff.NewOpts(t)
result := diff.Of(got, want, opts)
if result.Status() != diff.StatusUnchanged {
    t.Errorf("mismatch:\n%s", result.RenderDiff(opts))
}
```

## Status

Each `Result` carries a `Status` describing the comparison outcome:

| Status | Meaning | Prefix |
|---|---|---|
| `StatusUnchanged` | got == want | `  ` (two spaces) |
| `StatusChanged` | got and want differ | `~` (composite) or `-`/`+` (leaf) |
| `StatusMissing` | in want but not got | `-` |
| `StatusExtra` | in got but not want | `+` |

## Configuration

All configuration is done via functional options passed to `NewOpts`:

| Option | Default | Description |
|---|---|---|
| `WithShowUnchanged(bool)` | `false` | Include unchanged entries in composite diff output |
| `WithRenderOpts(*render.Opts)` | `render.NewOpts(tb)` | Override the render options used for output |
| `WithDiffer[T](Differ[T])` | — | Register a custom differ for type `T` |

```go
opts := diff.NewOpts(t,
    diff.WithShowUnchanged(true),
    diff.WithRenderOpts(render.NewOpts(t, render.DisableColour())),
)
```

## Custom diffing

There are two ways to customise how a type is diffed.

### Diffable interface

Types can implement `Diffable` to control their own diff logic:

```go
type Version struct {
    Major, Minor, Patch int
}

func (v Version) Diff(want Version, opts *diff.Opts) diff.Result {
    // custom comparison logic
}
```

### Differ[T]

Register an external differ for types you don't own, including interface types:

```go
opts := diff.NewOpts(t,
    diff.WithDiffer[time.Time](myTimeDiffer{}),
)
```

### Precedence

1. Exact-type `Differ[T]` match
2. Interface `Differ[T]` match
3. `Diffable` interface
4. Built-in reflection

## Result types

The package provides several result constructors for use in custom `Differ` and `Diffable` implementations:

| Constructor | Description |
|---|---|
| `ValueResult(got, want)` | Leaf value diff (unchanged or changed) |
| `MissingResult(want)` | Value present in want but not got |
| `ExtraResult(got)` | Value present in got but not want |
| `MapResult(type, entries)` | Composite map diff |
| `StructResult(type, fields)` | Composite struct diff |
| `SliceResult(type, entries)` | Composite slice/array diff |
| `PointerResult(got, want, inner)` | Pointer wrapper around an inner diff |
| `InterfaceResult(got, want, inner)` | Interface wrapper around an inner diff |

### CompositeResult

Result types that represent multi-line composite structures (maps, structs, slices) implement `CompositeResult`. This controls how parent diffs render them — composites use `~` prefix with nested output, while leaf types use `- want` / `+ got` pairs.

## Cross-type diffing

Maps and structs support diffing values of different types. When the got and want types differ, the diff header shows both types with `+`/`-` prefixes, and fields or keys are matched by name:

```
+ diff.testPerson({
- diff.testContact({
  - Name: string("Bob"),
  + Name: string("Alice"),
  + Age: int(30),
  - Email: string("bob@example.com"),
  })
```

## Supported types

Maps, structs, slices, arrays, pointers, interfaces, and all leaf types (strings, bools, ints, floats, etc.).

Channels and functions are unsupported and will fatal.

## Planned

- LCS/edit-distance based slice diffing (current implementation is positional only)
