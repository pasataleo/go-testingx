# v0.2.0

## FEATURES

- `CaptureAll` wraps several values in a `*Value` chain without calling a function, asserted on from the last value to the first as `Call` does. It is what `Capture` did before it became generic

## IMPROVEMENTS

- `Capture` is generic and takes a single value, which keeps its static type: a nil `error` is still an error to the error assertions, and `Capture()` on the result returns `T` rather than `interface{}`, as with `CallAs`

## BUG FIXES

- The error assertions accept a nil error passed to `CaptureAll`. Going through `...any` turns a nil `error` into a nil `interface{}`, so `NoError()` failed with `expected error type, got interface {}`, and `ErrorContains`, `ErrorCode` and the rest reported that instead of `got nil`. A nil value of any other type is still rejected as not an error

<!--
## SECURITY
Vulnerabilities or security-related changes addressed in this release.
-->

<!--
## DEPRECATIONS
Functionality that will be removed in a future release.
-->

## BREAKING CHANGES

- `Capture` takes exactly one value: `Capture[T any](t testing.TB, value T) *Value[T]`. Capturing several values at once moved to `CaptureAll`
- `Capture(t, nil)` no longer compiles, since there is nothing to infer `T` from

## UPGRADE NOTES

- `testingx.Capture(t, a, b)` becomes `testingx.CaptureAll(t, a, b)`
- `testingx.Capture(t, nil)` becomes `testingx.Capture[error](t, nil)`, or whichever type the value stands in for

# v0.1.0

## FEATURES

- `Call`, `CallAs`, and `Capture` for executing functions and chaining assertions on return values
- Fluent assertion chain with `Fatal` / `NonFatal` switching and `Capture` for extracting typed values
- Equality assertions: `Equal`, `NotEqual`
- Boolean assertions: `True`, `False`
- Nil assertions: `Nil`, `NotNil`
- Length assertions: `Len`, `Empty`, `NotEmpty`
- Containment assertions: `Contains`, `NotContains` (strings, slices, maps)
- Error assertions: `NoError`, `Error`, `MatchesError`, `MatchesErrorf`, `ErrorCode`, `ErrorContains`, `ErrorContainsf`, `HasError`, `HasErrorf`
- Custom validation via `Validate(fn)`
- Panic assertions: `Panics`, `NotPanics`, `PanicsAs`, `NotPanicsAs`
- `render` sub-package for human-readable value representations with colour support, type annotations, and custom renderers
- `diff` sub-package for structured deep diffs between arbitrary Go values
- `mocks` sub-package providing `mocks.T` for capturing assertion failure messages in tests
