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

func insert(node *Node, x int) {
	if node.val == x {
		return
	}
	if x < node.val {
		if node.left == nil {
			node.left = &Node{val: x}
		} else {
			insert(node.left, x)
		}
	}
	if x > node.val {
		if node.right == nil {
			node.right = &Node{val: x}
		} else {
			insert(node.right, x)
		}
	}
}

func (tree *Tree) Insert(x int) {
	if tree.root == nil {
		tree.root = &Node{val: x}
		return
	}
	insert(tree.root, x)
}

func find(node *Node, x int) bool {
	if node == nil {
		return false
	}
	if x < node.val {
		return find(node.left, x)
	}
	if x > node.val {
		return find(node.right, x)
	}
	return true
}

func (tree *Tree) Find(x int) bool {
	return find(tree.root, x)
}

func remove(node *Node, x int) *Node {
	if node == nil {
		return nil
	}
	if x < node.val {
		node.left = remove(node.left, x)
		return node
	}
	if x > node.val {
		node.right = remove(node.right, x)
		return node
	}

	if node.left == nil && node.right == nil {
		return nil
	}
	if node.left == nil {
		return node.right
	}
	if node.right == nil {
		return node.left
	}
	succ := node.right
	for succ.left != nil {
		succ = succ.left
	}
	node.val = succ.val
	node.right = remove(node.right, node.val)
	return node
}

func (tree *Tree) Remove(x int) {
	remove(tree.root, x)
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

func NewTree() Tree {
	return Tree{}
}
