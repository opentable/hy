package hy

import (
	"encoding/json"
	"testing"

	"github.com/kr/pretty"
)

func TestCodec_Read_json(t *testing.T) {
	c := NewCodec()

	v := TestStruct{}

	if err := c.Read("testdata/in", &v); err != nil {
		t.Fatal(err)
	}

	diffs := pretty.Diff(testDataSimple, v)
	for _, d := range diffs {
		t.Error(d)
	}
}

func TestCodec_Read_json_preserveOriginal(t *testing.T) {
	jsonWriter := JSONWriter
	jsonWriter.MarshalFunc = func(v interface{}) ([]byte, error) {
		return json.MarshalIndent(v, "", "  ")
	}
	c := NewCodec(func(c *Codec) {
		c.treeReader = NewFileTreeReader("json", "_")
		c.reader = jsonWriter
		c.writer = jsonWriter
	})

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
