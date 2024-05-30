package main

import (
	"math/rand"
	"sync"
	"testing"
	coarsegrained "trees/coarse-grained"
	finegrained "trees/fine-grained"
	"trees/optimistic"
)

var seed int64
var rnd *rand.Rand
var values []int
var number_size = 100_000

type BST interface {
	Insert(x int)
	Find(x int) bool
	Remove(x int)
	IsValid() bool
	IsEmpty() bool
}

func init() {
	seed = 42
	rnd = rand.New(rand.NewSource(seed))
	values = make([]int, number_size)
	for i := range values {
		values[i] = rnd.Int()
	}
}

func TestBST(t *testing.T) {
	trees := map[string]BST{
		"coarse-grained": coarsegrained.NewTree(),
		"fine-grained":   finegrained.NewTree(),
		"optimistic":     optimistic.NewTree(),
	}
	var wg sync.WaitGroup

	for name, tree := range trees {

		/*
			This tests check that after only insert operations tree will remain BST
			(left children's value is less than parent value, right children's -- is greater).
		*/
		t.Run("Insert", func(t *testing.T) {
			for _, v := range values {
				wg.Add(1)
				go func(x int) {
					defer wg.Done()
					tree.Insert(x)
				}(v)
			}
			wg.Wait()
			if !tree.IsValid() {
				t.Errorf("Insert test failed: %s tree is not valid.", name)
			}
		})

		/*
			This tests check that after insert operations all numbers inserted will present in BST.
		*/
		t.Run("Find after Insert", func(t *testing.T) {
			for _, v := range values {
				wg.Add(1)
				go func(x int) {
					defer wg.Done()
					tree.Insert(x)
				}(v)
			}
			wg.Wait()
			for _, v := range values {
				if !tree.Find(v) {
					t.Errorf("Find after Insert test failed: %s tree doesn't contain %d.", name, v)
				}
			}
		})

		/*
			This tests check what after inserting and removing the same elements BST will be empty.
		*/
		t.Run("Remove after Insert", func(t *testing.T) {
			for _, v := range values {
				wg.Add(1)
				go func(x int) {
					defer wg.Done()
					tree.Insert(x)
				}(v)
			}
			wg.Wait()
			for _, v := range values {
				wg.Add(1)
				go func(x int) {
					defer wg.Done()
					tree.Remove(x)
				}(v)
			}
			wg.Wait()
			if !tree.IsEmpty() {
				t.Errorf("Remove after Insert test failed: %s tree is not valid.", name)
			}
		})

		/*
			This tests check that during simultanious insert and remove operation tree will still be valid.
		*/
		t.Run("Insert and Remove", func(t *testing.T) {
			for i, v := range values {
				wg.Add(1)
				go func(x, i int) {
					defer wg.Done()
					if i%2 == 0 {
						tree.Insert(x)
					} else {
						tree.Remove(x)
					}
				}(v, i)
			}
			wg.Wait()
			if !tree.IsValid() {
				t.Errorf("Insert and Remove test failed: %s tree is not valid.", name)
			}
		})

		/*
			This tests check that during simultanious random operation tree will still be valid.
		*/
		t.Run("Random operations", func(t *testing.T) {
			for _, v := range values {
				wg.Add(1)
				go func(x int) {
					defer wg.Done()
					rnd := rand.New(rand.NewSource(int64(x)))
					op := rnd.Int()
					if op%3 == 0 {
						tree.Insert(x)
					} else if op%3 == 1 {
						tree.Find(x)
					} else {
						tree.Remove(x)
					}
				}(v)
			}
			wg.Wait()
			if !tree.IsValid() {
				t.Errorf("Random operations test failed: %s tree is not valid.", name)
			}
		})
	}
}
