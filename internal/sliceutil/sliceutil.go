package sliceutil

func Repeat[T any](value T, n int) []T {
	repeated := make([]T, 0, n)

	for i := 0; i < n; i++ {
		repeated = append(repeated, value)
	}

	return repeated
}

func Take[S ~[]T, T any](slice S, take int) S {
	take = min(take, len(slice))

	return slice[:take]
}

func Skip[S ~[]T, T any](slice S, skip int) S {
	if skip >= len(slice) {
		return nil
	}

	return slice[skip:]
}

// FlatMap manipulates a slice and transforms and flattens it to a slice of another type.
// The transform function can either return a slice or a `nil`, and in the `nil` case
// no value is added to the final slice.
func FlatMap[T any, R any](collection []T, iteratee func(item T) []R) []R {
	result := make([]R, 0, len(collection))

	for i := range collection {
		result = append(result, iteratee(collection[i])...)
	}

	return result
}
