package main

import (
	"fmt"
	finegrained "trees/fine-grained"
)

func main() {
	tree := finegrained.NewTree()
	tree.Insert(10)
	tree.Insert(5)
	tree.Insert(15)
	tree.Insert(20)
	fmt.Println(tree.Find(15))
	fmt.Println(tree.Find(11))
	tree.Remove(10)
	tree.Remove(11)
	fmt.Println(tree.Find(10))
	fmt.Println(tree.Find(11))
}
