package hy

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestNode_Write_struct(t *testing.T) {
	c := NewCodec()
	n, err := c.Analyse(TestStruct{})
	if err != nil {
		t.Fatal(err)
	}
	wc := NewWriteContext()
	v := reflect.ValueOf(testDataSimple)
	val := n.NewValFrom(v)
	if err := n.Write(wc, val); err != nil {
		t.Fatal(err)
	}
	targets := wc.targets
	expectedLen := len(testFileTargetsSnapshot)
	if targets.Len() != expectedLen {
		t.Errorf("got len %d; want %d", targets.Len(), expectedLen)
	}
	actualTargets := targets.Snapshot()
	for fileName, actual := range actualTargets {
		expected, ok := testFileTargetsSnapshot[fileName]
		if !ok {
			t.Errorf("extra file generated at %s:\n%s", fileName, actual.TestDump())
			continue
		}
		if actual.Value == nil && expected.Value == nil {
			continue
		}
		var actualType, expectedType reflect.Type
		if actual.Value != nil {
			actualType = reflect.ValueOf(actual.Value).Type()
			if expected.Value == nil {
				t.Errorf("at %q got: %v; want nil", fileName, actual.Value)
			}
		}
		if expected.Value != nil {
			expectedType = reflect.ValueOf(expected.Value).Type()
			if actual.Value == nil {
				t.Errorf("at %q got: nil; want:\n%v", fileName, expected.Value)
			}
		}

		if actualType != expectedType {
			t.Errorf("got type %s; want %s at %q", actualType, expectedType, fileName)
			t.Errorf("values: got:\n%# v\nwant:\n%# v", actual.Value, expected.Value)
		}
		if actual.TestDataDump() != expected.TestDataDump() {
			t.Errorf("\ngot rendered data:\n%s\nwant:\n%s\n",
				actual.TestDump(), expected.TestDump())
		}
	}
	for fileName := range testFileTargetsSnapshot {
		if _, ok := actualTargets[fileName]; !ok {
			t.Errorf("missing file %q", fileName)
		}
	}
}

func (ft FileTarget) TestDump() string {
	return fmt.Sprintf("file: %q\n%s\n", ft.FilePath, ft.TestDataDump())
}

func (ft FileTarget) TestDataDump() string {
	data, err := json.MarshalIndent(ft.Value, "  ", "  ")
	if err != nil {
		panic(err)
	}
	return string(data)
}

type NoFields struct {
	Child *NoFields `hy:"nofields"`
}

func TestNode_Write_noFields(t *testing.T) {
	c := NewCodec()
	thing := NoFields{&NoFields{}}
	node, err := c.Analyse(thing)
	if err != nil {
		t.Fatal(err)
	}
	ctx := NewWriteContext()
	node.Write(ctx, node.NewValFrom(reflect.ValueOf(thing)))
	if ctx.targets.Len() != 0 {
		t.Errorf("got %d targets; want 0", ctx.targets.Len())
		t.Error(ctx.targets.Paths())
	}
}

type NoChildren struct {
	B, C, A string
}

func TestNode_Write_orderedFields(t *testing.T) {
	c := NewCodec()
	thing := NoChildren{"B", "C", "A"}
	node, err := c.Analyse(thing)
	if err != nil {
		t.Fatal(err)
	}
	ctx := NewWriteContext()
	node.Write(ctx, node.NewValFrom(reflect.ValueOf(thing)))
	expectedNumTargets := 1
	if ctx.targets.Len() != expectedNumTargets {
		t.Errorf("got %d targets; want %d", ctx.targets.Len(), expectedNumTargets)
	}
	target, ok := ctx.targets.Snapshot()[""]
	if !ok {
		t.Fatal("root target (with empty path) not found")
	}
	actualData := target.Data()
	expectedData := thing
	if reflect.TypeOf(actualData) != reflect.TypeOf(expectedData) {
		t.Fatalf("root target data is %T; want %T", actualData, expectedData)
	}
}
