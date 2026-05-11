# go-testingx

A fluent assertion library for Go. Call a function, chain assertions on the return values.

## Installation

```sh
go get github.com/pasataleo/go-testingx
```

## Usage

### Calling functions

`Call` executes a function and returns a `*Value` chain you can assert against. Each return value becomes a link in the chain, with the last return value first.

```go
func TestSomething(t *testing.T) {
    testingx.Call(t, myFunc, arg1, arg2).NoError().Equal("expected")
}
```

`CallAs` does the same but with a typed `*Value[T]`, so `Capture()` returns `T` instead of `interface{}`.

```go
result := testingx.CallAs[MyStruct](t, myFunc).NoError().Capture()
```

`Capture` wraps raw values into a `*Value` chain without calling a function.

```go
testingx.Capture(t, val1, val2).Equal(val2).Equal(val1)
```

### Chaining

Each assertion method returns the next `*Value` in the chain, so for a function returning `(string, error)`:

```go
testingx.Call(t, func() (string, error) {
    return "hello", nil
}).NoError().Equal("hello")
```

The error (last return value) is asserted first, then the string.

### Fatal vs non-fatal

By default, assertion failures call `t.Errorf` (non-fatal). Call `Fatal()` to switch to `t.Fatalf`:

```go
testingx.Call(t, myFunc).Fatal().NoError().Equal("expected")
```

Use `NonFatal()` to switch back.

### Extracting values

`Capture()` extracts the raw value from the current position in the chain:

```go
result := testingx.CallAs[Config](t, loadConfig, path).NoError().Capture()
```

## Assertions

### Equality

- `Equal(want)` - asserts the value equals `want`
- `NotEqual(want)` - asserts the value does not equal `want`

### Booleans

- `True()` - asserts the value is `true`
- `False()` - asserts the value is `false`

### Nil

- `Nil(opts...)` - asserts the value is nil
- `NotNil()` - asserts the value is not nil

### Length

- `Len(n)` - asserts the length equals `n`
- `Empty()` - asserts the length is 0
- `NotEmpty()` - asserts the length is not 0

### Contains

- `Contains(want, opts...)` - asserts the value contains `want` (strings, slices, maps)
- `NotContains(want, opts...)` - asserts the value does not contain `want`

### Errors

- `NoError()` - asserts the value is a nil error
- `Error()` - asserts the value is a non-nil error
- `MatchesError(msg)` - asserts the error message equals `msg`
- `MatchesErrorf(format, args...)` - formatted version of `MatchesError`
- `ErrorCode(code)` - asserts the error has the given `errorsx.Code`
- `ErrorContains(substring)` - asserts the error message contains `substring`
- `ErrorContainsf(format, args...)` - formatted version of `ErrorContains`
- `HasError(msg)` - asserts an aggregated error contains a child with message `msg`
- `HasErrorf(format, args...)` - formatted version of `HasError`

### Custom validation

- `Validate(fn)` - calls `fn(t, value)` for custom validation logic; `fn` must be `func(testing.TB, T)` with no return values

## Panics

- `Panics(t, opts, fn, args...)` - asserts `fn` panics, returns the recovered value as a `*Value`
- `NotPanics(t, opts, fn, args...)` - asserts `fn` does not panic, returns the return values
- `PanicsAs[T](t, opts, fn, args...)` - asserts `fn` panics with a value of type `T`
- `NotPanicsAs[T](t, opts, fn, args...)` - typed version of `NotPanics`

The `opts` parameter is a `*render.Opts` that controls how values are rendered in failure messages. Pass `nil` for defaults.

## Sub-packages

### [render](pkg/render/)

Produces human-readable string representations of arbitrary Go values using reflection. Used internally for failure messages, and can be configured via `render.Opts` to control type annotations, colour, unexported field handling, and custom renderers.

### [diff](pkg/diff/)

Computes structured diffs between arbitrary Go values using reflection. Powers the `Equal` and `NotEqual` assertions, producing output that shows exactly what changed, what's missing, and what's extra.

### [mocks](pkg/mocks/)

Provides `mocks.T`, a `testing.T` wrapper that captures `Error`, `Errorf`, `Fatal`, and `Fatalf` messages instead of failing the test. Useful for testing that your own assertion helpers produce the right failure messages.
