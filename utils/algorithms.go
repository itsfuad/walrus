package utils

import "reflect"

func Has[T any](slice []T, item T) bool {
	for _, i := range slice {
		//use reflect.DeepEqual to compare the items
		if reflect.DeepEqual(i, item) {
			return true
		}
	}
	return false
}