package higher_order

import (
	"testing"
)

func testReduce[T any, R comparable](t *testing.T, name string, inputCollection []T, reducerFn func(acc R, val T) R, expectedResult R) {
	t.Run(name, func(t *testing.T) {
		if got := Reduce(inputCollection, reducerFn); got != expectedResult {
			t.Errorf("Reduce() = %v, want %v", got, expectedResult)
		}
	})
}

func TestReduce(t *testing.T) {
	testReduce(t, "Summing integers", []int{1, 2, 3, 4, 5}, func(acc, val int) int {
		return acc + val
	}, 15)

	testReduce(t, "Concatenating strings", []string{"a", "b", "c", "d", "e"}, func(acc, val string) string {
		return acc + val
	}, "abcde")

	testReduce(t, "Maintaining booleans", []int{1, 2, 3, 4, 5}, func(acc bool, val int) bool {
		return acc && val <= 3
	}, false)
}
