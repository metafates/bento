package sliceutil

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
