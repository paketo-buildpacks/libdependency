package collections

// FilterFunc returns an array containing only those elements of a for which f returns true
func FilterFunc[T any](a []T, f func(T) bool) []T {
	result := empty[T]()

	if f == nil {
		return result
	}

	for _, t := range a {
		if f(t) {
			result = append(result, t)
		}
	}

	return result
}

// TransformFunc returns an array containing the result of f for all elements of a
func TransformFunc[T, U any](a []T, f func(T) U) []U {
	result := empty[U]()

	if f == nil {
		return result
	}

	for _, t := range a {
		result = append(result, f(t))
	}

	return result
}

// TransformFuncWithError returns an array containing the result of f for all elements of a
// and allows f to return an error.
// If f returns an error at any point, an empty array and the error are returned.
func TransformFuncWithError[T, U any](a []T, f func(T) (U, error)) ([]U, error) {
	result := empty[U]()

	if f == nil {
		return result, nil
	}

	for _, t := range a {
		if u, err := f(t); err != nil {
			return empty[U](), err
		} else {
			result = append(result, u)
		}
	}

	return result, nil
}

func empty[T any]() []T {
	return make([]T, 0)
}
