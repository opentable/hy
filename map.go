package hy

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

// A MapNode represents a map node to be stored in a directory.
type MapNode struct {
	*DirNodeBase
	KeyType reflect.Type
}

// NewMapNode makes a new map node.
func (c *Codec) NewMapNode(base NodeBase) (Node, error) {
	n := &MapNode{
		DirNodeBase: &DirNodeBase{
			NodeBase: base,
		},
		KeyType: base.Type.Key(),
	}
	return n, errors.Wrap(n.AnalyseElemNode(n, c), "analysing map element node")
}

// ChildPathName returns the key as a string.
func (n *MapNode) ChildPathName(child Node, key, val reflect.Value) string {
	return fmt.Sprint(key)
}

// WriteTargets writes all map elements.
func (n *MapNode) WriteTargets(c WriteContext, key, val reflect.Value) (FileTargets, error) {
	fts := MakeFileTargets(val.Len())
	elemNode := *n.ElemNode
	for _, k := range val.MapKeys() {
		v := val.MapIndex(k)
		childContext := c.Push(elemNode.PathName(k, v))
		childTargets, err := elemNode.WriteTargets(childContext, k, v)
		if err != nil {
			return fts, errors.Wrapf(err, "writing map index %q failed", fmt.Sprint(k))
		}
		if err := fts.AddAll(childTargets); err != nil {
			return fts, errors.Wrapf(err, "adding children of key %q failed", fmt.Sprint(k))
		}
	}
	return fts, nil
}
