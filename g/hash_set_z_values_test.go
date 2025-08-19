package g

import (
	"testing"
)

func TestHashSet_Values(t *testing.T) {
	tests := []struct {
		name     string
		elements []int
		expected int
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
			name:     "Multiple elements",
			elements: []int{1, 2, 3, 4, 5},
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewHashSetFrom(tt.elements)

			// Test Values method
			values := set.Values()
			if len(values) != tt.expected {
				t.Errorf("HashSet.Values() returned %d values, want %d", len(values), tt.expected)
			}

			// Test that Values and Slice return the same result
			slice := set.Slice()
			if len(values) != len(slice) {
				t.Errorf("Values() length %d != Slice() length %d", len(values), len(slice))
			}

			// Test that all elements from the original slice are present
			valueSet := NewHashSetFrom(values)
			for _, elem := range tt.elements {
				if !valueSet.Contains(elem) {
					t.Errorf("Values() missing element %d", elem)
				}
			}
		})
	}
}

func TestHashSet_Values_String(t *testing.T) {
	set := NewHashSetFrom([]string{"a", "b", "c"})

	values := set.Values()
	expectedCount := 3

	if len(values) != expectedCount {
		t.Errorf("HashSet.Values() returned %d values, want %d", len(values), expectedCount)
	}

	// Verify all expected strings are present
	expectedStrings := map[string]bool{
		"a": true,
		"b": true,
		"c": true,
	}

	for _, value := range values {
		if !expectedStrings[value] {
			t.Errorf("Unexpected value: %s", value)
		}
	}
}

func TestHashSet_Values_Concurrent(t *testing.T) {
	set := NewHashSetFrom([]int{1, 2, 3, 4, 5}, true) // concurrent-safe

	// Test concurrent access
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			values := set.Values()
			if len(values) != 5 {
				t.Errorf("Expected 5 values, got %d", len(values))
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestHashSet_Values_Consistency(t *testing.T) {
	set := NewHashSetFrom([]int{3, 1, 4, 1, 5, 9, 2, 6})

	// Test that multiple calls to Values return consistent results
	values1 := set.Values()
	values2 := set.Values()

	if len(values1) != len(values2) {
		t.Errorf("Values() calls returned different lengths: %d vs %d", len(values1), len(values2))
	}

	// Create sets from both results and compare
	set1 := NewHashSetFrom(values1)
	set2 := NewHashSetFrom(values2)

	if !set1.Equals(set2) {
		t.Errorf("Values() calls returned different results")
	}
}
