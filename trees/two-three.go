package trees

/* TwoThreeTree is a binary search tree with the following properties:
 * 1. The left subtree of a node contains only nodes with keys less than the node's key.
 * 2. The right subtree of a node contains only nodes with keys greater than the node's key.
 * 3. The left and right subtrees of a node do not contain duplicate keys.
 */

type twoThreeNode[T any] struct {
	//Each node can have a maximum of two data values.
	firstData, secondData *T

	//Each node can have a maximum of three children.
	firstChild, secondChild, thirdChild *twoThreeNode[T]

	//Each node has a parent, except for the root node.
	parent *twoThreeNode[T]
}

/**
Returns a root node of a two-three tree.
*/
func TwoThreeNode[T any](init *T) *twoThreeNode[T] {
	return &twoThreeNode[T]{
		firstData:   init,
		secondData:  nil,
		firstChild:  nil,
		secondChild: nil,
		thirdChild:  nil,
		parent:      nil,
	}
}

func isLeaf[T any](node twoThreeNode[T]) bool {
	return node.firstChild == nil && node.secondChild == nil && node.thirdChild == nil
}
