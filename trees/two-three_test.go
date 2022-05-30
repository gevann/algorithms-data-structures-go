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

func (node *twoThreeNode[int]) setFD(i int) *twoThreeNode[int] {
	node.firstData = &i
	return node
}

func (node *twoThreeNode[int]) setSD(i int) *twoThreeNode[int] {
	node.secondData = &i
	return node
}

func (node *twoThreeNode[int]) setFC(child *twoThreeNode[int]) *twoThreeNode[int] {
	node.firstChild = child
	child.parent = node
	return node
}

func (node *twoThreeNode[int]) setSC(child *twoThreeNode[int]) *twoThreeNode[int] {
	node.secondChild = child
	child.parent = node
	return node
}

func (node *twoThreeNode[int]) setTC(child *twoThreeNode[int]) *twoThreeNode[int] {
	node.thirdChild = child
	child.parent = node
	return node
}

func ttni() *twoThreeNode[int] {
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
	rootData := bfs(root)
	otherData := bfs(other)

	if len(rootData) != len(otherData) {
		return false, "Lengths don't match"
	}

	for idx, value := range rootData {
		if otherValue := otherData[idx]; value != otherValue {
			return false, fmt.Sprintf("Values at index %d don't match: %d != %d", idx, value, otherValue)
		}
	}

	return true, ""
}

func bfs(root *twoThreeNode[int]) []int {
	queue := []*twoThreeNode[int]{root}
	var result []int
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		result = append(result, *node.firstData)
		if node.firstChild != nil {
			queue = append(queue, node.firstChild)
		}
		if node.secondChild != nil {
			queue = append(queue, node.secondChild)
		}
		if node.thirdChild != nil {
			queue = append(queue, node.thirdChild)
		}
	}
	return result
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

	leafWithNine := ttni().setFD(9)
	leafWithSix := ttni().setFD(6)
	leafWithTwelve := ttni().setFD(12)
	leafWithOneAndThree := ttni().setFD(1).setSD(3)

	threeNodeWithSevenAndEleven := ttni()
	threeNodeWithSevenAndEleven.setFD(7).setSD(11)
	threeNodeWithSevenAndEleven.setFC(leafWithSix).setSC(leafWithNine).setTC(leafWithTwelve)

	tree := ttni().
		setFD(4)
	tree.setFC(leafWithOneAndThree).
		setSC(threeNodeWithSevenAndEleven)

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
			if got, err := findLeaf(tree, tt.args.value); !(got.equals(tt.want)) {
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
				node: ttni().setFD(1).setSD(2),
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
				node:  ttni().setFD(1),
				value: 2,
			},
			want: ttni().setFD(1).setSD(2),
		},
		{
			name: "Reorders the data if the value being inserted is less than the first data",
			args: args{
				node:  ttni().setFD(2),
				value: 1,
			},
			want: ttni().setFD(1).setSD(2),
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
				node:  ttni().setFD(1).setSD(2),
				value: 3,
			},
			want:  1,
			want1: 2,
			want2: 3,
		},
		{
			name: "Returns [firstData, value, secondData] when value is greater than firstData and less than secondData",
			args: args{
				node:  ttni().setFD(1).setSD(3),
				value: 2,
			},
			want:  1,
			want1: 2,
			want2: 3,
		},
		{
			name: "Returns [value, firstData, secondData] when value is less than firstData",
			args: args{
				node:  ttni().setFD(1).setSD(2),
				value: 0,
			},
			want:  0,
			want1: 1,
			want2: 2,
		},
		{
			name: "Returns [value, firstData, secondData] when value equal to firstData",
			args: args{
				node:  ttni().setFD(1).setSD(2),
				value: 1,
			},
			want:  1,
			want1: 1,
			want2: 2,
		},
		{
			name: "Returns [firstData, value, secondData] when value equal to secondData",
			args: args{
				node:  ttni().setFD(1).setSD(2),
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

func TestInsert(t *testing.T) {
	type args struct {
		root  *twoThreeNode[int]
		value int
	}
	tests := []struct {
		name    string
		args    args
		want    *twoThreeNode[int]
		wantErr bool
	}{
		{
			name: "It inserts into the node when it is the root and has room",
			args: args{
				root:  ttni().setFD(1),
				value: 2,
			},
			want:    ttni().setFD(1).setSD(2),
			wantErr: false,
		},
		{
			name: "It splits the node when it is the root and has no room",
			args: args{
				root:  ttni().setFD(1).setSD(2),
				value: 3,
			},
			want:    ttni().setFD(2).setFC(ttni().setFD(1)).setSC(ttni().setFD(3)),
			wantErr: false,
		},
		{
			name: "It splits the first child correctly, and returns the root",
			args: args{
				root:  ttni().setFD(10).setFC(ttni().setFD(5).setSD(7)).setSC(ttni().setFD(15)),
				value: 8,
			},
			want:    ttni().setFD(7).setSD(10).setFC(ttni().setFD(5)).setSC(ttni().setFD(8)).setTC(ttni().setFD(15)),
			wantErr: false,
		},
		{
			name: "It splits the second child correctly, and returns the root",
			args: args{
				root:  ttni().setFD(10).setFC(ttni().setFD(5).setSD(7)).setSC(ttni().setFD(15).setSD(20)),
				value: 25,
			},
			want:    ttni().setFD(10).setSD(20).setFC(ttni().setFD(5).setSD(7)).setSC(ttni().setFD(15)).setTC(ttni().setFD(25)),
			wantErr: false,
		},
		{

			name: "It splits the root correctly when the leaf and root were both full",
			args: args{
				root:  ttni().setFD(10).setSD(20).setFC(ttni().setFD(5).setSD(7)).setSC(ttni().setFD(15)).setTC(ttni().setFD(25)),
				value: 8,
			},
			want:    ttni().setFD(10).setFC(ttni().setFD(7).setFC(ttni().setFD(5)).setSC(ttni().setFD(8))).setSC(ttni().setFD(20).setFC(ttni().setFD(15)).setSC(ttni().setFD(25))),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Insert(tt.args.root, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			matched, message := bfsEquals(got, tt.want)

			if !matched {
				t.Errorf("Insert() mismatch: %v", message)
				t.Errorf("\nGOT:%v\n\nWANT:%v\n", Print(got), Print(tt.want))
			}

			if err != nil != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInsertMultiple(t *testing.T) {
	type args struct {
		root       *twoThreeNode[int]
		valuesList []int
	}

	treeRoot7 := ttni().setFD(7).setFC(ttni().setFD(5)).setSC(ttni().setFD(8))
	treeRoot17 := ttni().setFD(17).setFC(ttni().setFD(15)).setSC(ttni().setFD(20))
	treeRoot40 := ttni().setFD(40).setFC(ttni().setFD(35)).setSC(ttni().setFD(45))

	/*
			            (10, 25)
			  ----------------------------
		     /           |               \
			(7)         (17)             (40)
			/  \        /  \            /  \
		(5)    (8)    (15)    (20)    (35)    (45)
	*/
	threeLevelTree := ttni().setFD(10).setSD(25).setFC(treeRoot7).setSC(treeRoot17).setTC(treeRoot40)

	test := struct {
		name       string
		args       args
		want       *twoThreeNode[int]
		wantHeight int
	}{

		name: "It correctly inserts multiple values",
		args: args{
			root: ttni().setFD(10),
			valuesList: []int{
				20, 5, 7, 25, 35, 40, 45, 8, 15, 17,
			},
		},
		want:       threeLevelTree,
		wantHeight: 3,
	}
	t.Run(test.name, func(t *testing.T) {
		root := test.args.root

		for _, value := range test.args.valuesList {
			root, _ = Insert(root, value)
		}

		matched, message := bfsEquals(root, test.want)

		if !matched {
			t.Errorf("Insert() mismatch: %v", message)
			t.Errorf("\nGOT:%v\n\nWANT:%v\n", Print(root), Print(test.want))
		}

		matchedHeight := root.height == test.wantHeight
		if !matchedHeight {
			t.Errorf("Insert() height mismatch: %v", message)
			t.Errorf("\nGOT:%v\n\nWANT:%v\n", root.height, test.wantHeight)
		}
	})
}

func Test_partitionChildNodes(t *testing.T) {
	type args struct {
		midValue   int
		childNodes []*twoThreeNode[int]
	}
	tests := []struct {
		name  string
		args  args
		want  []*twoThreeNode[int]
		want1 []*twoThreeNode[int]
	}{
		{
			name: "It partitions the child nodes correctly",
			args: args{
				midValue:   5,
				childNodes: []*twoThreeNode[int]{ttni().setFD(11), ttni().setFD(2), ttni().setFD(7), ttni().setFD(3)},
			},
			want:  []*twoThreeNode[int]{ttni().setFD(2), ttni().setFD(3)},
			want1: []*twoThreeNode[int]{ttni().setFD(7), ttni().setFD(11)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := partitionChildNodes(tt.args.midValue, tt.args.childNodes)
			for i, node := range got {
				if !reflect.DeepEqual(*node.firstData, *tt.want[i].firstData) {
					t.Errorf("partitionChildNodes() got = %d, want %d", *node.firstData, *tt.want[i].firstData)
				}
			}
			for i, node := range got1 {
				if !reflect.DeepEqual(*node.firstData, *tt.want1[i].firstData) {
					t.Errorf("partitionChildNodes() got1 = %d, want %d", *node.firstData, *tt.want1[i].firstData)
				}
			}
		})
	}
}

func Test_maxHeight(t *testing.T) {
	type args struct {
		nodes []*twoThreeNode[int]
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "It returns the correct max height",
			args: args{
				nodes: []*twoThreeNode[int]{
					{
						height: 1,
					},
					{
						height: 2,
					},
					{
						height: 3,
					},
				},
			},
			want: 3,
		},
		{
			name: "It returns 0 if there are no nodes",
			args: args{
				nodes: []*twoThreeNode[int]{},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxHeight(tt.args.nodes...); got != tt.want {
				t.Errorf("maxHeight() = %v, want %v", got, tt.want)
			}
		})
	}
}
