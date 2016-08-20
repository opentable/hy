package hy

// testDataSimple has all keys set, so it round-trips unchanged.
// This is in contrast to testDataUnsetKeys which has its keys  set on being
// written (see testFilesUnsetKeys for the read equivalent.
var testDataSimple = TestStruct{
	Name:        "Test struct writing",
	Int:         1,
	InlineSlice: []string{"a", "string", "slice"},
	InlineMap:   map[string]int{"one": 1, "two": 2, "three": 3},
	StructFile:  StructB{Name: "A file"},
	StringFile:  "A string in a file.",
	Nested: &TestStruct{
		Name: "A nested struct pointer.",
		Int:  2,
		Slice: []StructB{
			{Name: "Nested One"}, {Name: "Nested Two"},
		},
		Nested: &TestStruct{
			SliceFile: []string{"this", "is", "a", "slice", "in", "a", "file"},
			MapFile:   map[string]string{"deeply-nested": "map", "in a file": "yes"},
		},
		StructFile: StructB{
			Name: "Struct B file",
		},
		MapOfPtr: map[string]*StructB{
			"a-nil-file":           {Name: "a-nil-file"},
			"another-nil-file":     {Name: "another-nil-file"},
			"this-one-has-a-value": {Name: "this-one-has-a-value"},
		},
		Map: map[string]StructB{
			"a-zero-file":       {Name: "a-zero-file"},
			"another-zero-file": {Name: "another-zero-file"},
		},
	},
	Slice: []StructB{{Name: "One"}, {Name: "Two"}},
	Map: map[string]StructB{
		"First":  {Name: "First"},
		"Second": {Name: "Second"},
	},
	TextMarshalerKey: map[TextMarshaler]*StructB{
		// Slashes should translate to directories.
		{"Test/blah/blah", 1}: nil,
		{"Another/blah", 13}:  nil,
	},
	TextMarshalerPtrKey: map[*TextMarshaler]*StructB{
		{"Test", 2}:     nil,
		{"Another", 14}: nil,
	},
	SpecialMap: SpecialMap{m: map[TextMarshaler]*StructB{
		{"Special", 3}:  {Name: "Special"},
		{"Another", 15}: nil,
	}},
	SpecialMapPtr: &SpecialMap{m: map[TextMarshaler]*StructB{
		{"Special", 4}:  {Name: "Special Ptr"},
		{"Another", 16}: nil,
	}},
	SpecialPtrMap: SpecialPtrMap{m: map[TextMarshaler]*StructB{
		{"Special", 5}:  {Name: "Ptr Special"},
		{"Another", 17}: nil,
	}},
	SpecialPtrMapPtr: &SpecialPtrMap{m: map[TextMarshaler]*StructB{
		{"Special", 6}:  {Name: "Ptr Special Ptr"},
		{"Another", 18}: nil,
	}},
}

var testDataUnsetKeys = TestStruct{
	Name:        "Test struct writing",
	Int:         1,
	InlineSlice: []string{"a", "string", "slice"},
	InlineMap:   map[string]int{"one": 1, "two": 2, "three": 3},
	StructFile:  StructB{Name: "A file"},
	StringFile:  "A string in a file.",
	Nested: &TestStruct{
		Name: "A nested struct pointer.",
		Int:  2,
		Slice: []StructB{
			{Name: "Nested One"}, {Name: "Nested Two"},
		},
		Nested: &TestStruct{
			SliceFile: []string{"this", "is", "a", "slice", "in", "a", "file"},
			MapFile:   map[string]string{"deeply-nested": "map", "in a file": "yes"},
		},
		StructFile: StructB{
			Name: "Struct B file",
		},
		MapOfPtr: map[string]*StructB{
			"a-nil-file":           nil,
			"another-nil-file":     nil,
			"this-one-has-a-value": {},
		},
		Map: map[string]StructB{
			// Notice how we don't set the Name field here. Hy sets it in the write
			// data because of the ",Name" tag.
			"a-zero-file":       {},
			"another-zero-file": {},
		},
	},
	Slice: []StructB{{Name: "One"}, {Name: "Two"}},
	Map: map[string]StructB{
		// Notice how we don't set the Name field here. Hy sets it in the write
		// data because of the ",Name" tag.
		"First":  {},
		"Second": {},
	},
	TextMarshalerKey: map[TextMarshaler]*StructB{
		// Slashes should translate to directories.
		{"Test/blah/blah", 1}: nil,
		{"Another/blah", 13}:  nil,
	},
	TextMarshalerPtrKey: map[*TextMarshaler]*StructB{
		{"Test", 2}:     nil,
		{"Another", 14}: nil,
	},
	SpecialMap: SpecialMap{m: map[TextMarshaler]*StructB{
		{"Special", 3}:  {Name: "Special"},
		{"Another", 15}: nil,
	}},
	SpecialMapPtr: &SpecialMap{m: map[TextMarshaler]*StructB{
		{"Special", 4}:  {Name: "Special Ptr"},
		{"Another", 16}: nil,
	}},
	SpecialPtrMap: SpecialPtrMap{m: map[TextMarshaler]*StructB{
		{"Special", 5}:  {Name: "Ptr Special"},
		{"Another", 17}: nil,
	}},
	SpecialPtrMapPtr: &SpecialPtrMap{m: map[TextMarshaler]*StructB{
		{"Special", 6}:  {Name: "Ptr Special Ptr"},
		{"Another", 18}: nil,
	}},
}
