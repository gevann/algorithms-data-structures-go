package trees

import (
	"errors"
)

/* TwoThreeTree is a binary search tree where each node has either:
    - two children (aka 2-node) and one data element, or
	- three children (aka 3-node) and two data elements
*/

type twoThreeNode[T any] struct {
	//Each node can have a maximum of two data values.
	firstData, secondData *T

	//Each node can have a maximum of three children.
	firstChild, secondChild, thirdChild *twoThreeNode[T]

	//Each node has a parent, except for the root node.
	parent *twoThreeNode[T]

	// comparator is used to compare two values.
	// returns -1 if a < b, 0 if a == b, 1 if a > b
	comparator func(T, T) int
}

// TwoThreeNodeInt is a constructor for a two-three tree with int values.
// It returns a root node of a two-three tree.
func TwoThreeNodeInt(init *int) *twoThreeNode[int] {
	return &twoThreeNode[int]{
		firstData:   init,
		secondData:  nil,
		firstChild:  nil,
		secondChild: nil,
		thirdChild:  nil,
		parent:      nil,
		comparator:  intComparator,
	}
}

func intComparator(a, b int) int {
	if a < b {
		return -1
	} else if a == b {
		return 0
	} else {
		return 1
	}
}

func stringComparator(a, b string) int {
	if a < b {
		return -1
	} else if a == b {
		return 0
	} else {
		return 1
	}
}

func isLeaf[T any](node twoThreeNode[T]) bool {
	return node.firstChild == nil && node.secondChild == nil && node.thirdChild == nil
}

const (
	twoNode = iota + 2
	threeNode
)

func nodeType[T any](node twoThreeNode[T]) (int, error) {
	var zeroVal int
	firstChild := node.firstChild != nil
	secondChild := node.secondChild != nil
	thirdChild := node.thirdChild != nil
	firstDatum := node.firstData != nil
	secondDatum := node.secondData != nil

	if isLeaf(node) {
		return zeroVal, errors.New("node is a leaf")
	}

	if firstChild && secondChild && !thirdChild && firstDatum && !secondDatum {
		return twoNode, nil
	} else if firstChild && secondChild && thirdChild && firstDatum && secondDatum {
		return threeNode, nil
	} else {
		return zeroVal, errors.New("node is not a valid two-three tree node")
	}
}

// findLeaf locates the leaf node of the given twoThreeTree within which the value should be inserted.
// It returns the leaf node at which the value should be inserted.
// If no leaf node exists, it returns nil.
func findLeaf[T any](node twoThreeNode[T], value T) (*twoThreeNode[T], error) {
	if isLeaf(node) {
		return &node, nil
	}

	nt, err := nodeType(node)

	if err != nil {
		return nil, err
	}

	switch nt {
	case twoNode:
		if node.comparator(value, *node.firstData) < 0 {
			return findLeaf(*node.firstChild, value)
		} else {
			return findLeaf(*node.secondChild, value)
		}
	case threeNode:
		if node.comparator(value, *node.firstData) < 0 {
			return findLeaf(*node.firstChild, value)
		} else if node.comparator(value, *node.secondData) < 0 {
			return findLeaf(*node.secondChild, value)
		} else {
			return findLeaf(*node.thirdChild, value)
		}
	}

	return nil, errors.New("Unknown node type")
}
