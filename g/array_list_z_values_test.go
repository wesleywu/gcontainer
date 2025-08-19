package g

import (
	"testing"
)

func TestArrayList_Values(t *testing.T) {
	tests := []struct {
		name     string
		elements []int
		expected int
	}{
		{
			name:     "Empty array",
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
			array := NewArrayListFrom(tt.elements)

			// Test Values method
			values := array.Values()
			if len(values) != tt.expected {
				t.Errorf("ArrayList.Values() returned %d values, want %d", len(values), tt.expected)
			}

			// Test that Values and Slice return the same result
			slice := array.Slice()
			if len(values) != len(slice) {
				t.Errorf("Values() length %d != Slice() length %d", len(values), len(slice))
			}

			// Test that all elements from the original slice are present
			for i, elem := range tt.elements {
				if values[i] != elem {
					t.Errorf("Values()[%d] = %d, want %d", i, values[i], elem)
				}
			}
		})
	}
}

func TestArrayList_Values_String(t *testing.T) {
	array := NewArrayListFrom([]string{"a", "b", "c"})

	values := array.Values()
	expectedCount := 3

	if len(values) != expectedCount {
		t.Errorf("ArrayList.Values() returned %d values, want %d", len(values), expectedCount)
	}

	// Verify all expected strings are present
	expectedStrings := []string{"a", "b", "c"}
	for i, expected := range expectedStrings {
		if values[i] != expected {
			t.Errorf("Values()[%d] = %s, want %s", i, values[i], expected)
		}
	}
}

func TestArrayList_Values_Concurrent(t *testing.T) {
	array := NewArrayListFrom([]int{1, 2, 3, 4, 5}, true) // concurrent-safe

	// Test concurrent access
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			values := array.Values()
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

func TestArrayList_Values_Consistency(t *testing.T) {
	array := NewArrayListFrom([]int{3, 1, 4, 1, 5, 9, 2, 6})

	// Test that multiple calls to Values return consistent results
	values1 := array.Values()
	values2 := array.Values()

	if len(values1) != len(values2) {
		t.Errorf("Values() calls returned different lengths: %d vs %d", len(values1), len(values2))
	}

	// Verify both results are identical
	for i := range values1 {
		if values1[i] != values2[i] {
			t.Errorf("Values() calls returned different results at index %d: %d vs %d", i, values1[i], values2[i])
		}
	}
}
