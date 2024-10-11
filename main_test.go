package main

import (
	"fmt"
	"testing"
)

func ReverseIterator[E any](arr []E) func(func(int, E) bool) {
	return func(yield func(int, E) bool ) {
		for i := len(arr)-1; i >= 0; i-- {
			if !yield(i, arr[i]) {
				break
			}
		}
	}
}

func TestIterator(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	for _, v := range ReverseIterator(arr) {
		fmt.Println(v)
	}
}

