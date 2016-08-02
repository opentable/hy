package hy

import (
	"reflect"

	"github.com/pkg/errors"
)

// NodeSet is a set of Node pointers indexed by ID.
type NodeSet struct {
	nodes map[NodeID]*Node
}

// NewNodeSet creates a new node set.
func NewNodeSet() NodeSet {
	return NodeSet{nodes: map[NodeID]*Node{}}
}

// Register tries to register a node ID. If the ID is not yet registered, it
// returns a new node pointer and true. Otherwise it returns the already
// registered node pointer and false.
func (ns NodeSet) Register(id NodeID) (*Node, bool) {
	n, ok := ns.nodes[id]
	if ok {
		return n, false
	}
	n = new(Node)
	ns.nodes[id] = n
	return n, true
}

// Codec provides the primary encoding and decoding facility of this package.
type Codec struct {
	Nodes NodeSet
}

// NewCodec creates a new codec.
func NewCodec() *Codec {
	return &Codec{Nodes: NewNodeSet()}
}

// Analyse analyses a tree starting at root.
func (c *Codec) Analyse(root interface{}) (Node, error) {
	if root == nil {
		return nil, errors.New("cannot analyse nil")
	}
	n, err := c.analyse(nil, reflect.TypeOf(root), FieldInfo{})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to analyse %T", root)
	}
	return *n, err
}

func (c *Codec) analyse(parent Node, t reflect.Type, field FieldInfo) (*Node, error) {
	var isPtr bool
	k := t.Kind()
	if k == reflect.Ptr {
		isPtr = true
		t = t.Elem()
		k = t.Kind()
	}
	var parentType reflect.Type
	if parent != nil {
		parentType = parent.ID().Type
	}
	nodeID := NodeID{
		ParentType: parentType,
		Type:       t,
		IsPtr:      isPtr,
		FieldName:  field.Name,
	}
	n, new := c.Nodes.Register(nodeID)
	if !new {
		return n, nil
	}
	var err error
	base := NodeBase{NodeID: nodeID, Parent: parent, Tag: field.Tag}
	switch k {
	default:
		return nil, errors.Errorf("cannot analyse kind %s", k)
	case reflect.Struct:
		*n, err = c.analyseStruct(base)
	case reflect.Map:
		*n, err = c.analyseMap(base)
	case reflect.Slice:
		*n, err = c.analyseSlice(base)
	}
	return n, errors.Wrapf(err, "analysing %s failed", t)
}
