package finegrained

import (
	"fmt"
	"sync"
)

type Node struct {
	val   int
	left  *Node
	right *Node
	mutex sync.Mutex
}

func (node *Node) lock() {
	node.mutex.Lock()
}

func (node *Node) unlock() {
	node.mutex.Unlock()
}

type Tree struct {
	root  *Node
	mutex sync.Mutex
}

func (tree *Tree) lock_tree() {
	tree.mutex.Lock()
}

func (tree *Tree) unlock_tree() {
	tree.mutex.Unlock()
}

func (tree *Tree) findHelper(x int) (*Node, *Node) {
	tree.lock_tree()
	if tree.root == nil {
		return nil, nil
	}
	tree.root.lock()
	curr := tree.root
	var prev *Node = nil
	for curr != nil && curr.val != x {
		temp := prev
		prev = curr
		if x < curr.val {
			if curr.left != nil {
				curr.left.lock()
			}
			curr = curr.left
		} else {
			if curr.right != nil {
				curr.right.lock()
			}
			curr = curr.right
		}
		if temp == nil {
			tree.unlock_tree()
		} else {
			temp.unlock()
		}
	}
	return curr, prev
}

func (tree *Tree) Insert(x int) {
	curr, prev := tree.findHelper(x)
	node := &Node{val: x}
	if tree.root == nil {
		defer tree.unlock_tree()
		tree.root = node
		return
	}
	if prev == nil {
		defer tree.unlock_tree()
	} else {
		defer prev.unlock()
	}
	if curr != nil {
		defer curr.unlock()
		return
	}
	if x < prev.val {
		prev.left = node
	} else {
		prev.right = node
	}
}

func (tree *Tree) Find(x int) bool {
	curr, prev := tree.findHelper(x)
	if prev == nil {
		defer tree.unlock_tree()
	} else {
		defer prev.unlock()
	}
	if curr != nil {
		defer curr.unlock()
		return true
	}
	return false
}

func (tree *Tree) Remove(x int) {
	curr, prev := tree.findHelper(x)
	if prev == nil {
		defer tree.unlock_tree()
	} else {
		defer prev.unlock()
	}
	if curr == nil {
		return
	}

	if curr.left == nil && curr.right == nil {
		if curr == tree.root {
			tree.root = nil
		} else if curr.val < prev.val {
			prev.left = nil
		} else {
			prev.right = nil
		}
		return
	}

	if curr.left == nil {
		if curr == tree.root {
			tree.root = curr.right
		} else if curr.val < prev.val {
			prev.left = curr.right
		} else {
			prev.right = curr.right
		}
		return
	}

	if curr.right == nil {
		if curr == tree.root {
			tree.root = curr.left
		} else if curr.val < prev.val {
			prev.left = curr.left
		} else {
			prev.right = curr.left
		}
		return
	}
	defer curr.unlock()
	curr.right.lock()
	succ_parent := curr
	succ := curr.right
	for succ.left != nil {
		temp := succ_parent
		succ_parent = succ
		succ.left.lock()
		succ = succ.left
		if temp != nil && temp != curr {
			temp.unlock()
		}
	}
	if succ_parent != curr {
		defer succ_parent.unlock()
		succ_parent.left = succ.right
	} else {
		succ_parent.right = succ.right
	}
	curr.val = succ.val
}

// Seems like it is optimal to lock one node only
func inOrderPrint(node *Node) {
	if node == nil {
		return
	}
	node.lock()
	defer node.unlock()
	inOrderPrint(node.left)
	fmt.Print(node.val, " ")
	inOrderPrint(node.right)
}

func (tree *Tree) InOrderPrint() {
	inOrderPrint(tree.root)
	fmt.Println()
}

func isValid(node *Node) bool {
	//At this point node is locked, overwise deadlock is possible
	if node == nil {
		return true
	}
	if node.left != nil {
		node.left.lock()
	}
	if node.right != nil {
		node.right.lock()
	}
	defer node.unlock()
	if node.left != nil && node.left.val >= node.val {
		return false
	}
	if node.right != nil && node.right.val <= node.val {
		return false
	}
	return isValid(node.left) && isValid(node.right)
}

func (tree *Tree) IsValid() bool {
	if tree.root == nil {
		return true
	}
	tree.root.lock()
	return isValid(tree.root)
}

func (tree *Tree) IsEmpty() bool {
	return tree.root == nil
}

func NewTree() *Tree {
	return &Tree{}
}
