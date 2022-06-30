package higher_order

import (
	"reflect"
	"strings"
	"testing"
)

func testMap[T any, R comparable](t *testing.T, name string, inputCollection []T, mapperFn func(val T) R, expectedResult []R) {
	t.Run(name, func(t *testing.T) {
		if got := Map(inputCollection, mapperFn); !reflect.DeepEqual(got, expectedResult) {
			t.Errorf("Reduce() = %v, want %v", got, expectedResult)
		}
	})
}

func TestMap(t *testing.T) {
	testMap(t, "Doubling integers", []int{1, 2, 3, 4, 5}, func(val int) int {
		return val * 2
	}, []int{2, 4, 6, 8, 10})

	testMap(t, "Capitalizing strings", []string{"a", "b", "c", "d", "e"}, func(val string) string {
		return strings.ToUpper(val)
	}, []string{"A", "B", "C", "D", "E"})

	testMap(t, "Flipping booleans", []bool{true, true, false, false}, func(val bool) bool {
		return !val
	}, []bool{false, false, true, true})
}
