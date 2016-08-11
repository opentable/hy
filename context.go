package hy

import (
	"path"

	"github.com/pkg/errors"
)

// WriteContext is context collected during a write opration.
type WriteContext struct {
	// Targets is the collected targets in this write context.
	Targets FileTargets
	// Parent is the parent write context.
	Parent *WriteContext
	// PathName is the name of this section of the path.
	PathName string
}

// NewWriteContext returns a new write context.
func NewWriteContext() WriteContext {
	return WriteContext{Targets: MakeFileTargets(0)}
}

// Push creates a derivative node context.
func (c WriteContext) Push(pathName string) WriteContext {
	return WriteContext{
		Targets:  c.Targets,
		Parent:   &c,
		PathName: pathName,
	}
}

// Path returns the path of this context.
func (c WriteContext) Path() string {
	if c.Parent == nil {
		return c.PathName
	}
	return path.Join(c.Parent.Path(), c.PathName)
}

// SetValue sets the value of the current path.
func (c WriteContext) SetValue(v interface{}) error {
	t := &FileTarget{FilePath: c.Path(), Value: v}
	return errors.Wrapf(c.Targets.Add(t), "setting value at %q", c.Path())
}
