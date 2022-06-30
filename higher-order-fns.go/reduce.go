package higher_order

func Reduce[T any, R any](arr []T, fn func(acc R, val T) R) R {
	var acc R

	for _, v := range arr {
		acc = fn(acc, v)
	}

	return acc
}
