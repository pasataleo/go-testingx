package render

// Renderable is implemented by types that know how to render themselves
// as a human-readable string.
type Renderable interface {
	Render(opts *Opts) string
}
