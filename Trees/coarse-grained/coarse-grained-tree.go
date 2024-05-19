package coarsegrained

import (
	"sync"
	"trees/sequential"
)

type Tree struct {
	mutex    sync.Mutex
	seq_tree sequential.Tree
}

func (tree *Tree) Insert(x int) {
	tree.mutex.Lock()
	tree.seq_tree.Insert(x)
	tree.mutex.Unlock()
}

func (tree *Tree) Find(x int) bool {
	tree.mutex.Lock()
	val := tree.seq_tree.Find(x)
	tree.mutex.Unlock()
	return val
}

func (tree *Tree) Remove(x int) {
	tree.mutex.Lock()
	tree.seq_tree.Remove(x)
	tree.mutex.Unlock()
}

func (tree *Tree) InOrderPrint() {
	tree.mutex.Lock()
	tree.seq_tree.InOrderPrint()
	tree.mutex.Unlock()
}

func NewTree() Tree {
	return Tree{}
}
