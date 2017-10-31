package share

import (
	"errors"
)

func (p *Tree) Get(name string) *Tree {
	child, ok := p.children[name]
	if !ok {
		child = &Tree{
			parent:   p,
			children: make(map[string]*Tree),
		}
		p.children[name] = child
	}
	return child
}

func (p *Tree) Value() string {
	return p.value
}

func NewTree() *Tree {
	return &Tree{
		children: make(map[string]*Tree),
	}
}

type Tree struct {
	parent   *Tree
	children map[string]*Tree
	value    string
}
