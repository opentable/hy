package hy

import (
	"encoding/json"
	"os"
	"sync/atomic"
	"testing"
)

func counter() *int64    { var c int64; return &c }
func increment(c *int64) { atomic.AddInt64(c, 1) }

const prefix = "testdata/out"

type NewCodecExpectation struct {
	Conf     func(*Codec)
	Expected Codec
}

var newCodecCalls = []NewCodecExpectation{
	{nil,
		Codec{FileExtension: "json", RootFileName: "_"}},
	{func(c *Codec) {},
		Codec{FileExtension: "json", RootFileName: "_"}},
	{func(c *Codec) { c.FileExtension = "blah" },
		Codec{FileExtension: "blah", RootFileName: "_"}},
	{func(c *Codec) { c.RootFileName = "blah2" },
		Codec{FileExtension: "json", RootFileName: "blah2"}},
}

func TestNewCodec_success(t *testing.T) {
	for _, expectation := range newCodecCalls {
		actual := NewCodec(expectation.Conf)
		expected := expectation.Expected
		if actual.FileExtension != expected.FileExtension {
			t.Errorf("got FileExtension %q; want %q", actual.FileExtension, expected.FileExtension)
		}
		if actual.RootFileName != expected.RootFileName {
			t.Errorf("got RootFileName %q; want %q", actual.RootFileName, expected.RootFileName)
		}
	}
}

func TestCodec_Write(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	if err := os.RemoveAll(prefix); err != nil {
		t.Fatalf("failed to remove output dir: %s", err)
	}

	var numCalls int64
	c := NewCodec(func(c *Codec) {
		c.MarshalFunc = func(v interface{}) ([]byte, error) {
			increment(&numCalls)
			return json.MarshalIndent(v, "", "  ")
		}

	})

	if err := c.Write(prefix, testData); err != nil {
		t.Fatal(err)
	}

	expectedNumCalls := int64(31)
	if numCalls != expectedNumCalls {
		t.Errorf("MarshalFunc called %d times; want %d", numCalls, expectedNumCalls)
	}
}
