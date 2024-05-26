package optimistic

import (
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

func (tree *Tree) validate(x int, node, parent *Node) bool {
	if node == nil && parent == nil {
		return tree.root == nil
	}
	curr := tree.root
	var prev *Node = nil
	for curr != nil && curr.val != x && curr != node {
		prev = curr
		if x < curr.val {
			curr = curr.left
		} else {
			curr = curr.right
		}
	}
	return curr == node && prev == parent
}

func (tree *Tree) findHelper(x int) (*Node, *Node) {
	for {
		tree.lock_tree()
		if tree.root == nil {
			return nil, nil
		}
		curr := tree.root
		var prev *Node = nil
		for curr != nil && curr.val != x {
			temp := prev
			prev = curr
			if x < curr.val {
				curr = curr.left
			} else {
				curr = curr.right
			}
			if temp == nil {
				tree.unlock_tree()
			}
		}

		if prev != nil {
			prev.lock()
		}
		if curr != nil {
			curr.lock()
		}

		if tree.validate(x, curr, prev) {
			return curr, prev
		}

		if curr != nil {
			curr.unlock()
		}
		if prev != nil {
			prev.unlock()
		}
	}
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
	defer curr.unlock()

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

	defer succ.unlock()
	if succ_parent != curr {
		defer succ_parent.unlock()
		succ_parent.left = succ.right
	} else {
		succ_parent.right = succ.right
	}
	curr.val = succ.val
}

// IsValid and IsEmpty does not support optimistic-lock semantic, so they are purely sequential
func isValid(node *Node) bool {
	if node == nil {
		return true
	}
	if node.left != nil && node.left.val >= node.val {
		return false
	}
	if node.right != nil && node.right.val <= node.val {
		return false
	}
	return isValid(node.left) && isValid(node.right)
}

func (tree *Tree) IsValid() bool {
	return isValid(tree.root)
}

func (tree *Tree) IsEmpty() bool {
	return tree.root == nil
}

func NewTree() *Tree {
	return &Tree{}
}
