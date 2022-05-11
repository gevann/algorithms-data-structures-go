package trees

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTwoThreeNode(t *testing.T) {
	type args struct {
		init *int
	}
	var init *int = new(int)
	tests := []struct {
		name string
		args args
		want *twoThreeNode[int]
	}{
		{
			name: "Test 1",
			args: args{
				init: init,
			},
			want: &twoThreeNode[int]{
				firstData:   init,
				secondData:  nil,
				firstChild:  nil,
				secondChild: nil,
				thirdChild:  nil,
				parent:      nil,
				comparator:  intComparator,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TwoThreeNodeInt(tt.args.init)
			if got.firstData != tt.want.firstData {
				t.Errorf("TwoThreeNodeInt() = %v, want %v", got.firstData, tt.want.firstData)
			}
			if comparison := got.comparator(*got.firstData, *tt.want.firstData); comparison != 0 {
				t.Errorf("TwoThreeNodeInt.comparator(x, x) = %v, want 0", comparison)
			}
		})
	}
}

func Test_isLeaf(t *testing.T) {
	type args struct {
		node twoThreeNode[string]
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Returns true when the node has no children.",
			args: args{
				node: twoThreeNode[string]{
					firstData:   new(string),
					secondData:  nil,
					firstChild:  nil,
					secondChild: nil,
					thirdChild:  nil,
					parent:      nil,
					comparator:  stringComparator,
				},
			},
			want: true,
		},
		{
			name: "Returns false when the node has at least one child",
			args: args{
				node: twoThreeNode[string]{
					firstData:   new(string),
					secondData:  nil,
					firstChild:  &twoThreeNode[string]{},
					secondChild: nil,
					thirdChild:  nil,
					parent:      nil,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLeaf(tt.args.node); got != tt.want {
				t.Errorf("isLeaf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeType(t *testing.T) {
	type args struct {
		node twoThreeNode[int]
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Returns 3 when the node has three children and two datum.",
			args: args{
				node: twoThreeNode[int]{
					firstData:   new(int),
					secondData:  new(int),
					firstChild:  &twoThreeNode[int]{},
					secondChild: &twoThreeNode[int]{},
					thirdChild:  &twoThreeNode[int]{},
				},
			},
			want:    threeNode,
			wantErr: false,
		},
		{
			name: "Returns 2 when the node has two children and one datum.",
			args: args{
				node: twoThreeNode[int]{
					firstData:   new(int),
					secondData:  nil,
					firstChild:  &twoThreeNode[int]{},
					secondChild: &twoThreeNode[int]{},
					thirdChild:  nil,
				},
			},
			want:    twoNode,
			wantErr: false,
		},
		{
			name: "Errors when the node is a leaf",
			args: args{
				node: twoThreeNode[int]{
					firstData:   new(int),
					secondData:  nil,
					firstChild:  nil,
					secondChild: nil,
					thirdChild:  nil,
				},
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Errors in all other cases",
			args: args{
				node: twoThreeNode[int]{
					firstData:   new(int),
					secondData:  nil,
					firstChild:  nil,
					secondChild: &twoThreeNode[int]{},
					thirdChild:  nil,
				},
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nodeType(tt.args.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("nodeType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("nodeType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (node *twoThreeNode[int]) setFirstData(i int) *twoThreeNode[int] {
	node.firstData = &i
	return node
}

func (node *twoThreeNode[int]) setSecondData(i int) *twoThreeNode[int] {
	node.secondData = &i
	return node
}

func (node *twoThreeNode[int]) setFirstChild(child *twoThreeNode[int]) *twoThreeNode[int] {
	node.firstChild = child
	child.parent = node
	return node
}

func (node *twoThreeNode[int]) setSecondChild(child *twoThreeNode[int]) *twoThreeNode[int] {
	node.secondChild = child
	child.parent = node
	return node
}

func (node *twoThreeNode[int]) setThirdChild(child *twoThreeNode[int]) *twoThreeNode[int] {
	node.thirdChild = child
	child.parent = node
	return node
}

func twoThreeNodeInt() *twoThreeNode[int] {
	return &twoThreeNode[int]{
		comparator: intComparator,
	}
}

func (node *twoThreeNode[string]) equals(other *twoThreeNode[string]) bool {
	return node.firstData == other.firstData &&
		node.secondData == other.secondData &&
		node.firstChild == other.firstChild &&
		node.secondChild == other.secondChild &&
		node.thirdChild == other.thirdChild
}

func bfsEquals(root *twoThreeNode[int], other *twoThreeNode[int]) (bool, string) {
	queue := []*twoThreeNode[int]{root}
	otherQueue := []*twoThreeNode[int]{other}
	for len(queue) > 0 && len(otherQueue) > 0 {
		node := queue[0]
		queue = queue[1:]
		otherNode := otherQueue[0]
		otherQueue = otherQueue[1:]

		safeDeref := func(ptr *int) int {
			if ptr == nil {
				return 0
			}
			return *ptr
		}

		if !(safeDeref(node.firstData) == safeDeref(otherNode.firstData)) {
			return false, fmt.Sprintf("%v != %v", *node.firstData, *otherNode.firstData)
		}
		if !(safeDeref(node.secondData) == safeDeref(otherNode.secondData)) {
			return false, fmt.Sprintf("%v != %v", *node.secondData, *otherNode.secondData)
		}
		if node.firstChild != nil {
			queue = append(queue, node.firstChild)
			otherQueue = append(otherQueue, otherNode.firstChild)
		}
		if node.secondChild != nil {
			queue = append(queue, node.secondChild)
			otherQueue = append(otherQueue, otherNode.secondChild)
		}
		if node.thirdChild != nil {
			queue = append(queue, node.thirdChild)
			otherQueue = append(otherQueue, otherNode.thirdChild)
		}
	}
	return true, ""
}

func (node *twoThreeNode[int]) toString() string {
	var firstData int
	var secondData int
	if node.firstData != nil {
		firstData = *node.firstData
	}
	if node.secondData != nil {
		secondData = *node.secondData
	}

	return fmt.Sprintf("TwoThreeNode: [%v, %v]", firstData, secondData)
}

func Test_findLeaf(t *testing.T) {
	type args struct {
		value int
	}

	leafWithNine := twoThreeNodeInt().setFirstData(9)
	leafWithSix := twoThreeNodeInt().setFirstData(6)
	leafWithTwelve := twoThreeNodeInt().setFirstData(12)
	leafWithOneAndThree := twoThreeNodeInt().setFirstData(1).setSecondData(3)

	threeNodeWithSevenAndEleven := twoThreeNodeInt()
	threeNodeWithSevenAndEleven.setFirstData(7).setSecondData(11)
	threeNodeWithSevenAndEleven.setFirstChild(leafWithSix).setSecondChild(leafWithNine).setThirdChild(leafWithTwelve)

	tree := twoThreeNodeInt().
		setFirstData(4)
	tree.setFirstChild(leafWithOneAndThree).
		setSecondChild(threeNodeWithSevenAndEleven)

	tests := []struct {
		name string
		args args
		want *twoThreeNode[int]
	}{
		{
			name: "Returns the correct leaf",
			args: args{
				value: 10,
			},
			want: leafWithNine,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := findLeaf(*tree, tt.args.value); !(got.equals(tt.want)) {
				t.Errorf("findLeaf() = %v, want %v. Error: %v", got.toString(), tt.want.toString(), err)
			}
		})
	}
}

func Test_datumCount(t *testing.T) {
	type args struct {
		node *twoThreeNode[int]
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Returns the correct number of datums",
			args: args{
				node: twoThreeNodeInt().setFirstData(1).setSecondData(2),
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := datumCount(tt.args.node); got != tt.want {
				t.Errorf("datumCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_insertIntoSingleDatumNode(t *testing.T) {
	type args struct {
		node  *twoThreeNode[int]
		value int
	}
	tests := []struct {
		name string
		args args
		want *twoThreeNode[int]
	}{
		{
			name: "Inserts the value into the node",
			args: args{
				node:  twoThreeNodeInt().setFirstData(1),
				value: 2,
			},
			want: twoThreeNodeInt().setFirstData(1).setSecondData(2),
		},
		{
			name: "Reorders the data if the value being inserted is less than the first data",
			args: args{
				node:  twoThreeNodeInt().setFirstData(2),
				value: 1,
			},
			want: twoThreeNodeInt().setFirstData(1).setSecondData(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := insertIntoSingleDatumNode(tt.args.node, tt.args.value)
			gotFirstData, gotSecondData := *got.firstData, *got.secondData
			wantFirstData, wantSecondData := *tt.want.firstData, *tt.want.secondData

			if !(gotFirstData == wantFirstData && gotSecondData == wantSecondData) {
				t.Errorf("insertIntoSingleDatumNode() = %v, want %v", got.toString(), tt.want.toString())
			}
		})
	}
}

func Test_sortData(t *testing.T) {
	type args struct {
		node  *twoThreeNode[int]
		value int
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 int
		want2 int
	}{
		{
			name: "Returns [firstData, secondData, value] when value is greater than secondData",
			args: args{
				node:  twoThreeNodeInt().setFirstData(1).setSecondData(2),
				value: 3,
			},
			want:  1,
			want1: 2,
			want2: 3,
		},
		{
			name: "Returns [firstData, value, secondData] when value is greater than firstData and less than secondData",
			args: args{
				node:  twoThreeNodeInt().setFirstData(1).setSecondData(3),
				value: 2,
			},
			want:  1,
			want1: 2,
			want2: 3,
		},
		{
			name: "Returns [value, firstData, secondData] when value is less than firstData",
			args: args{
				node:  twoThreeNodeInt().setFirstData(1).setSecondData(2),
				value: 0,
			},
			want:  0,
			want1: 1,
			want2: 2,
		},
		{
			name: "Returns [value, firstData, secondData] when value equal to firstData",
			args: args{
				node:  twoThreeNodeInt().setFirstData(1).setSecondData(2),
				value: 1,
			},
			want:  1,
			want1: 1,
			want2: 2,
		},
		{
			name: "Returns [firstData, value, secondData] when value equal to secondData",
			args: args{
				node:  twoThreeNodeInt().setFirstData(1).setSecondData(2),
				value: 2,
			},
			want:  1,
			want1: 2,
			want2: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref1, ref2, ref3 := sortData(tt.args.node, tt.args.value)
			got, got1, got2 := *ref1, *ref2, *ref3
			if !reflect.DeepEqual([]int{got, got1, got2}, []int{tt.want, tt.want1, tt.want2}) {
				t.Errorf("sortData() = %v, want %v", []int{got, got1, got2}, []int{tt.want, tt.want1, tt.want2})
			}
		})
	}
}

func Test_insertIntoFullNode(t *testing.T) {
	leafWithOneAndFive := twoThreeNodeInt().setFirstData(1).setSecondData(5)
	leafWithNine := twoThreeNodeInt().setFirstData(9)
	_ = twoThreeNodeInt().setFirstData(7).setFirstChild(leafWithOneAndFive).setSecondChild(leafWithNine)

	expectedLeafWithFive := twoThreeNodeInt().setFirstData(5)
	expectedLeafWithOne := twoThreeNodeInt().setFirstData(1)
	expectedLeafWithNine := twoThreeNodeInt().setFirstData(9)
	expectedTree := twoThreeNodeInt().setFirstData(4).setSecondData(7).setFirstChild(expectedLeafWithOne).setSecondChild(expectedLeafWithFive).setThirdChild(expectedLeafWithNine)

	type args struct {
		value int
	}
	tests := []struct {
		name string
		args args
		want *twoThreeNode[int]
	}{
		{
			name: "Inserts the value into the tree",
			args: args{
				value: 4,
			},
			want: expectedTree,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := insertIntoFullNode(leafWithOneAndFive, tt.args.value)
			matched, message := bfsEquals(got, expectedTree)
			if !matched {
				t.Errorf("insertIntoFullNode() mismatch: %v", message)
				t.Errorf("got = %v", Print(got))
			}
		})
	}
}
