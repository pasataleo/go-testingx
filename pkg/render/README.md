# render

Package `render` produces human-readable string representations of arbitrary Go values using reflection. It is designed for use in test output, diffs, and diagnostics.

## Usage

```go
opts := render.NewOpts(t)
fmt.Println(render.Render(myValue, opts))
```

## Configuration

All configuration is done via functional options passed to `NewOpts`:

| Option | Default | Description |
|---|---|---|
| `WithTypes(bool)` | `true` | Wrap values with type info, e.g. `int(42)` vs `42` |
| `WithPointers(bool)` | `true` | Show `&`/`*` pointer metadata |
| `WithContents(bool)` | `true` | Expand composites or collapse to `{...}`/`[...]` |
| `WithIndices(bool)` | `true` | Prefix slice/array elements with `[i]:` |
| `WithSkipUnexported(bool)` | `false` | Silently skip unexported struct fields (default: fatal) |
| `WithSkipField(string)` | — | Skip a field by name (`"Age"`) or qualified name (`"Person.Age"`) |
| `WithRenderer[T](Renderer[T])` | — | Register a custom renderer for type `T` |
| `DisableColour()` | — | Disable colour output |

```go
opts := render.NewOpts(t,
    render.WithTypes(false),
    render.WithSkipField("CreatedAt"),
)
```

## Custom rendering

There are two ways to customise how a type is rendered.

### Renderable interface

Types can implement `Renderable` to control their own output:

```go
type UserID string

func (id UserID) Render(opts *render.Opts) string {
    return fmt.Sprintf("user:%s", string(id))
}
```

### Renderer[T]

Register an external renderer for types you don't own, including interface types:

```go
opts := render.NewOpts(t,
    render.WithRenderer[fmt.Stringer](myStringerRenderer{}),
)
```

`RendererFunc[T]` is a function adapter for simple cases:

```go
opts := render.NewOpts(t,
    render.WithRenderer[time.Time](render.RendererFunc[time.Time](
        func(v time.Time, opts *render.Opts) string {
            return v.Format(time.RFC3339)
        },
    )),
)
```

### Precedence

1. Exact-type `Renderer[T]` match
2. Interface `Renderer[T]` match
3. `Renderable` interface
4. Built-in reflection

## Supported types

Strings, bools, all int/uint/float/complex variants, uintptr, pointers, interfaces, structs, maps, slices, and arrays.

Channels and functions are unsupported and will fatal.
