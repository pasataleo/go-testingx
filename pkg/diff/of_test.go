package diff

import (
	"reflect"
	"testing"

	"github.com/pasataleo/go-testingx/pkg/render"
)

type testPerson struct {
	Name string
	Age  int
}

type testContact struct {
	Name  string
	Email string
}

type testEmpty struct{}

// testVersion implements Diffable[testVersion].
type testVersion struct {
	Major, Minor int
}

func (v *testVersion) Diff(want *testVersion, opts *Opts) Result {
	return ValueResult(reflect.ValueOf(v), reflect.ValueOf(want))
}

// testScore is used with a registered Differ.
type testScore struct {
	Value int
}

type testScoreDiffer struct{}

func (d testScoreDiffer) Diff(got, want *testScore, opts *Opts) Result {
	return ValueResult(reflect.ValueOf(got), reflect.ValueOf(want))
}

type testScorer interface {
	Score() string
}

type testScorerDiffer struct{}

func (d testScorerDiffer) Diff(got, want testScorer, opts *Opts) Result {
	return ValueResult(reflect.ValueOf(got.Score()), reflect.ValueOf(want.Score()))
}

func TestOf(t *testing.T) {
	tcs := map[string]struct {
		got      reflect.Value
		want     reflect.Value
		diffOpts []OptsFn

		status     Status
		renderGot  *string
		renderWant *string
		renderDiff string
	}{
		"missing int": {
			want:       reflect.ValueOf(42),
			status:     StatusMissing,
			renderWant: ptr("int(42)"),
			renderDiff: "- int(42)",
		},
		"missing string": {
			want:       reflect.ValueOf("hello"),
			status:     StatusMissing,
			renderWant: ptr(`string("hello")`),
			renderDiff: `- string("hello")`,
		},
		"extra int": {
			got:        reflect.ValueOf(42),
			status:     StatusExtra,
			renderGot:  ptr("int(42)"),
			renderDiff: "+ int(42)",
		},
		"extra string": {
			got:        reflect.ValueOf("hello"),
			status:     StatusExtra,
			renderGot:  ptr(`string("hello")`),
			renderDiff: `+ string("hello")`,
		},
		"changed int": {
			got:        reflect.ValueOf(1),
			want:       reflect.ValueOf(2),
			status:     StatusChanged,
			renderGot:  ptr("int(1)"),
			renderWant: ptr("int(2)"),
			renderDiff: "- int(2)\n+ int(1)",
		},
		"changed string": {
			got:        reflect.ValueOf("got"),
			want:       reflect.ValueOf("want"),
			status:     StatusChanged,
			renderGot:  ptr(`string("got")`),
			renderWant: ptr(`string("want")`),
			renderDiff: "- string(\"want\")\n+ string(\"got\")",
		},
		"unchanged int": {
			got:        reflect.ValueOf(42),
			want:       reflect.ValueOf(42),
			status:     StatusUnchanged,
			renderGot:  ptr("int(42)"),
			renderWant: ptr("int(42)"),
			renderDiff: "  int(42)",
		},
		"unchanged string": {
			got:        reflect.ValueOf("same"),
			want:       reflect.ValueOf("same"),
			status:     StatusUnchanged,
			renderGot:  ptr(`string("same")`),
			renderWant: ptr(`string("same")`),
			renderDiff: `  string("same")`,
		},

		// Maps
		"map unchanged": {
			got:    reflect.ValueOf(map[string]int{"a": 1}),
			want:   reflect.ValueOf(map[string]int{"a": 1}),
			status: StatusUnchanged,
			renderGot: ptr("map[string]int({\n" +
				"  string(\"a\"): int(1),\n" +
				"})"),
			renderWant: ptr("map[string]int({\n" +
				"  string(\"a\"): int(1),\n" +
				"})"),
			renderDiff: "  map[string]int({...})",
		},
		"map changed value": {
			got:    reflect.ValueOf(map[string]int{"a": 2}),
			want:   reflect.ValueOf(map[string]int{"a": 1}),
			status: StatusChanged,
			renderGot: ptr("map[string]int({\n" +
				"  string(\"a\"): int(2),\n" +
				"})"),
			renderWant: ptr("map[string]int({\n" +
				"  string(\"a\"): int(1),\n" +
				"})"),
			renderDiff: "~ map[string]int({\n" +
				"  - string(\"a\"): int(1),\n" +
				"  + string(\"a\"): int(2),\n" +
				"  })",
		},
		"map extra key": {
			got:    reflect.ValueOf(map[string]int{"a": 1, "b": 2}),
			want:   reflect.ValueOf(map[string]int{"a": 1}),
			status: StatusChanged,
			renderGot: ptr("map[string]int({\n" +
				"  string(\"a\"): int(1),\n" +
				"  string(\"b\"): int(2),\n" +
				"})"),
			renderWant: ptr("map[string]int({\n" +
				"  string(\"a\"): int(1),\n" +
				"})"),
			renderDiff: "~ map[string]int({\n" +
				"  + string(\"b\"): int(2),\n" +
				"  })",
		},
		"map missing key": {
			got:    reflect.ValueOf(map[string]int{"a": 1}),
			want:   reflect.ValueOf(map[string]int{"a": 1, "b": 2}),
			status: StatusChanged,
			renderGot: ptr("map[string]int({\n" +
				"  string(\"a\"): int(1),\n" +
				"})"),
			renderWant: ptr("map[string]int({\n" +
				"  string(\"a\"): int(1),\n" +
				"  string(\"b\"): int(2),\n" +
				"})"),
			renderDiff: "~ map[string]int({\n" +
				"  - string(\"b\"): int(2),\n" +
				"  })",
		},

		"map different types overlapping keys": {
			got:    reflect.ValueOf(map[string]string{"a": "hello"}),
			want:   reflect.ValueOf(map[string]int{"a": 1}),
			status: StatusChanged,
			renderGot: ptr("map[string]string({\n" +
				"  string(\"a\"): string(\"hello\"),\n" +
				"})"),
			renderWant: ptr("map[string]int({\n" +
				"  string(\"a\"): int(1),\n" +
				"})"),
			renderDiff: "- map[string]int({\n" +
				"+ map[string]string({\n" +
				"  - string(\"a\"): int(1),\n" +
				"  + string(\"a\"): string(\"hello\"),\n" +
				"  })",
		},

		// Nil maps
		"nil map vs empty map": {
			got:        reflect.ValueOf(map[string]int(nil)),
			want:       reflect.ValueOf(map[string]int{}),
			status:     StatusChanged,
			renderGot:  ptr("nil(map[string]int)"),
			renderWant: ptr("map[string]int({})"),
			renderDiff: "- map[string]int({})\n" +
				"+ nil(map[string]int)",
		},
		"empty map vs nil map": {
			got:        reflect.ValueOf(map[string]int{}),
			want:       reflect.ValueOf(map[string]int(nil)),
			status:     StatusChanged,
			renderGot:  ptr("map[string]int({})"),
			renderWant: ptr("nil(map[string]int)"),
			renderDiff: "- nil(map[string]int)\n" +
				"+ map[string]int({})",
		},
		"nil map vs nil map": {
			got:        reflect.ValueOf(map[string]int(nil)),
			want:       reflect.ValueOf(map[string]int(nil)),
			status:     StatusUnchanged,
			renderGot:  ptr("nil(map[string]int)"),
			renderWant: ptr("nil(map[string]int)"),
			renderDiff: "  nil(map[string]int)",
		},

		// Nested maps
		"nested map changed inner value": {
			got: reflect.ValueOf(map[string]map[string]int{
				"outer": {"a": 2},
			}),
			want: reflect.ValueOf(map[string]map[string]int{
				"outer": {"a": 1},
			}),
			status: StatusChanged,
			renderGot: ptr("map[string]map[string]int({\n" +
				"  string(\"outer\"): map[string]int({\n" +
				"    string(\"a\"): int(2),\n" +
				"  }),\n" +
				"})"),
			renderWant: ptr("map[string]map[string]int({\n" +
				"  string(\"outer\"): map[string]int({\n" +
				"    string(\"a\"): int(1),\n" +
				"  }),\n" +
				"})"),
			renderDiff: "~ map[string]map[string]int({\n" +
				"  ~ string(\"outer\"): map[string]int({\n" +
				"    - string(\"a\"): int(1),\n" +
				"    + string(\"a\"): int(2),\n" +
				"    }),\n" +
				"  })",
		},

		// ShowUnchanged
		"map show unchanged": {
			got:      reflect.ValueOf(map[string]int{"a": 1, "b": 2}),
			want:     reflect.ValueOf(map[string]int{"a": 1, "b": 3}),
			diffOpts: []OptsFn{WithShowUnchanged(true)},
			status:   StatusChanged,
			renderGot: ptr("map[string]int({\n" +
				"  string(\"a\"): int(1),\n" +
				"  string(\"b\"): int(2),\n" +
				"})"),
			renderWant: ptr("map[string]int({\n" +
				"  string(\"a\"): int(1),\n" +
				"  string(\"b\"): int(3),\n" +
				"})"),
			renderDiff: "~ map[string]int({\n" +
				"    string(\"a\"): int(1),\n" +
				"  - string(\"b\"): int(3),\n" +
				"  + string(\"b\"): int(2),\n" +
				"  })",
		},

		// Pointers
		"pointer unchanged": {
			got:        reflect.ValueOf(ptr(42)),
			want:       reflect.ValueOf(ptr(42)),
			status:     StatusUnchanged,
			renderGot:  ptr("&int(42)"),
			renderWant: ptr("&int(42)"),
			renderDiff: "  &int(42)",
		},
		"pointer changed": {
			got:        reflect.ValueOf(ptr(1)),
			want:       reflect.ValueOf(ptr(2)),
			status:     StatusChanged,
			renderGot:  ptr("&int(1)"),
			renderWant: ptr("&int(2)"),
			renderDiff: "- &int(2)\n+ &int(1)",
		},
		"pointer nil got": {
			got:        reflect.ValueOf((*int)(nil)),
			want:       reflect.ValueOf(ptr(42)),
			status:     StatusChanged,
			renderGot:  ptr("nil(int)"),
			renderWant: ptr("&int(42)"),
			renderDiff: "- &int(42)\n+ nil(int)",
		},
		"pointer nil want": {
			got:        reflect.ValueOf(ptr(42)),
			want:       reflect.ValueOf((*int)(nil)),
			status:     StatusChanged,
			renderGot:  ptr("&int(42)"),
			renderWant: ptr("nil(int)"),
			renderDiff: "- nil(int)\n+ &int(42)",
		},
		"pointer both nil": {
			got:        reflect.ValueOf((*int)(nil)),
			want:       reflect.ValueOf((*int)(nil)),
			status:     StatusUnchanged,
			renderGot:  ptr("nil(int)"),
			renderWant: ptr("nil(int)"),
			renderDiff: "  nil(int)",
		},
		"pointer to map changed": {
			got:    reflect.ValueOf(&map[string]int{"a": 2}),
			want:   reflect.ValueOf(&map[string]int{"a": 1}),
			status: StatusChanged,
			renderGot: ptr("&map[string]int({\n" +
				"  string(\"a\"): int(2),\n" +
				"})"),
			renderWant: ptr("&map[string]int({\n" +
				"  string(\"a\"): int(1),\n" +
				"})"),
			renderDiff: "~ &map[string]int({\n" +
				"  - string(\"a\"): int(1),\n" +
				"  + string(\"a\"): int(2),\n" +
				"  })",
		},
		"pointer nested in map": {
			got: reflect.ValueOf(map[string]*int{
				"a": ptr(2),
				"b": ptr(3),
			}),
			want: reflect.ValueOf(map[string]*int{
				"a": ptr(2),
				"c": ptr(3),
			}),
			status: StatusChanged,
			renderGot: ptr("map[string]*int({\n" +
				"  string(\"a\"): &int(2),\n" +
				"  string(\"b\"): &int(3),\n" +
				"})"),
			renderWant: ptr("map[string]*int({\n" +
				"  string(\"a\"): &int(2),\n" +
				"  string(\"c\"): &int(3),\n" +
				"})"),
			renderDiff: "~ map[string]*int({\n" +
				"  + string(\"b\"): &int(3),\n" +
				"  - string(\"c\"): &int(3),\n" +
				"  })",
		},

		// Interfaces
		"interface unchanged": {
			got:        iface(42),
			want:       iface(42),
			status:     StatusUnchanged,
			renderGot:  ptr("interface {}(int(42))"),
			renderWant: ptr("interface {}(int(42))"),
			renderDiff: "  interface {}(int(42))",
		},
		"interface changed same type": {
			got:        iface(1),
			want:       iface(2),
			status:     StatusChanged,
			renderGot:  ptr("interface {}(int(1))"),
			renderWant: ptr("interface {}(int(2))"),
			renderDiff: "- interface {}(int(2))\n+ interface {}(int(1))",
		},
		"interface changed different type": {
			got:        iface(42),
			want:       iface("hello"),
			status:     StatusChanged,
			renderGot:  ptr("interface {}(int(42))"),
			renderWant: ptr(`interface {}(string("hello"))`),
			renderDiff: "- interface {}(string(\"hello\"))\n+ interface {}(int(42))",
		},
		"interface nil got": {
			got:        ifaceNil(),
			want:       iface(42),
			status:     StatusChanged,
			renderGot:  ptr("nil(interface {})"),
			renderWant: ptr("interface {}(int(42))"),
			renderDiff: "- interface {}(int(42))\n+ nil(interface {})",
		},
		"interface nil want": {
			got:        iface(42),
			want:       ifaceNil(),
			status:     StatusChanged,
			renderGot:  ptr("interface {}(int(42))"),
			renderWant: ptr("nil(interface {})"),
			renderDiff: "- nil(interface {})\n+ interface {}(int(42))",
		},
		"interface both nil": {
			got:        ifaceNil(),
			want:       ifaceNil(),
			status:     StatusUnchanged,
			renderGot:  ptr("nil(interface {})"),
			renderWant: ptr("nil(interface {})"),
			renderDiff: "  nil(interface {})",
		},
		"interface wrapping pointer": {
			got:        iface(ptr(1)),
			want:       iface(ptr(2)),
			status:     StatusChanged,
			renderGot:  ptr("interface {}(&int(1))"),
			renderWant: ptr("interface {}(&int(2))"),
			renderDiff: "- interface {}(&int(2))\n+ interface {}(&int(1))",
		},

		// Structs
		"struct unchanged": {
			got:    reflect.ValueOf(testPerson{Name: "Alice", Age: 30}),
			want:   reflect.ValueOf(testPerson{Name: "Alice", Age: 30}),
			status: StatusUnchanged,
			renderGot: ptr("diff.testPerson({\n" +
				"  Name: string(\"Alice\"),\n" +
				"  Age : int(30),\n" +
				"})"),
			renderWant: ptr("diff.testPerson({\n" +
				"  Name: string(\"Alice\"),\n" +
				"  Age : int(30),\n" +
				"})"),
			renderDiff: "  diff.testPerson({...})",
		},
		"struct changed field": {
			got:    reflect.ValueOf(testPerson{Name: "Alice", Age: 31}),
			want:   reflect.ValueOf(testPerson{Name: "Alice", Age: 30}),
			status: StatusChanged,
			renderGot: ptr("diff.testPerson({\n" +
				"  Name: string(\"Alice\"),\n" +
				"  Age : int(31),\n" +
				"})"),
			renderWant: ptr("diff.testPerson({\n" +
				"  Name: string(\"Alice\"),\n" +
				"  Age : int(30),\n" +
				"})"),
			renderDiff: "~ diff.testPerson({\n" +
				"  - Age: int(30),\n" +
				"  + Age: int(31),\n" +
				"  })",
		},
		"struct all fields changed": {
			got:    reflect.ValueOf(testPerson{Name: "Bob", Age: 25}),
			want:   reflect.ValueOf(testPerson{Name: "Alice", Age: 30}),
			status: StatusChanged,
			renderGot: ptr("diff.testPerson({\n" +
				"  Name: string(\"Bob\"),\n" +
				"  Age : int(25),\n" +
				"})"),
			renderWant: ptr("diff.testPerson({\n" +
				"  Name: string(\"Alice\"),\n" +
				"  Age : int(30),\n" +
				"})"),
			renderDiff: "~ diff.testPerson({\n" +
				"  - Name: string(\"Alice\"),\n" +
				"  + Name: string(\"Bob\"),\n" +
				"  - Age: int(30),\n" +
				"  + Age: int(25),\n" +
				"  })",
		},
		"struct show unchanged": {
			got:      reflect.ValueOf(testPerson{Name: "Alice", Age: 31}),
			want:     reflect.ValueOf(testPerson{Name: "Alice", Age: 30}),
			diffOpts: []OptsFn{WithShowUnchanged(true)},
			status:   StatusChanged,
			renderGot: ptr("diff.testPerson({\n" +
				"  Name: string(\"Alice\"),\n" +
				"  Age : int(31),\n" +
				"})"),
			renderWant: ptr("diff.testPerson({\n" +
				"  Name: string(\"Alice\"),\n" +
				"  Age : int(30),\n" +
				"})"),
			renderDiff: "~ diff.testPerson({\n" +
				"    Name: string(\"Alice\"),\n" +
				"  - Age: int(30),\n" +
				"  + Age: int(31),\n" +
				"  })",
		},
		"empty struct unchanged": {
			got:        reflect.ValueOf(testEmpty{}),
			want:       reflect.ValueOf(testEmpty{}),
			status:     StatusUnchanged,
			renderGot:  ptr("diff.testEmpty({})"),
			renderWant: ptr("diff.testEmpty({})"),
			renderDiff: "  diff.testEmpty({...})",
		},
		"struct with nested map changed": {
			got: reflect.ValueOf(struct{ M map[string]int }{
				M: map[string]int{"a": 2},
			}),
			want: reflect.ValueOf(struct{ M map[string]int }{
				M: map[string]int{"a": 1},
			}),
			status: StatusChanged,
			renderGot: ptr("struct { M map[string]int }({\n" +
				"  M: map[string]int({\n" +
				"    string(\"a\"): int(2),\n" +
				"  }),\n" +
				"})"),
			renderWant: ptr("struct { M map[string]int }({\n" +
				"  M: map[string]int({\n" +
				"    string(\"a\"): int(1),\n" +
				"  }),\n" +
				"})"),
			renderDiff: "~ struct { M map[string]int }({\n" +
				"  ~ M: map[string]int({\n" +
				"    - string(\"a\"): int(1),\n" +
				"    + string(\"a\"): int(2),\n" +
				"    }),\n" +
				"  })",
		},
		"struct different types overlapping fields": {
			got:    reflect.ValueOf(testPerson{Name: "Alice", Age: 30}),
			want:   reflect.ValueOf(testContact{Name: "Bob", Email: "bob@example.com"}),
			status: StatusChanged,
			renderGot: ptr("diff.testPerson({\n" +
				"  Name: string(\"Alice\"),\n" +
				"  Age : int(30),\n" +
				"})"),
			renderWant: ptr("diff.testContact({\n" +
				"  Name : string(\"Bob\"),\n" +
				"  Email: string(\"bob@example.com\"),\n" +
				"})"),
			renderDiff: "- diff.testContact({\n" +
				"+ diff.testPerson({\n" +
				"  - Name: string(\"Bob\"),\n" +
				"  + Name: string(\"Alice\"),\n" +
				"  + Age: int(30),\n" +
				"  - Email: string(\"bob@example.com\"),\n" +
				"  })",
		},
		"pointer to struct changed": {
			got:    reflect.ValueOf(&testPerson{Name: "Bob", Age: 25}),
			want:   reflect.ValueOf(&testPerson{Name: "Alice", Age: 30}),
			status: StatusChanged,
			renderGot: ptr("&diff.testPerson({\n" +
				"  Name: string(\"Bob\"),\n" +
				"  Age : int(25),\n" +
				"})"),
			renderWant: ptr("&diff.testPerson({\n" +
				"  Name: string(\"Alice\"),\n" +
				"  Age : int(30),\n" +
				"})"),
			renderDiff: "~ &diff.testPerson({\n" +
				"  - Name: string(\"Alice\"),\n" +
				"  + Name: string(\"Bob\"),\n" +
				"  - Age: int(30),\n" +
				"  + Age: int(25),\n" +
				"  })",
		},

		// Slices
		"slice unchanged": {
			got:    reflect.ValueOf([]int{1, 2, 3}),
			want:   reflect.ValueOf([]int{1, 2, 3}),
			status: StatusUnchanged,
			renderGot: ptr("[]int([\n" +
				"  [0]: int(1),\n" +
				"  [1]: int(2),\n" +
				"  [2]: int(3),\n" +
				"])"),
			renderWant: ptr("[]int([\n" +
				"  [0]: int(1),\n" +
				"  [1]: int(2),\n" +
				"  [2]: int(3),\n" +
				"])"),
			renderDiff: "  []int([...])",
		},
		"slice changed element": {
			got:    reflect.ValueOf([]int{1, 9, 3}),
			want:   reflect.ValueOf([]int{1, 2, 3}),
			status: StatusChanged,
			renderGot: ptr("[]int([\n" +
				"  [0]: int(1),\n" +
				"  [1]: int(9),\n" +
				"  [2]: int(3),\n" +
				"])"),
			renderWant: ptr("[]int([\n" +
				"  [0]: int(1),\n" +
				"  [1]: int(2),\n" +
				"  [2]: int(3),\n" +
				"])"),
			renderDiff: "~ []int([\n" +
				"  - [1]: int(2),\n" +
				"  + [1]: int(9),\n" +
				"  ])",
		},
		"slice extra elements": {
			got:    reflect.ValueOf([]int{1, 2, 3}),
			want:   reflect.ValueOf([]int{1}),
			status: StatusChanged,
			renderGot: ptr("[]int([\n" +
				"  [0]: int(1),\n" +
				"  [1]: int(2),\n" +
				"  [2]: int(3),\n" +
				"])"),
			renderWant: ptr("[]int([\n" +
				"  [0]: int(1),\n" +
				"])"),
			renderDiff: "~ []int([\n" +
				"  + [1]: int(2),\n" +
				"  + [2]: int(3),\n" +
				"  ])",
		},
		"slice missing elements": {
			got:    reflect.ValueOf([]int{1}),
			want:   reflect.ValueOf([]int{1, 2, 3}),
			status: StatusChanged,
			renderGot: ptr("[]int([\n" +
				"  [0]: int(1),\n" +
				"])"),
			renderWant: ptr("[]int([\n" +
				"  [0]: int(1),\n" +
				"  [1]: int(2),\n" +
				"  [2]: int(3),\n" +
				"])"),
			renderDiff: "~ []int([\n" +
				"  - [1]: int(2),\n" +
				"  - [2]: int(3),\n" +
				"  ])",
		},
		"empty slice unchanged": {
			got:        reflect.ValueOf([]int{}),
			want:       reflect.ValueOf([]int{}),
			status:     StatusUnchanged,
			renderGot:  ptr("[]int([])"),
			renderWant: ptr("[]int([])"),
			renderDiff: "  []int([...])",
		},
		"slice show unchanged": {
			got:      reflect.ValueOf([]int{1, 9}),
			want:     reflect.ValueOf([]int{1, 2}),
			diffOpts: []OptsFn{WithShowUnchanged(true)},
			status:   StatusChanged,
			renderGot: ptr("[]int([\n" +
				"  [0]: int(1),\n" +
				"  [1]: int(9),\n" +
				"])"),
			renderWant: ptr("[]int([\n" +
				"  [0]: int(1),\n" +
				"  [1]: int(2),\n" +
				"])"),
			renderDiff: "~ []int([\n" +
				"    [0]: int(1),\n" +
				"  - [1]: int(2),\n" +
				"  + [1]: int(9),\n" +
				"  ])",
		},
		"slice with nested map changed": {
			got: reflect.ValueOf([]map[string]int{
				{"a": 2},
			}),
			want: reflect.ValueOf([]map[string]int{
				{"a": 1},
			}),
			status: StatusChanged,
			renderGot: ptr("[]map[string]int([\n" +
				"  [0]: map[string]int({\n" +
				"         string(\"a\"): int(2),\n" +
				"       }),\n" +
				"])"),
			renderWant: ptr("[]map[string]int([\n" +
				"  [0]: map[string]int({\n" +
				"         string(\"a\"): int(1),\n" +
				"       }),\n" +
				"])"),
			renderDiff: "~ []map[string]int([\n" +
				"  ~ [0]: map[string]int({\n" +
				"    - string(\"a\"): int(1),\n" +
				"    + string(\"a\"): int(2),\n" +
				"    }),\n" +
				"  ])",
		},
		"slice different types": {
			got:    reflect.ValueOf([]string{"hello"}),
			want:   reflect.ValueOf([]int{1}),
			status: StatusChanged,
			renderGot: ptr("[]string([\n" +
				"  [0]: string(\"hello\"),\n" +
				"])"),
			renderWant: ptr("[]int([\n" +
				"  [0]: int(1),\n" +
				"])"),
			renderDiff: "- []int([\n" +
				"+ []string([\n" +
				"  - [0]: int(1),\n" +
				"  + [0]: string(\"hello\"),\n" +
				"  ])",
		},

		// Arrays
		"array unchanged": {
			got:    reflect.ValueOf([2]int{1, 2}),
			want:   reflect.ValueOf([2]int{1, 2}),
			status: StatusUnchanged,
			renderGot: ptr("[2]int([\n" +
				"  [0]: int(1),\n" +
				"  [1]: int(2),\n" +
				"])"),
			renderWant: ptr("[2]int([\n" +
				"  [0]: int(1),\n" +
				"  [1]: int(2),\n" +
				"])"),
			renderDiff: "  [2]int([...])",
		},
		"array changed element": {
			got:    reflect.ValueOf([2]int{1, 9}),
			want:   reflect.ValueOf([2]int{1, 2}),
			status: StatusChanged,
			renderGot: ptr("[2]int([\n" +
				"  [0]: int(1),\n" +
				"  [1]: int(9),\n" +
				"])"),
			renderWant: ptr("[2]int([\n" +
				"  [0]: int(1),\n" +
				"  [1]: int(2),\n" +
				"])"),
			renderDiff: "~ [2]int([\n" +
				"  - [1]: int(2),\n" +
				"  + [1]: int(9),\n" +
				"  ])",
		},

		// Diffable interface
		"diffable unchanged": {
			got:    reflect.ValueOf(&testVersion{Major: 1, Minor: 0}),
			want:   reflect.ValueOf(&testVersion{Major: 1, Minor: 0}),
			status: StatusUnchanged,
			renderGot: ptr("&diff.testVersion({\n" +
				"  Major: int(1),\n" +
				"  Minor: int(0),\n" +
				"})"),
			renderWant: ptr("&diff.testVersion({\n" +
				"  Major: int(1),\n" +
				"  Minor: int(0),\n" +
				"})"),
			renderDiff: "  &diff.testVersion({\n" +
				"  Major: int(1),\n" +
				"  Minor: int(0),\n" +
				"})",
		},
		"diffable changed": {
			got:    reflect.ValueOf(&testVersion{Major: 2, Minor: 0}),
			want:   reflect.ValueOf(&testVersion{Major: 1, Minor: 0}),
			status: StatusChanged,
			renderGot: ptr("&diff.testVersion({\n" +
				"  Major: int(2),\n" +
				"  Minor: int(0),\n" +
				"})"),
			renderWant: ptr("&diff.testVersion({\n" +
				"  Major: int(1),\n" +
				"  Minor: int(0),\n" +
				"})"),
			renderDiff: "- &diff.testVersion({\n" +
				"  - Major: int(1),\n" +
				"  - Minor: int(0),\n" +
				"})\n" +
				"+ &diff.testVersion({\n" +
				"  + Major: int(2),\n" +
				"  + Minor: int(0),\n" +
				"})",
		},

		"diffable vs different type": {
			got:    reflect.ValueOf(&testVersion{Major: 1, Minor: 0}),
			want:   reflect.ValueOf("hello"),
			status: StatusChanged,
			renderGot: ptr("&diff.testVersion({\n" +
				"  Major: int(1),\n" +
				"  Minor: int(0),\n" +
				"})"),
			renderWant: ptr(`string("hello")`),
			renderDiff: "- string(\"hello\")\n" +
				"+ &diff.testVersion({\n" +
				"  + Major: int(1),\n" +
				"  + Minor: int(0),\n" +
				"})",
		},

		"diffable vs nil": {
			got:    reflect.ValueOf(&testVersion{Major: 1, Minor: 0}),
			want:   reflect.ValueOf((*testVersion)(nil)),
			status: StatusChanged,
			renderGot: ptr("&diff.testVersion({\n" +
				"  Major: int(1),\n" +
				"  Minor: int(0),\n" +
				"})"),
			renderWant: ptr("nil(diff.testVersion)"),
			renderDiff: "- nil(diff.testVersion)\n" +
				"+ &diff.testVersion({\n" +
				"  + Major: int(1),\n" +
				"  + Minor: int(0),\n" +
				"})",
		},

		// Registered Differ
		"differ unchanged": {
			got:      reflect.ValueOf(&testScore{Value: 100}),
			want:     reflect.ValueOf(&testScore{Value: 100}),
			diffOpts: []OptsFn{WithDiffer[*testScore](testScoreDiffer{})},
			status:   StatusUnchanged,
			renderGot: ptr("&diff.testScore({\n" +
				"  Value: int(100),\n" +
				"})"),
			renderWant: ptr("&diff.testScore({\n" +
				"  Value: int(100),\n" +
				"})"),
			renderDiff: "  &diff.testScore({\n" +
				"  Value: int(100),\n" +
				"})",
		},
		"differ changed": {
			got:      reflect.ValueOf(&testScore{Value: 200}),
			want:     reflect.ValueOf(&testScore{Value: 100}),
			diffOpts: []OptsFn{WithDiffer[*testScore](testScoreDiffer{})},
			status:   StatusChanged,
			renderGot: ptr("&diff.testScore({\n" +
				"  Value: int(200),\n" +
				"})"),
			renderWant: ptr("&diff.testScore({\n" +
				"  Value: int(100),\n" +
				"})"),
			renderDiff: "- &diff.testScore({\n" +
				"  - Value: int(100),\n" +
				"})\n" +
				"+ &diff.testScore({\n" +
				"  + Value: int(200),\n" +
				"})",
		},

		"differ vs different type": {
			got:      reflect.ValueOf(&testScore{Value: 100}),
			want:     reflect.ValueOf("hello"),
			diffOpts: []OptsFn{WithDiffer[*testScore](testScoreDiffer{})},
			status:   StatusChanged,
			renderGot: ptr("&diff.testScore({\n" +
				"  Value: int(100),\n" +
				"})"),
			renderWant: ptr(`string("hello")`),
			renderDiff: "- string(\"hello\")\n" +
				"+ &diff.testScore({\n" +
				"  + Value: int(100),\n" +
				"})",
		},

		"differ vs nil": {
			got:    reflect.ValueOf(&testScore{Value: 100}),
			want:   reflect.ValueOf((*testScore)(nil)),
			status: StatusChanged,
			renderGot: ptr("&diff.testScore({\n" +
				"  Value: int(100),\n" +
				"})"),
			renderWant: ptr("nil(diff.testScore)"),
			renderDiff: "- nil(diff.testScore)\n" +
				"+ &diff.testScore({\n" +
				"  + Value: int(100),\n" +
				"})",
		},

		"nil vs differ": {
			want:   reflect.ValueOf(&testScore{Value: 100}),
			got:    reflect.ValueOf((*testScore)(nil)),
			status: StatusChanged,
			renderWant: ptr("&diff.testScore({\n" +
				"  Value: int(100),\n" +
				"})"),
			renderGot: ptr("nil(diff.testScore)"),
			renderDiff: "- &diff.testScore({\n" +
				"  - Value: int(100),\n" +
				"})\n" +
				"+ nil(diff.testScore)",
		},

		"differ vs nil interface": {
			got:    reflect.ValueOf(&testScore{Value: 100}),
			want:   ifaceNil(),
			status: StatusChanged,
			renderGot: ptr("&diff.testScore({\n" +
				"  Value: int(100),\n" +
				"})"),
			renderWant: ptr("nil(interface {})"),
			renderDiff: "- nil(interface {})\n" +
				"+ &diff.testScore({\n" +
				"  + Value: int(100),\n" +
				"})",
		},

		"nil interface vs differ": {
			want:   reflect.ValueOf(&testScore{Value: 100}),
			got:    ifaceNil(),
			status: StatusChanged,
			renderWant: ptr("&diff.testScore({\n" +
				"  Value: int(100),\n" +
				"})"),
			renderGot: ptr("nil(interface {})"),
			renderDiff: "- &diff.testScore({\n" +
				"  - Value: int(100),\n" +
				"})\n" +
				"+ nil(interface {})",
		},

		// Kind mismatch
		"map got string want": {
			got:    reflect.ValueOf(map[string]int{"a": 1}),
			want:   reflect.ValueOf("hello"),
			status: StatusChanged,
			renderGot: ptr("map[string]int({\n" +
				"  string(\"a\"): int(1),\n" +
				"})"),
			renderWant: ptr(`string("hello")`),
			renderDiff: "- string(\"hello\")\n" +
				"+ map[string]int({\n" +
				"  + string(\"a\"): int(1),\n" +
				"})",
		},

		"nested nil differ": {
			got: reflect.ValueOf(struct {
				String string
				Scorer testScorer
			}{
				String: "hello",
			}),
			want: reflect.ValueOf(struct {
				String string
				Scorer testScorer
			}{
				String: "hello",
			}),
			status:     StatusUnchanged,
			diffOpts:   []OptsFn{WithDiffer[testScorer](testScorerDiffer{})},
			renderGot:  ptr("struct { String string; Scorer diff.testScorer }({\n  String: string(\"hello\"),\n  Scorer: nil(diff.testScorer),\n})"),
			renderWant: ptr("struct { String string; Scorer diff.testScorer }({\n  String: string(\"hello\"),\n  Scorer: nil(diff.testScorer),\n})"),
			renderDiff: "  struct { String string; Scorer diff.testScorer }({...})",
		},
	}
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			renderOpts := render.NewOpts(t, render.DisableColour())
			diffOpts := NewOpts(t, append(tc.diffOpts, WithRenderOpts(renderOpts))...)
			result := of(tc.got, tc.want, diffOpts)

			if result.Status() != tc.status {
				t.Errorf("status: got %v, want %v", result.Status(), tc.status)
			}

			if tc.renderGot != nil {
				got := result.RenderGot(renderOpts)
				if got != *tc.renderGot {
					t.Errorf("RenderGot:\n\tgot  %q\n\twant %q", got, *tc.renderGot)
				}
			} else {
				assertPanics(t, "RenderGot", func() { result.RenderGot(renderOpts) })
			}

			if tc.renderWant != nil {
				got := result.RenderWant(renderOpts)
				if got != *tc.renderWant {
					t.Errorf("RenderWant:\n\tgot  %q\n\twant %q", got, *tc.renderWant)
				}
			} else {
				assertPanics(t, "RenderWant", func() { result.RenderWant(renderOpts) })
			}

			got := result.RenderDiff(diffOpts)
			if got != tc.renderDiff {
				t.Errorf("RenderDiff:\n\tgot  %q\n\twant %q", got, tc.renderDiff)
			}
		})
	}
}

func assertPanics(t *testing.T, name string, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("%s: expected panic", name)
		}
	}()
	fn()
}

func ptr[T any](v T) *T {
	return &v
}

func iface(v any) reflect.Value {
	type wrapper struct{ V any }
	return reflect.ValueOf(wrapper{V: v}).Field(0)
}

func ifaceNil() reflect.Value {
	type wrapper struct{ V any }
	return reflect.ValueOf(wrapper{}).Field(0)
}
