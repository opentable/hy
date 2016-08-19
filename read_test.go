package hy

import (
	"encoding/json"
	"os"
	"testing"
)

func TestCodec_Read_json_roundTrip(t *testing.T) {
	c := NewCodec(func(c *Codec) {
		c.MarshalFunc = func(v interface{}) ([]byte, error) {
			return json.MarshalIndent(v, "", "  ")
		}
	})

	v := TestStruct{}

	if err := c.Read("testdata/in", &v); err != nil {
		t.Fatal(err)
	}

	if err := os.RemoveAll("testdata/roundtripped"); err != nil {
		t.Fatal(err)
	}
	if err := c.Write("testdata/roundtripped", v); err != nil {
		t.Fatal(err)
	}

	v2 := TestStruct{}
	if err := c.Read("testdata/roundtripped", &v2); err != nil {
		t.Fatal(err)
	}

	if err := os.RemoveAll("testdata/roundtripped2"); err != nil {
		t.Fatal(err)
	}
	if err := c.Write("testdata/roundtripped2", &v2); err != nil {
		t.Fatal(err)
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
