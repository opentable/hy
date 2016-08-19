package hy

import (
	"encoding/json"
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

	actualJSONBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	expectedJSONBytes, err := json.MarshalIndent(testData, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	actualJSON := string(actualJSONBytes)
	expectedJSON := string(expectedJSONBytes)
	if actualJSON != expectedJSON {
		t.Errorf("\n\ngot:\n\n%s\n\nwant:\n\n%s", actualJSON, expectedJSON)
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
