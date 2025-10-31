package main

func arraymap[T, R any](t []T, f func(T) R) []R {
	result := []R{}

	for _, i := range t {
		result = append(result, f(i))
	}

	return result
}

// reduce takes an array of T, and returns an array of T including only those members where b returns true
func reduce[T any](t []T, b func(T) bool) []T {
	reduced := []T{}
	for _, i := range t {
		if b(i) {
			reduced = append(reduced, i)
		}
	}

	return reduced
}

func mapreduce[T, R any](t []T, b func(T) bool, f func(T) R) []R {
	return arraymap(reduce(t, b), f)
}
