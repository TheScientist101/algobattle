package utils

import (
	"iter"
	"sync"
)

const (
	RED   = true
	BLACK = false
)

type node[T any] struct {
	value T
	left  *node[T]
	right *node[T]
	color bool
}

type Comparator[T any] func(a, b T) int

// TreeSet is a thread-safe Red-Black Tree implementation of a set
type TreeSet[T any] struct {
	mu         sync.RWMutex // Protects concurrent access
	root       *node[T]
	comparator Comparator[T]
}

// NewTreeSet creates a new TreeSet with the given comparator
func NewTreeSet[T any](comparator Comparator[T]) *TreeSet[T] {
	return &TreeSet[T]{
		comparator: comparator,
	}
}

// Insert adds values to the TreeSet
func (t *TreeSet[T]) Insert(values ...T) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, value := range values {
		t.root = t.insert(t.root, value)
		if t.root != nil {
			t.root.color = BLACK
		}
	}
}

// Remove removes values from the TreeSet
func (t *TreeSet[T]) Remove(values ...T) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, value := range values {
		t.root = t.delete(t.root, value)
		if t.root != nil {
			t.root.color = BLACK
		}
	}
}

// Contains checks if the TreeSet contains the given value
func (t *TreeSet[T]) Contains(value T) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	x := t.root
	for x != nil {
		if compareResult := t.comparator(value, x.value); compareResult < 0 {
			x = x.left
		} else if compareResult > 0 {
			x = x.right
		} else {
			return true
		}
	}
	return false
}

// All returns an iterator over all values in the TreeSet
func (t *TreeSet[T]) All() iter.Seq[T] {
	// Take a snapshot of the tree while holding the lock
	t.mu.RLock()
	values := t.AsSlice()
	t.mu.RUnlock()

	// Return an iterator over the snapshot
	return func(yield func(T) bool) {
		for _, v := range values {
			if !yield(v) {
				return
			}
		}
	}
}

// AsSlice returns all values in the TreeSet as a slice
func (t *TreeSet[T]) AsSlice() []T {
	t.mu.RLock()
	defer t.mu.RUnlock()

	slice := make([]T, 0)
	var traverse func(*node[T])
	traverse = func(n *node[T]) {
		if n == nil {
			return
		}
		traverse(n.left)
		slice = append(slice, n.value)
		traverse(n.right)
	}
	traverse(t.root)
	return slice
}

func (n *node[T]) isRed() bool {
	if n == nil {
		return false
	}
	return n.color == RED
}

func (n *node[T]) rotateLeft() *node[T] {
	x := n.right
	n.right = x.left
	x.left = n
	x.color = n.color
	n.color = RED
	return x
}

func (n *node[T]) rotateRight() *node[T] {
	x := n.left
	n.left = x.right
	x.right = n
	x.color = n.color
	n.color = RED
	return x
}

func (n *node[T]) flipColors() {
	n.color = RED
	if n.left != nil {
		n.left.color = BLACK
	}
	if n.right != nil {
		n.right.color = BLACK
	}
}

func (t *TreeSet[T]) insert(n *node[T], value T) *node[T] {
	if n == nil {
		return &node[T]{value: value, color: RED}
	}

	if compareResult := t.comparator(value, n.value); compareResult < 0 {
		n.left = t.insert(n.left, value)
	} else if t.comparator(value, n.value) > 0 {
		n.right = t.insert(n.right, value)
	}

	if n.right != nil && n.right.isRed() && (n.left == nil || !n.left.isRed()) {
		n = n.rotateLeft()
	}
	if n.left != nil && n.left.isRed() && n.left.left != nil && n.left.left.isRed() {
		n = n.rotateRight()
	}
	if n.left != nil && n.left.isRed() && n.right != nil && n.right.isRed() {
		n.flipColors()
	}

	return n
}

func (t *TreeSet[T]) delete(h *node[T], value T) *node[T] {
	if h == nil {
		return nil
	}

	if compareResult := t.comparator(value, h.value); compareResult < 0 {
		if (h.left == nil || !h.left.isRed()) && (h.left == nil || h.left.left == nil || !h.left.left.isRed()) {
			h = h.moveRedLeft()
		}
		if h.left != nil {
			h.left = t.delete(h.left, value)
		}
	} else {
		if h.left != nil && h.left.isRed() {
			h = h.rotateRight()
		}
		if t.comparator(value, h.value) == 0 && h.right == nil {
			return nil
		}
		if (h.right == nil || !h.right.isRed()) && (h.right == nil || h.right.left == nil || !h.right.left.isRed()) {
			h = h.moveRedRight()
		}
		if t.comparator(value, h.value) == 0 {
			smallest := h.right.min()
			h.value = smallest.value
			h.right = h.right.deleteMin()
		} else {
			if h.right != nil {
				h.right = t.delete(h.right, value)
			}
		}
	}
	return h.balance()
}

// min finds the minimum node in the subtree
func (n *node[T]) min() *node[T] {
	if n == nil {
		return nil // Handle empty subtree
	}
	if n.left == nil {
		return n
	}
	return n.left.min()
}

func (n *node[T]) deleteMin() *node[T] {
	if n == nil {
		return nil
	}
	if n.left == nil {
		return nil
	}

	if (n.left == nil || !n.left.isRed()) && (n.left.left == nil || !n.left.left.isRed()) {
		n = n.moveRedLeft()
	}

	n.left = n.left.deleteMin()
	return n.balance()
}

func (n *node[T]) moveRedLeft() *node[T] {
	n.flipColors()
	if n.right != nil && n.right.left != nil && n.right.left.isRed() {
		n.right = n.right.rotateRight()
		n = n.rotateLeft()
		n.flipColors()
	}
	return n
}

func (n *node[T]) moveRedRight() *node[T] {
	n.flipColors()
	if n.left != nil && n.left.left != nil && n.left.left.isRed() {
		n = n.rotateRight()
		n.flipColors()
	}
	return n
}

func (n *node[T]) balance() *node[T] {
	if n == nil {
		return nil
	}
	if n.right != nil && n.right.isRed() {
		n = n.rotateLeft()
	}

	if n.left != nil && n.left.isRed() && (n.left.left != nil && n.left.left.isRed()) {
		n = n.rotateRight()
	}

	if n.left != nil && n.left.isRed() && n.right != nil && n.right.isRed() {
		n.flipColors()
	}

	return n
}
