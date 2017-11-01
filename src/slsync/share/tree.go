package share

import (
	"fmt"
)

type Leaves map[string]Leaf

type Leaf struct {
	Clocks string
	Value  string
}

type Tree struct {
	parent   *Tree
	children map[string]*Tree
	clocks   Clocks
	value    string
}

func (p *Tree) GetLeaves(prefix string) (leaves Leaves) {
	leaves = make(Leaves)
	p.leaves(prefix, leaves)
	return
}

func (p *Tree) leaves(prefix string, result Leaves) {
	if len(p.clocks) != 0 {
		result[prefix] = Leaf{
			Clocks: p.clocks.ToString(),
			Value:  p.value,
		}
	}

	fmt.Println("p.child:", p.children)

	for k, v := range p.children {
		fmt.Println("child:", k, p.children[k])
		path := k
		if len(prefix) != 0 {
			path = prefix + "/" + path
		}
		v.leaves(path, result)
	}
}

func (p *Tree) Get(name string) *Tree {
	child, ok := p.children[name]
	if !ok {
		child = &Tree{
			parent:   p,
			children: make(map[string]*Tree),
			clocks:   NewClocks(),
		}
		p.children[name] = child
	}
	return child
}

func (p *Tree) Clocks() Clocks {
	return p.clocks
}

func (p *Tree) Value() string {
	return p.value
}

func NewTree() *Tree {
	return &Tree{
		children: make(map[string]*Tree),
		clocks:   NewClocks(),
	}
}
