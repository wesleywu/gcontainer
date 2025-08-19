package g

import (
	"testing"

	"github.com/wesleywu/gcontainer/utils/comparators"
)

func TestTreeSet_Values(t *testing.T) {
	tests := []struct {
		name       string
		elements   []int
		comparator comparators.Comparator[int]
		expected   int
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
			name:       "Multiple elements",
			elements:   []int{1, 2, 3, 4, 5},
			comparator: comparators.ComparatorInt,
			expected:   5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewTreeSetFrom(tt.elements, tt.comparator)

			// Test Values method
			values := set.Values()
			if len(values) != tt.expected {
				t.Errorf("TreeSet.Values() returned %d values, want %d", len(values), tt.expected)
			}

			// Test that Values and Slice return the same result
			slice := set.Slice()
			if len(values) != len(slice) {
				t.Errorf("Values() length %d != Slice() length %d", len(values), len(slice))
			}

			// Test that all elements from the original slice are present
			valueSet := NewTreeSetFrom(values, tt.comparator)
			for _, elem := range tt.elements {
				if !valueSet.Contains(elem) {
					t.Errorf("Values() missing element %d", elem)
				}
			}
		})
	}
}

func TestTreeSet_Values_String(t *testing.T) {
	set := NewTreeSetFrom([]string{"a", "b", "c"}, comparators.ComparatorString)

	values := set.Values()
	expectedCount := 3

	if len(values) != expectedCount {
		t.Errorf("TreeSet.Values() returned %d values, want %d", len(values), expectedCount)
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

func TestTreeSet_Values_Concurrent(t *testing.T) {
	set := NewTreeSetFrom([]int{1, 2, 3, 4, 5}, comparators.ComparatorInt, true) // concurrent-safe

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

func TestTreeSet_Values_Consistency(t *testing.T) {
	set := NewTreeSetFrom([]int{3, 1, 4, 1, 5, 9, 2, 6}, comparators.ComparatorInt)

	// Test that multiple calls to Values return consistent results
	values1 := set.Values()
	values2 := set.Values()

	if len(values1) != len(values2) {
		t.Errorf("Values() calls returned different lengths: %d vs %d", len(values1), len(values2))
	}

	// Create sets from both results and compare
	set1 := NewTreeSetFrom(values1, comparators.ComparatorInt)
	set2 := NewTreeSetFrom(values2, comparators.ComparatorInt)

	if !set1.Equals(set2) {
		t.Errorf("Values() calls returned different results")
	}
}

func TestTreeSet_Values_Ordering(t *testing.T) {
	// Test that Values maintains the ordering defined by the comparator
	set := NewTreeSetFrom([]int{5, 2, 8, 1, 9}, comparators.ComparatorInt)

	values := set.Values()

	// Check that values are in ascending order
	for i := 1; i < len(values); i++ {
		if values[i] < values[i-1] {
			t.Errorf("Values not in ascending order: %d < %d at position %d", values[i], values[i-1], i)
		}
	}

	// Verify specific ordering
	expectedOrder := []int{1, 2, 5, 8, 9}
	for i, expected := range expectedOrder {
		if values[i] != expected {
			t.Errorf("Values[%d] = %d, want %d", i, values[i], expected)
		}
	}
}
