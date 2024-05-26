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
	seed = 0xC0FFEE
	rnd = rand.New(rand.NewSource(seed))
	values = make([]int, number_size)
	for i := range values {
		values[i] = rnd.Intn(1000)
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

		t.Run("Random operations", func(t *testing.T) {
			for _, v := range values {
				wg.Add(1)
				go func(x int) {
					defer wg.Done()
					op := rnd.Intn(100)
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
