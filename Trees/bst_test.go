package main

import (
	"math/rand"
	"sync"
	"testing"
	coarsegrained "trees/coarse-grained"
	finegrained "trees/fine-grained"
)

var values []int
var numberSize = 100_000

type BST interface {
	Insert(x int)
	Find(x int) bool
	Remove(x int)
	IsValid() bool
	IsEmpty() bool
}

func init() {
	rnd := rand.New(rand.NewSource(42))
	values = make([]int, numberSize)
	for i := range values {
		values[i] = rnd.Intn(1000)
	}
}

func TestBST(t *testing.T) {
	trees := map[string]BST{
		"coarse-grained": coarsegrained.NewTree(),
		"fine-grained":   finegrained.NewTree(),
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
			for _, v := range values {
				wg.Add(1)
				go func(x int) {
					defer wg.Done()
					tree.Insert(x)
				}(v)
			}
			wg.Wait()
			if !tree.IsValid() {
				t.Errorf("Insert only test failed: %s tree is not valid.", name)
			}
		})
	}
}
