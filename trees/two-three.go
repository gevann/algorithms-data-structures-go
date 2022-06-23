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

type TwoThreeNode[T any] struct {
	//Each node can have a maximum of two data values.
	firstData, secondData *T

	//Each node can have a maximum of three children.
	firstChild, secondChild, thirdChild *TwoThreeNode[T]

	//Each node has a parent, except for the root node.
	parent *TwoThreeNode[T]

	// comparator is used to compare two values.
	// returns -1 if a < b, 0 if a == b, 1 if a > b
	comparator func(T, T) int

	// the height of the tree.
	height int
}

// TwoThreeNodeInt is a constructor for a two-three tree with int values.
// It returns a root node of a two-three tree.
func TwoThreeNodeInt(init *int) *TwoThreeNode[int] {
	return &TwoThreeNode[int]{
		firstData:   init,
		secondData:  nil,
		firstChild:  nil,
		secondChild: nil,
		thirdChild:  nil,
		parent:      nil,
		comparator:  intComparator,
		height:      0,
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

func isLeaf[T any](node TwoThreeNode[T]) bool {
	return node.firstChild == nil && node.secondChild == nil && node.thirdChild == nil
}

const (
	twoNode = iota + 2
	threeNode
)

func nodeType[T any](node TwoThreeNode[T]) (int, error) {
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
func findLeaf[T any](node *TwoThreeNode[T], value T) (*TwoThreeNode[T], error) {
	if isLeaf(*node) {
		return node, nil
	}

	nt, err := nodeType(*node)

	if err != nil {
		return nil, err
	}

	switch nt {
	case twoNode:
		if node.comparator(value, *node.firstData) < 0 {
			return findLeaf(node.firstChild, value)
		} else {
			return findLeaf(node.secondChild, value)
		}
	case threeNode:
		if node.comparator(value, *node.firstData) < 0 {
			return findLeaf(node.firstChild, value)
		} else if node.comparator(value, *node.secondData) < 0 {
			return findLeaf(node.secondChild, value)
		} else {
			return findLeaf(node.thirdChild, value)
		}
	}

	return nil, errors.New("Unknown node type")
}

// datumCount returns the number of data values in the given twoThreeNode.
func datumCount[T any](node *TwoThreeNode[T]) int {
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
func sortData[T any](node *TwoThreeNode[T], value T) (*T, *T, *T) {
	if node.comparator(value, *node.firstData) <= 0 {
		return &value, node.firstData, node.secondData
	}
	if node.comparator(value, *node.secondData) <= 0 {
		return node.firstData, &value, node.secondData
	}
	return node.firstData, node.secondData, &value
}

func findRoot[T any](node *TwoThreeNode[T]) *TwoThreeNode[T] {
	if node.parent == nil {
		return node
	}
	node.height = 1 + maxHeight(node.firstChild, node.secondChild, node.thirdChild)
	return findRoot(node.parent)
}

// rebalance rebalances the tree after a node has been inserted.
// It recurses up the tree until it finds a node that is not full, or the root node.
// It returns the new root of the tree.
func rebalance[T any](node *TwoThreeNode[T], value T, tmpChildNode *TwoThreeNode[T]) *TwoThreeNode[T] {
	finalSplit := node.parent == nil || datumCount(node.parent) == 1

	if datumCount(node) == 1 {
		insertIntoSingleDatumNode(node, value)
		return findRoot(node)
	}

	// node cannot be split, so recurse up the tree
	min, mid, max := sortData(node, value)
	parent := node.parent

	// We can split this node.
	// it:
	// - has two data values
	// - is the root node, or
	// - its parent has only one data value
	node.firstData = min
	node.secondData = nil

	otherNode := TwoThreeNode[T]{
		firstData:   max,
		secondData:  nil,
		firstChild:  nil,
		secondChild: nil,
		thirdChild:  nil,
		parent:      parent,
		comparator:  node.comparator,
	}

	leftChildren, rightChildren := partitionChildNodes(*mid, []*TwoThreeNode[T]{node.firstChild, node.secondChild, node.thirdChild, tmpChildNode})

	node.firstChild = nil
	node.secondChild = nil
	node.thirdChild = nil

	for i, child := range leftChildren {
		child.parent = node
		if i == 0 {
			node.firstChild = child
		} else {
			node.secondChild = child
		}
	}

	for i, child := range rightChildren {
		child.parent = &otherNode
		if i == 0 {
			otherNode.firstChild = child
		} else {
			otherNode.thirdChild = child
		}
	}

	node.height = 1 + maxHeight(node.firstChild, node.secondChild, node.thirdChild)
	otherNode.height = 1 + maxHeight(otherNode.firstChild, otherNode.secondChild, otherNode.thirdChild)

	if parent == nil {
		parent := TwoThreeNode[T]{
			firstData:   mid,
			secondData:  nil,
			firstChild:  node,
			secondChild: &otherNode,
			thirdChild:  nil,
			parent:      nil,
			comparator:  node.comparator,
			height:      1 + maxHeight(node, &otherNode),
		}
		node.parent = &parent
		otherNode.parent = &parent
		return &parent
	}

	if finalSplit {
		if node.parent.firstChild == node {
			// node is the firstChild of its parent
			// - set otherNode as the secondChild of the parent and
			// - set secondChild of node as the thirdChild of the parent
			parent.thirdChild = parent.secondChild
			parent.secondChild = &otherNode
		} else {
			// node is the secondChild of its parent
			parent.thirdChild = &otherNode
			otherNode.parent = parent
		}
		return rebalance(parent, *mid, nil)
	} else {
		return rebalance(parent, *mid, &otherNode)
	}
}

func maxHeight[T any](nodes ...*TwoThreeNode[T]) int {
	max := 0

	if len(nodes) == 0 {
		return max
	}

	for _, node := range nodes {
		if node != nil && node.height > max {
			max = node.height
		}
	}

	return max
}

// insertIntoSingleDatumNode inserts a value into a single-datum node.
// It returns node after inserting the value.
func insertIntoSingleDatumNode[T any](node *TwoThreeNode[T], value T) *TwoThreeNode[T] {
	if node.comparator(value, *node.firstData) < 0 {
		node.secondData = node.firstData
		node.firstData = &value
	} else {
		node.secondData = &value
	}
	node.height = maxHeight(node.firstChild, node.secondChild) + 1
	return node
}

// partitionChildeNodes partitions the child nodes of a node into two groups based on the mid value given.
// It returns the two groups.
func partitionChildNodes[T any](midValue T, childNodes []*TwoThreeNode[T]) ([]*TwoThreeNode[T], []*TwoThreeNode[T]) {
	var leftChildNodes []*TwoThreeNode[T]
	var rightChildNodes []*TwoThreeNode[T]

	appendOrPrepend := func(childNode *TwoThreeNode[T], childNodes []*TwoThreeNode[T]) []*TwoThreeNode[T] {
		if len(childNodes) == 0 || childNodes[0].comparator(*childNodes[0].firstData, *childNode.firstData) < 0 {
			childNodes = append(childNodes, childNode)
		} else {
			childNodes = append([]*TwoThreeNode[T]{childNode}, childNodes...)
		}

		return childNodes
	}

	for _, childNode := range childNodes {
		if childNode == nil {
			continue
		}
		if childNode.comparator(*childNode.firstData, midValue) < 0 {
			leftChildNodes = appendOrPrepend(childNode, leftChildNodes)
		} else {
			rightChildNodes = appendOrPrepend(childNode, rightChildNodes)
		}
	}

	return leftChildNodes, rightChildNodes
}

// Insert inserts a value into the tree.
// Note that the root of the tree may be modified by this operation.
// It returns the root node of the tree.
func Insert[T any](root *TwoThreeNode[T], value T) (*TwoThreeNode[T], error) {
	node, err := findLeaf(root, value)
	if err != nil {
		goto EXIT_ERROR
	}

	return rebalance(node, value, nil), nil

EXIT_ERROR:
	return nil, err
}

type queueElement[T any] struct {
	node  *TwoThreeNode[T]
	level int
}

// reference converts the node's hex-string address to a base 62 string.
// It returns the base 62 string.
func reference[T any](node *TwoThreeNode[T]) string {
	hexAddress := fmt.Sprintf("%p", node)
	i, _ := strconv.ParseInt(hexAddress[2:], 16, 64)
	return big.NewInt(i).Text(62)
}

// ToString returns a string representation of the given twoThreeNode.
// It includes the reference to the parent node, if any, and the reference to itself if it is not a leaf.
// Returns a string where the references are base 62 numbers of the nodes' memory addresses.
func ToString[T any](node *TwoThreeNode[T]) string {
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
	return fmt.Sprintf("{%v[%v, %v] {h: %d} %v}\t", ref, fd, sd, node.height, parentRef)
}

// Print traverses the tree in breadth-first order.
// It returns the string representation of the tree.
func Print[T any](node *TwoThreeNode[T]) string {
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
