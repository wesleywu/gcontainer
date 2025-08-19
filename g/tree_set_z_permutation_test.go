package g

import (
	"testing"

	"github.com/wesleywu/gcontainer/utils/comparators"
)

func TestTreeSet_Permutation(t *testing.T) {
	tests := []struct {
		name       string
		elements   []int
		comparator comparators.Comparator[int]
		expected   int // expected number of permutations
	}{
		{
			name:       "Empty set",
			elements:   []int{},
			comparator: comparators.ComparatorInt,
			expected:   0,
		},
		{
			name:       "Single element",
			elements:   []int{1},
			comparator: comparators.ComparatorInt,
			expected:   1,
		},
		{
			name:       "Two elements",
			elements:   []int{1, 2},
			comparator: comparators.ComparatorInt,
			expected:   2, // 2! = 2
		},
		{
			name:       "Three elements",
			elements:   []int{1, 2, 3},
			comparator: comparators.ComparatorInt,
			expected:   6, // 3! = 6
		},
		{
			name:       "Four elements",
			elements:   []int{1, 2, 3, 4},
			comparator: comparators.ComparatorInt,
			expected:   24, // 4! = 24
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewTreeSetFrom(tt.elements, tt.comparator)
			permutations := set.Permutation()

			if len(permutations) != tt.expected {
				t.Errorf("TreeSet.Permutation() returned %d permutations, want %d", len(permutations), tt.expected)
			}

			// Verify that all permutations are unique
			seen := make(map[string]bool)
			for _, perm := range permutations {
				permStr := ""
				for _, elem := range perm {
					permStr += string(rune(elem + '0'))
				}
				if seen[permStr] {
					t.Errorf("Duplicate permutation found: %v", perm)
				}
				seen[permStr] = true
			}

			// Verify that each permutation contains all elements
			for _, perm := range permutations {
				if len(perm) != len(tt.elements) {
					t.Errorf("Permutation length %d, want %d", len(perm), len(tt.elements))
				}

				// Check if all elements are present
				permSet := NewTreeSetFrom(perm, tt.comparator)
				for _, elem := range tt.elements {
					if !permSet.Contains(elem) {
						t.Errorf("Permutation %v missing element %d", perm, elem)
					}
				}
			}
		})
	}
}

func TestTreeSet_Permutation_String(t *testing.T) {
	set := NewTreeSetFrom([]string{"a", "b", "c"}, comparators.ComparatorString)
	permutations := set.Permutation()

	expectedCount := 6 // 3! = 6
	if len(permutations) != expectedCount {
		t.Errorf("TreeSet.Permutation() returned %d permutations, want %d", len(permutations), expectedCount)
	}

	// Verify specific permutations for string set
	expectedPerms := map[string]bool{
		"abc": true,
		"acb": true,
		"bac": true,
		"bca": true,
		"cab": true,
		"cba": true,
	}

	for _, perm := range permutations {
		permStr := ""
		for _, elem := range perm {
			permStr += elem
		}
		if !expectedPerms[permStr] {
			t.Errorf("Unexpected permutation: %s", permStr)
		}
	}
}

func TestTreeSet_Permutation_Concurrent(t *testing.T) {
	set := NewTreeSetFrom([]int{1, 2, 3, 4, 5}, comparators.ComparatorInt, true) // concurrent-safe

	// Test concurrent access
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			permutations := set.Permutation()
			if len(permutations) != 120 { // 5! = 120
				t.Errorf("Expected 120 permutations, got %d", len(permutations))
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestTreeSet_Permutation_DefaultComparator(t *testing.T) {
	set := NewTreeSetDefault[int]() // using default comparator
	set.Add(3, 1, 2)

	permutations := set.Permutation()
	expectedCount := 6 // 3! = 6

	if len(permutations) != expectedCount {
		t.Errorf("TreeSet.Permutation() returned %d permutations, want %d", len(permutations), expectedCount)
	}

	// Verify that all permutations are unique and contain all elements
	seen := make(map[string]bool)
	for _, perm := range permutations {
		permStr := ""
		for _, elem := range perm {
			permStr += string(rune(elem + '0'))
		}
		if seen[permStr] {
			t.Errorf("Duplicate permutation found: %v", perm)
		}
		seen[permStr] = true

		// Check if all elements are present
		if len(perm) != 3 {
			t.Errorf("Permutation length %d, want 3", len(perm))
		}
	}
}
