package hy

import (
	"reflect"
	"testing"

	"github.com/kr/pretty"
)

func TestNode_Read_struct(t *testing.T) {
	c := NewCodec()
	n, err := c.Analyse(TestStruct{})
	if err != nil {
		t.Fatal(err)
	}
	reader := c.reader
	targets := testFileTargets
	rc := NewReadContext("testdata/in", targets, reader)
	v := reflect.ValueOf(testDataSimple)
	val := n.NewValFrom(v)
	if err := n.Read(rc, val); err != nil {
		t.Fatal(err)
	}
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

func TestCodec_Read_json(t *testing.T) {
	t.Skip()
	c := NewCodec()

	actual := TestStruct{}
	expected := testDataSimple

	if err := c.Read("testdata/in", &actual); err != nil {
		t.Fatal(err)
	}

	diffs := pretty.Diff(expected, actual)
	for _, d := range diffs {
		t.Error(d, "\n")
	}
}

func TestCodec_Read_json_preserveOriginalFieldsWhenNotRead(t *testing.T) {
	c := NewCodec()
	v := TestStruct{
		unexportedField: "this field should remain",
	}
	if err := c.Read("testdata/in", &v); err != nil {
		t.Fatal(err)
	}
	if v.unexportedField != "this field should remain" {
		t.Fatal("unexported field not preserved")
	}
}
