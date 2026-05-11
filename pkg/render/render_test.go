package render

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type testPerson struct {
	Name string
	Age  int
}

type testEmpty struct{}

type testUnexported struct {
	Name     string
	hidden   int
	Exported bool
}

func TestRender(t *testing.T) {
	tcs := map[string]struct {
		input interface{}
		opts  []OptsFn
		want  string
	}{
		"string": {
			input: "hello world",
			want:  `string("hello world")`,
		},
		"string pointer": {
			input: func() *string {
				s := "hello world"
				return &s
			}(),
			want: `&string("hello world")`,
		},
		"nil string": {
			input: (*string)(nil),
			want:  "nil(string)",
		},
		"bool true": {
			input: true,
			want:  "bool(true)",
		},
		"bool false": {
			input: false,
			want:  "bool(false)",
		},
		"int": {
			input: 42,
			want:  "int(42)",
		},
		"int negative": {
			input: -7,
			want:  "int(-7)",
		},
		"int8": {
			input: int8(127),
			want:  "int8(127)",
		},
		"int16": {
			input: int16(32000),
			want:  "int16(32000)",
		},
		"int32": {
			input: int32(100000),
			want:  "int32(100000)",
		},
		"int64": {
			input: int64(9999999999),
			want:  "int64(9999999999)",
		},
		"uint": {
			input: uint(42),
			want:  "uint(42)",
		},
		"uint8": {
			input: uint8(255),
			want:  "uint8(255)",
		},
		"uint16": {
			input: uint16(65535),
			want:  "uint16(65535)",
		},
		"uint32": {
			input: uint32(100000),
			want:  "uint32(100000)",
		},
		"uint64": {
			input: uint64(9999999999),
			want:  "uint64(9999999999)",
		},
		"float32": {
			input: float32(1.5),
			want:  "float32(1.5)",
		},
		"float64": {
			input: float64(2.718281828),
			want:  "float64(2.718281828)",
		},
		"float64 whole": {
			input: float64(42),
			want:  "float64(42)",
		},
		"complex64": {
			input: complex64(1 + 2i),
			want:  "complex64((1+2i))",
		},
		"complex128": {
			input: complex128(3 + 4i),
			want:  "complex128((3+4i))",
		},
		"nil int pointer": {
			input: (*int)(nil),
			want:  "nil(int)",
		},
		"int pointer": {
			input: func() *int {
				v := 42
				return &v
			}(),
			want: "&int(42)",
		},
		"nil bool pointer": {
			input: (*bool)(nil),
			want:  "nil(boolean)",
		},
		"interface": {
			input: func() reflect.Value {
				type wrapper struct{ V any }
				w := wrapper{V: 42}
				return reflect.ValueOf(w).Field(0)
			}(),
			want: "interface {}(int(42))",
		},
		"nil interface": {
			input: func() reflect.Value {
				type wrapper struct{ V error }
				w := wrapper{}
				return reflect.ValueOf(w).Field(0)
			}(),
			want: "nil(error)",
		},
		"mixed interface and pointers": {
			input: func() interface{} {
				i := (**int)(nil)
				var j interface{} = &i
				var k interface{} = &j
				return k
			}(),
			want: "&interface {}(&nil(*int))",
		},

		// ShowTypes=false
		"no types string": {
			input: "hello",
			opts:  []OptsFn{WithTypes(false)},
			want:  `"hello"`,
		},
		"no types bool": {
			input: true,
			opts:  []OptsFn{WithTypes(false)},
			want:  "true",
		},
		"no types int": {
			input: 42,
			opts:  []OptsFn{WithTypes(false)},
			want:  "42",
		},
		"no types float64": {
			input: float64(3.14),
			opts:  []OptsFn{WithTypes(false)},
			want:  "3.14",
		},
		"no types complex128": {
			input: complex128(1 + 2i),
			opts:  []OptsFn{WithTypes(false)},
			want:  "(1+2i)",
		},
		"no types nil pointer": {
			input: (*string)(nil),
			opts:  []OptsFn{WithTypes(false)},
			want:  "nil",
		},
		"no types pointer": {
			input: func() *int {
				v := 42
				return &v
			}(),
			opts: []OptsFn{WithTypes(false)},
			want: "&42",
		},
		"no types nil interface": {
			input: func() reflect.Value {
				type wrapper struct{ V error }
				w := wrapper{}
				return reflect.ValueOf(w).Field(0)
			}(),
			opts: []OptsFn{WithTypes(false)},
			want: "nil",
		},
		"no types interface": {
			input: func() reflect.Value {
				type wrapper struct{ V any }
				w := wrapper{V: 42}
				return reflect.ValueOf(w).Field(0)
			}(),
			opts: []OptsFn{WithTypes(false)},
			want: "42",
		},

		// ShowPointers=false
		"no pointers string pointer": {
			input: func() *string {
				s := "hello"
				return &s
			}(),
			opts: []OptsFn{WithPointers(false)},
			want: `string("hello")`,
		},
		"no pointers int pointer": {
			input: func() *int {
				v := 42
				return &v
			}(),
			opts: []OptsFn{WithPointers(false)},
			want: "int(42)",
		},
		"no pointers nil pointer type": {
			input: (**int)(nil),
			opts:  []OptsFn{WithPointers(false)},
			want:  "nil(int)",
		},

		// ShowTypes=false, ShowPointers=false
		"no types no pointers pointer": {
			input: func() *int {
				v := 42
				return &v
			}(),
			opts: []OptsFn{WithTypes(false), WithPointers(false)},
			want: "42",
		},

		// Maps
		"empty map": {
			input: map[string]int{},
			want:  "map[string]int({})",
		},
		"nil map": {
			input: map[string]int(nil),
			want:  "nil(map[string]int)",
		},
		"map single entry": {
			input: map[string]int{"a": 1},
			want: "map[string]int({\n" +
				"  string(\"a\"): int(1),\n" +
				"})",
		},
		"map multiple entries sorted": {
			input: map[string]int{"b": 2, "a": 1},
			want: "map[string]int({\n" +
				"  string(\"a\"): int(1),\n" +
				"  string(\"b\"): int(2),\n" +
				"})",
		},
		"map nested": {
			input: map[string]map[string]int{
				"outer": {"inner": 42},
			},
			want: "map[string]map[string]int({\n" +
				"  string(\"outer\"): map[string]int({\n" +
				"    string(\"inner\"): int(42),\n" +
				"  }),\n" +
				"})",
		},
		"map multiline key": {
			input: map[testPerson]string{
				{Name: "Alice", Age: 30}: "hello",
			},
			opts: []OptsFn{WithSkipUnexported(true)},
			want: "map[render.testPerson]string({\n" +
				"  render.testPerson({\n" +
				"    Name: string(\"Alice\"),\n" +
				"    Age : int(30),\n" +
				"  }): string(\"hello\"),\n" +
				"})",
		},
		"map aligned keys": {
			input: map[string]int{"a": 1, "longer": 2},
			want: "map[string]int({\n" +
				"  string(\"a\")     : int(1),\n" +
				"  string(\"longer\"): int(2),\n" +
				"})",
		},
		"no types empty map": {
			input: map[string]int{},
			opts:  []OptsFn{WithTypes(false)},
			want:  "{}",
		},
		"no types map": {
			input: map[string]int{"a": 1},
			opts:  []OptsFn{WithTypes(false)},
			want: "{\n" +
				"  \"a\": 1,\n" +
				"}",
		},

		// ShowContents=false
		"no contents map": {
			input: map[string]int{"a": 1, "b": 2},
			opts:  []OptsFn{WithContents(false)},
			want:  "map[string]int({...})",
		},
		"no contents empty map": {
			input: map[string]int{},
			opts:  []OptsFn{WithContents(false)},
			want:  "map[string]int({})",
		},
		"no contents nil map": {
			input: map[string]int(nil),
			opts:  []OptsFn{WithContents(false)},
			want:  "nil(map[string]int)",
		},
		"no contents no types map": {
			input: map[string]int{"a": 1},
			opts:  []OptsFn{WithContents(false), WithTypes(false)},
			want:  "{...}",
		},

		// Structs
		"struct": {
			input: testPerson{Name: "Alice", Age: 30},
			want: "render.testPerson({\n" +
				"  Name: string(\"Alice\"),\n" +
				"  Age : int(30),\n" +
				"})",
		},
		"empty struct": {
			input: testEmpty{},
			want:  "render.testEmpty({})",
		},
		"struct unexported fields skipped": {
			input: testUnexported{Name: "Bob", hidden: 42, Exported: true},
			opts:  []OptsFn{WithSkipUnexported(true)},
			want: "render.testUnexported({\n" +
				"  Name    : string(\"Bob\"),\n" +
				"  Exported: bool(true),\n" +
				"})",
		},
		"struct pointer": {
			input: &testPerson{Name: "Alice", Age: 30},
			want: "&render.testPerson({\n" +
				"  Name: string(\"Alice\"),\n" +
				"  Age : int(30),\n" +
				"})",
		},
		"no types struct": {
			input: testPerson{Name: "Alice", Age: 30},
			opts:  []OptsFn{WithTypes(false)},
			want: "{\n" +
				"  Name: \"Alice\",\n" +
				"  Age : 30,\n" +
				"}",
		},
		"no contents struct": {
			input: testPerson{Name: "Alice", Age: 30},
			opts:  []OptsFn{WithContents(false)},
			want:  "render.testPerson({...})",
		},
		"skip field by name": {
			input: testPerson{Name: "Alice", Age: 30},
			opts:  []OptsFn{WithSkipField("Age")},
			want: "render.testPerson({\n" +
				"  Name: string(\"Alice\"),\n" +
				"})",
		},
		"skip field by struct.field": {
			input: testPerson{Name: "Alice", Age: 30},
			opts:  []OptsFn{WithSkipField("testPerson.Age")},
			want: "render.testPerson({\n" +
				"  Name: string(\"Alice\"),\n" +
				"})",
		},
		"skip field wrong struct": {
			input: testPerson{Name: "Alice", Age: 30},
			opts:  []OptsFn{WithSkipField("testEmpty.Age")},
			want: "render.testPerson({\n" +
				"  Name: string(\"Alice\"),\n" +
				"  Age : int(30),\n" +
				"})",
		},
		// Slices
		"nil slice": {
			input: []string(nil),
			want:  "nil([]string)",
		},
		"empty slice": {
			input: []string{},
			want:  "[]string([])",
		},
		"slice single": {
			input: []int{42},
			want: "[]int([\n" +
				"  [0]: int(42),\n" +
				"])",
		},
		"slice multiple": {
			input: []string{"a", "b", "c"},
			want: "[]string([\n" +
				"  [0]: string(\"a\"),\n" +
				"  [1]: string(\"b\"),\n" +
				"  [2]: string(\"c\"),\n" +
				"])",
		},
		"slice nested": {
			input: [][]int{{1, 2}, {3}},
			want: "[][]int([\n" +
				"  [0]: []int([\n" +
				"         [0]: int(1),\n" +
				"         [1]: int(2),\n" +
				"       ]),\n" +
				"  [1]: []int([\n" +
				"         [0]: int(3),\n" +
				"       ]),\n" +
				"])",
		},
		"no types slice": {
			input: []int{1, 2},
			opts:  []OptsFn{WithTypes(false)},
			want: "[\n" +
				"  [0]: 1,\n" +
				"  [1]: 2,\n" +
				"]",
		},
		"no contents slice": {
			input: []int{1, 2, 3},
			opts:  []OptsFn{WithContents(false)},
			want:  "[]int([...])",
		},
		"no indices slice": {
			input: []int{1, 2},
			opts:  []OptsFn{WithIndices(false)},
			want: "[]int([\n" +
				"  int(1),\n" +
				"  int(2),\n" +
				"])",
		},

		// Arrays
		"array": {
			input: [3]int{1, 2, 3},
			want: "[3]int([\n" +
				"  [0]: int(1),\n" +
				"  [1]: int(2),\n" +
				"  [2]: int(3),\n" +
				"])",
		},
		"empty array": {
			input: [0]int{},
			want:  "[0]int([])",
		},

		// Uintptr
		"uintptr": {
			input: uintptr(57005),
			want:  "uintptr(57005)",
		},

		// Renderable
		"renderable": {
			input: testRenderable{Value: "hello"},
			want:  "custom:hello",
		},
		"renderable via pointer": {
			input: &testRenderable{Value: "hello"},
			want:  "custom:hello",
		},
		"renderable in slice": {
			input: []testRenderable{{Value: "a"}, {Value: "b"}},
			want: "[]render.testRenderable([\n" +
				"  [0]: custom:a,\n" +
				"  [1]: custom:b,\n" +
				"])",
		},
		"renderable in map": {
			input: map[string]testRenderable{
				"key": {Value: "val"},
			},
			want: "map[string]render.testRenderable({\n" +
				"  string(\"key\"): custom:val,\n" +
				"})",
		},

		// Renderer
		"exact type renderer": {
			input: 42,
			opts:  []OptsFn{WithRenderer[int](intRenderer{})},
			want:  "int:42",
		},
		"interface renderer": {
			input: testStringer{Value: "hello"},
			opts:  []OptsFn{WithRenderer[fmt.Stringer](stringerRenderer{})},
			want:  "renderer:stringer:hello",
		},
		"renderer takes precedence over renderable": {
			input: testRenderable{Value: "hello"},
			opts: []OptsFn{WithRenderer[testRenderable](RendererFunc[testRenderable](func(v testRenderable, opts *Opts) string {
				return "override:" + v.Value
			}))},
			want: "override:hello",
		},
		"renderer func": {
			input: "hello",
			opts: []OptsFn{WithRenderer[string](RendererFunc[string](func(value string, opts *Opts) string {
				return "fn:" + value
			}))},
			want: "fn:hello",
		},
		"exact type renderer takes precedence over interface renderer": {
			input: testStringer{Value: "hello"},
			opts: []OptsFn{
				WithRenderer[testStringer](RendererFunc[testStringer](func(v testStringer, opts *Opts) string {
					return "exact:" + v.Value
				})),
				WithRenderer[fmt.Stringer](stringerRenderer{}),
			},
			want: "exact:hello",
		},
	}
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			opts := NewOpts(t, append([]OptsFn{DisableColour()}, tc.opts...)...)

			got := Render(tc.input, opts)
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

type testRenderable struct {
	Value string
}

func (r testRenderable) Render(opts *Opts) string {
	return "custom:" + r.Value
}

type testStringer struct {
	Value string
}

func (s testStringer) String() string {
	return "stringer:" + s.Value
}

type stringerRenderer struct{}

func (stringerRenderer) Render(value fmt.Stringer, opts *Opts) string {
	return "renderer:" + value.String()
}

type intRenderer struct{}

func (intRenderer) Render(value int, opts *Opts) string {
	return "int:" + fmt.Sprintf("%d", value)
}

func TestRenderNilOpts(t *testing.T) {
	// Passing nil opts should not panic and should produce output.
	got := Render("hello", nil)
	if got == "" {
		t.Error("expected non-empty output with nil opts")
	}
	if !strings.Contains(got, `"hello"`) {
		t.Errorf("expected output to contain the value, got %q", got)
	}
}

func TestRenderReflectValue(t *testing.T) {
	v := reflect.ValueOf(42)
	opts := NewOpts(t, DisableColour())
	got := Render(v, opts)
	want := "int(42)"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRenderChannelFatal(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for channel")
		}
	}()
	opts := NewOpts(nil, DisableColour())
	Render(make(chan int), opts)
}

func TestRenderFunctionFatal(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for function")
		}
	}()
	opts := NewOpts(nil, DisableColour())
	Render(func() {}, opts)
}

func TestRenderUnexportedFatal(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for unexported field")
		}
	}()
	opts := NewOpts(nil, DisableColour())
	Render(testUnexported{}, opts)
}
