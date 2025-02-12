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
func Some[T any](slice []T, predicate func(T) bool) (T, bool) {
	for _, item := range slice {
		if predicate(item) {
			return item, true
		}
	}
	var zeroValue T
	return zeroValue, false
}

// None checks if a slice does not contain an item using the default equality function.
func None[T any](slice []T, predicate func(T) bool) bool {
	if _, found := Some(slice, predicate); !found {
		return true
	} else {
		return false
	}
}
