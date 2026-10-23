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

<!--
## IMPROVEMENTS
Enhancements to existing functionality.
-->

<!--
## BUG FIXES
Issues that have been resolved.
-->

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
