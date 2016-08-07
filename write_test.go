package hy

import (
	"encoding/json"
	"reflect"
	"testing"
)

type TestWriteStruct struct {
	Name         string
	Int          int
	InlineSlice  []string
	InlineMap    map[string]int
	IgnoredField string             `hy:"-"`              // not output anywhere
	StructFile   StructB            `hy:"a-file"`         // a single file
	StringFile   string             `hy:"a-string-file"`  // a single file
	SliceFile    []string           `hy:"a-slice-file"`   // a single file
	MapFile      map[string]string  `hy:"a-map-file"`     // a single file
	Nested       *TestWriteStruct   `hy:"nested"`         // like a new root
	Slice        []StructB          `hy:"slice/"`         // file per element
	NamedSlice2  []StructB          `hy:"a-named-slice/"` // file per element
	Map          map[string]StructB `hy:"map/"`           // file per element
	NamedMap     map[string]StructB `hy:"a-named-map/"`   // file per element
}

var testWriteStructData = TestWriteStruct{
	Name:        "Test struct writing",
	Int:         1,
	InlineSlice: []string{"a", "string", "slice"},
	InlineMap:   map[string]int{"one": 1, "two": 2, "three": 3},
	StructFile:  StructB{Name: "A file"},
	StringFile:  "A string in a file.",
	Nested: &TestWriteStruct{
		Name: "A nested struct pointer.",
		Int:  2,
		Slice: []StructB{
			{Name: "Nested One"}, {Name: "Nested Two"},
		},
		Nested: &TestWriteStruct{
			SliceFile: []string{"this", "is", "a", "slice", "in", "a", "file"},
			MapFile:   map[string]string{"deeply-nested": "map", "in a file": "yes"},
		},
	},
	Slice: []StructB{{Name: "One"}, {Name: "Two"}},
	Map: map[string]StructB{
		"First":  StructB{},
		"Second": StructB{},
	},
}

var testWriteFileTargets = map[string]FileTarget{
	"TestWriteStruct": FileTarget{
		Data: map[string]interface{}{
			"Name":        "Test struct writing",
			"Int":         1,
			"InlineSlice": []string{"a", "string", "slice"},
			"InlineMap":   map[string]int{"one": 1, "two": 2, "three": 3},
		},
	},
	"slice/1": FileTarget{
		Data: map[string]interface{}{
			"Name": "One",
		},
	},
	"slice/2": FileTarget{
		Data: map[string]interface{}{
			"Name": "Two",
		},
	},
	"map/First": FileTarget{
		Data: map[string]interface{}{
			"Name": "First",
		},
	},
	"map/Second": FileTarget{
		Data: map[string]interface{}{
			"Name": "Second",
		},
	},
}

func TestNode_WriteTargets_struct(t *testing.T) {
	c := NewCodec()
	n, err := c.Analyse(TestWriteStruct{})
	if err != nil {
		t.Fatal(err)
	}
	wc := NewWriteContext()
	v := reflect.ValueOf(testWriteStructData)
	if err := n.WriteTargets(wc, reflect.Value{}, v); err != nil {
		t.Fatal(err)
	}
	targets := wc.Targets
	expectedLen := 21
	if targets.Len() != expectedLen {
		t.Errorf("got len %d; want %d", targets.Len(), expectedLen)
		for k, ft := range targets.Snapshot() {
			data, err := json.MarshalIndent(ft.Data, "  ", "  ")
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("file: %s\n%s\n", k, data)
		}
	}
}

var testWriteFS = `
file: TestWriteStruct.yaml
Name: Test struct writing
Int: 1
InlineSlice:
	- a
	- string
	- slice
InlineMap:
	one: 1
	two: 2
	three: 3

file: Slice/1.yaml
Name: One

file: Slice/2.yaml
Name: Two

file: Map/First.yaml
Name: First

file: Map/Second.yaml
Name: Second
`
