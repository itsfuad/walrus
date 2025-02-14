package utils

import (
	"testing"
)

const errMsg = "Has() = %v, want %v"

type Test struct {
	name     string
	slice    interface{}
	item     interface{}
	equalFn  interface{}
	expected bool
}

func TestHas(t *testing.T) {
	equalInt := func(a, b int) bool {
		return a == b
	}

	equalString := func(a, b string) bool {
		return a == b
	}

	equalStruct := func(a, b struct{ a, b int }) bool {
		return a == b
	}

	tests := []Test{
		{
			name:     "int slice contains item",
			slice:    []int{1, 2, 3, 4, 5},
			item:     3,
			equalFn:  equalInt,
			expected: true,
		},
		{
			name:     "int slice does not contain item",
			slice:    []int{1, 2, 3, 4, 5},
			item:     6,
			equalFn:  equalInt,
			expected: false,
		},
		{
			name:     "string slice contains item",
			slice:    []string{"a", "b", "c"},
			item:     "b",
			equalFn:  equalString,
			expected: true,
		},
		{
			name:     "string slice does not contain item",
			slice:    []string{"a", "b", "c"},
			item:     "d",
			equalFn:  equalString,
			expected: false,
		},
		{
			name:     "struct slice contains item",
			slice:    []struct{ a, b int }{{1, 2}, {3, 4}, {5, 6}},
			item:     struct{ a, b int }{3, 4},
			equalFn:  equalStruct,
			expected: true,
		},
		{
			name:     "struct slice does not contain item",
			slice:    []struct{ a, b int }{{1, 2}, {3, 4}, {5, 6}},
			item:     struct{ a, b int }{7, 8},
			equalFn:  equalStruct,
			expected: false,
		},
		{
			name:     "empty slice",
			slice:    []int{},
			item:     1,
			equalFn:  equalInt,
			expected: false,
		},
	}

	RunLoop(t, tests)

}

type SomeTest struct {
	name      string
	slice     interface{}
	predicate interface{}
	expected  bool
}

func TestSome(t *testing.T) {
	tests := []SomeTest{
		{
			name:      "int slice contains item greater than 3",
			slice:     []int{1, 2, 3, 4, 5},
			predicate: func(i int) bool { return i > 3 },
			expected:  true,
		},
		{
			name:      "int slice does not contain item greater than 5",
			slice:     []int{1, 2, 3, 4, 5},
			predicate: func(i int) bool { return i > 5 },
			expected:  false,
		},
		{
			name:      "string slice contains item with length 1",
			slice:     []string{"a", "bb", "ccc"},
			predicate: func(s string) bool { return len(s) == 1 },
			expected:  true,
		},
		{
			name:      "string slice does not contain item with length 4",
			slice:     []string{"a", "bb", "ccc"},
			predicate: func(s string) bool { return len(s) == 4 },
			expected:  false,
		},
		{
			name:      "struct slice contains item with a > 2",
			slice:     []struct{ a, b int }{{1, 2}, {3, 4}, {5, 6}},
			predicate: func(s struct{ a, b int }) bool { return s.a > 2 },
			expected:  true,
		},
		{
			name:      "struct slice does not contain item with a > 5",
			slice:     []struct{ a, b int }{{1, 2}, {3, 4}, {5, 6}},
			predicate: func(s struct{ a, b int }) bool { return s.a > 5 },
			expected:  false,
		},
		{
			name:      "empty slice",
			slice:     []int{},
			predicate: func(i int) bool { return i > 0 },
			expected:  false,
		},
	}

	RunSomeLoop(t, tests)
}

func RunSomeLoop(t *testing.T, tests []SomeTest) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runSomeTest(t, tt)
		})
	}
}

func runSomeTest(t *testing.T, tt SomeTest) {
	switch s := tt.slice.(type) {
	case []int:
		runSomeIntTest(t, s, tt)
	case []string:
		runSomeStringTest(t, s, tt)
	case []struct{ a, b int }:
		runSomeStructTest(t, s, tt)
	}
}

func runSomeIntTest(t *testing.T, s []int, tt SomeTest) {
	if _, got := Some(s, tt.predicate.(func(int) bool)); got != tt.expected {
		t.Errorf(errMsg, got, tt.expected)
	}
}

func runSomeStringTest(t *testing.T, s []string, tt SomeTest) {
	if _, got := Some(s, tt.predicate.(func(string) bool)); got != tt.expected {
		t.Errorf(errMsg, got, tt.expected)
	}
}

func runSomeStructTest(t *testing.T, s []struct{ a, b int }, tt SomeTest) {
	if _, got := Some(s, tt.predicate.(func(struct{ a, b int }) bool)); got != tt.expected {
		t.Errorf(errMsg, got, tt.expected)
	}
}

func RunLoop(t *testing.T, tests []Test) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch s := tt.slice.(type) {
			case []int:
				runIntTest(t, s, tt)
			case []string:
				runStringTest(t, s, tt)
			case []struct{ a, b int }:
				runStructTest(t, s, tt)
			}
		})
	}
}

func runIntTest(t *testing.T, s []int, tt Test) {
	if got := Has(s, tt.item.(int), tt.equalFn.(func(a, b int) bool)); got != tt.expected {
		t.Errorf(errMsg, got, tt.expected)
	}
}

func runStringTest(t *testing.T, s []string, tt Test) {
	if got := Has(s, tt.item.(string), tt.equalFn.(func(a, b string) bool)); got != tt.expected {
		t.Errorf(errMsg, got, tt.expected)
	}
}

func runStructTest(t *testing.T, s []struct{ a, b int }, tt Test) {
	if got := Has(s, tt.item.(struct{ a, b int }), tt.equalFn.(func(a, b struct{ a, b int }) bool)); got != tt.expected {
		t.Errorf(errMsg, got, tt.expected)
	}
}
