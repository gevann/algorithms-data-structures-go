package trees

import (
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TwoThreeNode(tt.args.init); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TwoThreeNode() = %v, want %v", got, tt.want)
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
