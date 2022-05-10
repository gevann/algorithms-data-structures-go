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
