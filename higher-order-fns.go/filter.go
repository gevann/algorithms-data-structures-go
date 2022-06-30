package higher_order

func Filter[T any](arr []T, fn func(val T) bool) []T {
	return Reduce(arr, func(acc []T, val T) []T {
		if fn(val) {
			return append(acc, val)
		}
		return acc
	})
}
