# v0.1.1

<!--
## FEATURES
New functionality added in this release.
-->

<!--
## IMPROVEMENTS
Enhancements to existing functionality.
-->

## BUG FIXES

- The error assertions accept a nil error passed to `Capture`. Going through `...any` turns a nil `error` into a nil `interface{}`, so `Capture(t, err).NoError()` failed with `expected error type, got interface {}` when `err` was nil, and `ErrorContains`, `ErrorCode` and the rest reported that instead of `got nil`. A nil value of any other type is still rejected as not an error

<!--
## SECURITY
Vulnerabilities or security-related changes addressed in this release.
-->

<!--
## DEPRECATIONS
Functionality that will be removed in a future release.
-->

<!--
## BREAKING CHANGES
Changes that are not backwards compatible and require updates from consumers.
-->

<!--
## UPGRADE NOTES
Steps required when upgrading from a previous version.
-->

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
