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
	//if n.Field != nil && n.Field.KeyField != "" {
	//	n.Field.GetKeyFunc.Call([]reflect.Value{val})
	//}
	return fmt.Sprint(key)
}

// WriteTargets writes all map elements.
func (n *MapNode) WriteTargets(c WriteContext, key, val reflect.Value) error {
	elemNode := *n.ElemNode
	for _, k := range val.MapKeys() {
		v := val.MapIndex(k)
		// make an addressable copy of v
		// this is ripe for refactoring so we don't need to jump through hoops
		// to get the address.
		newVal := reflect.New(v.Type()).Elem()
		newVal.Set(v)
		v = newVal
		vAddr := v
		if vAddr.Kind() != reflect.Ptr {
			vAddr = v.Addr()
		}
		if n.Field != nil && n.Field.KeyField != "" {
			n.Field.SetKeyFunc.Call([]reflect.Value{vAddr, k})
		}
		//log.Printf("Writing %s[%s] = %+ v\n", n.Type, k, v)
		childContext := c.Push(elemNode.PathName(k, v))
		if err := elemNode.Write(childContext, k, v); err != nil {
			return errors.Wrapf(err, "writing map index %q failed", fmt.Sprint(k))
		}
	}
	return nil
}
