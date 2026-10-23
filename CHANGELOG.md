# v0.2.0

## FEATURES

- `CaptureAll`, for capturing several values at once.

## IMPROVEMENTS

- `Capture` is generic over a single value, so it keeps the value's type.

## BUG FIXES

- The error assertions accept a nil error passed to `CaptureAll`.

<!--
## SECURITY
Vulnerabilities or security-related changes addressed in this release.
-->

<!--
## DEPRECATIONS
Functionality that will be removed in a future release.
-->

## BREAKING CHANGES

- `Capture` takes exactly one value.

## UPGRADE NOTES

- `Capture(t, a, b)` becomes `CaptureAll(t, a, b)`, and `Capture(t, nil)` becomes `Capture[error](t, nil)`.

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
