package utils

// Has checks if a slice contains an item using a custom equality function.
func Has[T any](slice []T, item T, equalFn func(a, b T) bool) bool {
	for _, i := range slice {
		if equalFn(i, item) {
			return true
		}
	}
	return false
}

// Some checks if a slice contains an item using the default equality function.
func Some[T any](slice []T, predicate func(T) bool) bool {
	for _, i := range slice {
		if predicate(i) {
			return true
		}
	}
	return false
}

// None checks if a slice does not contain an item using the default equality function.
func None[T any](slice []T, predicate func(T) bool) bool {
	return !Some(slice, predicate)
}