package higher_order

func Map[T any, R any](arr []T, fn func(val T) R) []R {
	return Reduce(arr, func(acc []R, val T) []R {
		acc = append(acc, fn(val))
		return acc
	})
}
