package coarsegrained

import (
	"sync"
	"trees/sequential"
)

type Tree struct {
	mutex    sync.Mutex
	seq_tree sequential.Tree
}

func (tree *Tree) lock() {
	tree.mutex.Lock()
}

func (tree *Tree) unlock() {
	tree.mutex.Unlock()
}

func (tree *Tree) Insert(x int) {
	tree.lock()
	defer tree.unlock()
	tree.seq_tree.Insert(x)
}

func (tree *Tree) Find(x int) bool {
	tree.lock()
	defer tree.unlock()
	return tree.seq_tree.Find(x)
}

func (tree *Tree) Remove(x int) {
	tree.lock()
	defer tree.unlock()
	tree.seq_tree.Remove(x)
}

func (tree *Tree) InOrderPrint() {
	tree.lock()
	defer tree.unlock()
	tree.seq_tree.InOrderPrint()
}

func (tree *Tree) IsValid() bool {
	tree.lock()
	defer tree.unlock()
	return tree.seq_tree.IsValid()
}

func (tree *Tree) IsEmpty() bool {
	tree.lock()
	defer tree.unlock()
	return tree.seq_tree.IsEmpty()
}

func NewTree() *Tree {
	return &Tree{}
}
