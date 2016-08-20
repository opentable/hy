package hy

import (
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
)

// TestStruct is the primary data structure used by tests, it should illustrate
// every significant combination of field type and tag.
type TestStruct struct {
	unexportedField     string                      // not output anywhere
	Name                string                      // regular field
	Int                 int                         // regular field
	InlineSlice         []string                    // regular field
	InlineMap           map[string]int              // regular field
	StructB             StructB                     // regular field
	StructBPtr          *StructB                    // regular field
	IgnoredField        string                      `hy:"-"`                 // not output anywhere
	StructFile          StructB                     `hy:"a-file"`            // a single file
	StringFile          string                      `hy:"a-string-file"`     // a single file
	SliceFile           []string                    `hy:"a-slice-file"`      // a single file
	MapFile             map[string]string           `hy:"a-map-file"`        // a single file
	Nested              *TestStruct                 `hy:"nested"`            // like a new root
	Slice               []StructB                   `hy:"slice/"`            // file per element
	Map                 map[string]StructB          `hy:"map/,Name"`         // file per element
	MapOfPtr            map[string]*StructB         `hy:"map-of-ptr/,Name"`  // file per element
	TextMarshalerKey    map[TextMarshaler]*StructB  `hy:"textmarshaler/"`    // file per element
	TextMarshalerPtrKey map[*TextMarshaler]*StructB `hy:"textmarshalerptr/"` // file per element
	SpecialMap          SpecialMap                  `hy:"specialmap/"`       // file per element
	SpecialMapPtr       *SpecialMap                 `hy:"specialmapptr/"`    // file per element
	SpecialPtrMap       SpecialPtrMap               `hy:"specialptrmap/"`    // file per element
	SpecialPtrMapPtr    *SpecialPtrMap              `hy:"specialptrmapptr/"` // file per element
}

// SpecialMap is a wrapper around a map. Hy treats it like a map because it has
// GetAll and SetAll methods. The return and accept the same map type,
// respectively.
type SpecialMap struct {
	m map[TextMarshaler]*StructB
}

func (s SpecialMap) SetAll(m map[TextMarshaler]*StructB) { s.m = m }
func (s SpecialMap) GetAll() map[TextMarshaler]*StructB  { return s.m }

// SpecialPtrMap is very similar to SpecialMap, except that we define its GetAll
// and SetAll methods on its pointer receiver.
type SpecialPtrMap struct {
	m map[TextMarshaler]*StructB
}

func (s *SpecialPtrMap) SetAll(m map[TextMarshaler]*StructB) { s.m = m }
func (s *SpecialPtrMap) GetAll() map[TextMarshaler]*StructB  { return s.m }

// TextMarshaler is used to illustrate using non-string map keys. Hy will write
// any map key which is either a string, or implements encoding.MarshalText.
// Likewise, Hy will read any map key which is either a string or implements
// encoding.UnmarshalText.
type TextMarshaler struct {
	String string
	Int    int
}

func (tm TextMarshaler) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%s-%d", tm.String, tm.Int)), nil
}

func (tm *TextMarshaler) UnmarshalText(text []byte) error {
	s := string(text)
	s = strings.Replace(s, "-", " ", -1)
	n, err := fmt.Sscanf(s, "%s %d", &tm.String, &tm.Int)
	if err != nil && err != io.EOF {
		return errors.Wrapf(err, "unmarshaling %s", s)
	}
	if n != 2 {
		return errors.Errorf("%s has %d missing fields", s, 2-n)
	}
	return nil
}
