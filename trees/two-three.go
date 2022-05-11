package trees

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
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

// datumCount returns the number of data values in the given twoThreeNode.
func datumCount[T any](node *twoThreeNode[T]) int {
	count := 0
	for _, datum := range []*T{node.firstData, node.secondData} {
		if datum != nil {
			count++
		}
	}
	return count
}

// sortData returns the data of a node and a given value, sorted in ascending order.
// It returns the sorted data
func sortData[T any](node *twoThreeNode[T], value T) (*T, *T, *T) {
	if node.comparator(value, *node.firstData) <= 0 {
		return &value, node.firstData, node.secondData
	}
	if node.comparator(value, *node.secondData) <= 0 {
		return node.firstData, &value, node.secondData
	}
	return node.firstData, node.secondData, &value
}

// insertIntoFullNode inserts a value into a full node, when the parent node is not full.
// It returns the parent node.
func insertIntoFullNode[T any](node *twoThreeNode[T], value T) *twoThreeNode[T] {
	min, mid, max := sortData(node, value)
	parent := node.parent

	insertIntoSingleDatumNode(parent, *mid)

	node.firstData = min
	node.secondData = nil

	nodePrime := twoThreeNode[T]{
		firstData:   max,
		secondData:  nil,
		firstChild:  nil,
		secondChild: nil,
		thirdChild:  nil,
		parent:      parent,
		comparator:  node.comparator,
	}

	// if node is the firstChild of its parent:
	// - set nodePrime as the secondChild of the parent and
	// - set secondChild of node as the thirdChild of the parent
	if node.parent.firstChild == node {
		secondChild := node.parent.secondChild
		node.parent.thirdChild = secondChild
		node.parent.secondChild = &nodePrime
	}

	return parent
}

// insertIntoSingleDatumNode inserts a value into a single-datum node.
// It returns node after inserting the value.
func insertIntoSingleDatumNode[T any](node *twoThreeNode[T], value T) *twoThreeNode[T] {
	if node.comparator(value, *node.firstData) < 0 {
		node.secondData = node.firstData
		node.firstData = &value
	} else {
		node.secondData = &value
	}
	return node
}

// Insert inserts a value into the tree.
// Note that the root of the tree may be modified by this operation.
// It returns the root node of the tree.
func Insert[T any](root *twoThreeNode[T], value T) (*twoThreeNode[T], error) {
	node, err := findLeaf(*root, value)
	if err != nil {
		goto EXIT_ERROR
	}

	if datumCount(node) == 1 {
		insertIntoSingleDatumNode(node, value)
		return root, nil
	}
	return node, nil
EXIT_ERROR:
	return nil, err
}

type queueElement[T any] struct {
	node  *twoThreeNode[T]
	level int
}

// reference converts the node's hex-string address to a base 62 string.
// It returns the base 62 string.
func reference[T any](node *twoThreeNode[T]) string {
	hexAddress := fmt.Sprintf("%p", node)
	i, _ := strconv.ParseInt(hexAddress[2:], 16, 64)
	return big.NewInt(i).Text(62)
}

// ToString returns a string representation of the given twoThreeNode.
// It includes the reference to the parent node, if any, and the reference to itself if it is not a leaf.
// Returns a string where the references are base 62 numbers of the nodes' memory addresses.
func ToString[T any](node *twoThreeNode[T]) string {
	fd, sd, parentRef, ref := "_", "_", "", ""

	if !isLeaf(*node) {
		ref = fmt.Sprintf("@<%v>", reference(node))
	}

	if node.parent != nil {
		parentRef = fmt.Sprintf(" (parent: %s)", reference(node.parent))
	}

	if node.firstData != nil {
		fd = fmt.Sprintf("%v", *node.firstData)
	}
	if node.secondData != nil {
		sd = fmt.Sprintf("%v", *node.secondData)
	}
	return fmt.Sprintf("{%v[%v, %v]%v}\t", ref, fd, sd, parentRef)
}

// Print traverses the tree in breadth-first order.
// It returns the string representation of the tree.
func Print[T any](node *twoThreeNode[T]) string {
	var queue []queueElement[T]
	var str string = "\n0:\t"
	queue = append(queue, queueElement[T]{node, 0})
	highestLevel := 0

	dataStr := func(elem queueElement[T]) string {
		start := ""
		if elem.level > highestLevel {
			highestLevel = elem.level
			start = fmt.Sprintf("\n%d:\t", elem.level)
		}
		return fmt.Sprintf("%v%s\t", start, ToString(elem.node))
	}

	for len(queue) > 0 {
		elem := queue[0]
		queue = queue[1:]
		str += fmt.Sprintf("%v", dataStr(elem))
		nextLevel := elem.level + 1
		if elem.node.firstChild != nil {
			queue = append(queue, queueElement[T]{elem.node.firstChild, nextLevel})
		}
		if elem.node.secondChild != nil {
			queue = append(queue, queueElement[T]{elem.node.secondChild, nextLevel})
		}
		if elem.node.thirdChild != nil {
			queue = append(queue, queueElement[T]{elem.node.thirdChild, nextLevel})
		}
	}
	return str
}
