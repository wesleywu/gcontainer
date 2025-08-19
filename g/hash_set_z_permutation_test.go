package g

import (
	"testing"
)

func TestHashSet_Permutation(t *testing.T) {
	tests := []struct {
		name     string
		elements []int
		expected int // expected number of permutations
	}{
		{
			name:     "Empty set",
			elements: []int{},
			expected: 0,
		},
		{
			name:     "Single element",
			elements: []int{1},
			expected: 1,
		},
		{
			name:     "Two elements",
			elements: []int{1, 2},
			expected: 2, // 2! = 2
		},
		{
			name:     "Three elements",
			elements: []int{1, 2, 3},
			expected: 6, // 3! = 6
		},
		{
			name:     "Four elements",
			elements: []int{1, 2, 3, 4},
			expected: 24, // 4! = 24
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewHashSetFrom(tt.elements)
			permutations := set.Permutation()

			if len(permutations) != tt.expected {
				t.Errorf("HashSet.Permutation() returned %d permutations, want %d", len(permutations), tt.expected)
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
				permSet := NewHashSetFrom(perm)
				for _, elem := range tt.elements {
					if !permSet.Contains(elem) {
						t.Errorf("Permutation %v missing element %d", perm, elem)
					}
				}
			}
		})
	}
}

func TestHashSet_Permutation_String(t *testing.T) {
	set := NewHashSetFrom([]string{"a", "b", "c"})
	permutations := set.Permutation()

	expectedCount := 6 // 3! = 6
	if len(permutations) != expectedCount {
		t.Errorf("HashSet.Permutation() returned %d permutations, want %d", len(permutations), expectedCount)
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

func TestHashSet_Permutation_Concurrent(t *testing.T) {
	set := NewHashSetFrom([]int{1, 2, 3, 4, 5}, true) // concurrent-safe

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
