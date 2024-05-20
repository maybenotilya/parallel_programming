package sequential

import "fmt"

type Node struct {
	val   int
	left  *Node
	right *Node
}

type Tree struct {
	root *Node
}

func (tree *Tree) findHelper(x int) (*Node, *Node) {
	if tree.root == nil {
		return nil, nil
	}
	curr := tree.root
	var prev *Node = nil
	for curr != nil && curr.val != x {
		prev = curr
		if x < curr.val {
			curr = curr.left
		} else {
			curr = curr.right
		}
	}
	return curr, prev
}

func (tree *Tree) Insert(x int) {
	curr, prev := tree.findHelper(x)
	if tree.root == nil {
		tree.root = &Node{val: x}
		return
	}
	if curr != nil {
		return
	}
	node := &Node{val: x}
	if x < prev.val {
		prev.left = node
	} else {
		prev.right = node
	}
}

func (tree *Tree) Find(x int) bool {
	curr, _ := tree.findHelper(x)
	return curr != nil
}

func (tree *Tree) Remove(x int) {
	curr, prev := tree.findHelper(x)
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
	succ_parent := curr
	succ := curr.right
	for succ.left != nil {
		succ_parent = succ
		succ = succ.left
	}
	if succ_parent != curr {
		succ_parent.left = succ.right
	} else {
		succ_parent.right = succ.right
	}
	curr.val = succ.val
}

func inOrderPrint(node *Node) {
	if node == nil {
		return
	}
	inOrderPrint(node.left)
	fmt.Print(node.val, " ")
	inOrderPrint(node.right)
}

func (tree *Tree) InOrderPrint() {
	inOrderPrint(tree.root)
	fmt.Println()
}

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
