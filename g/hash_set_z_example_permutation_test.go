package g

import (
	"fmt"
	"sort"
)

func ExampleHashSet_Permutation() {
	// Create a new HashSet with some elements
	set := NewHashSetFrom([]int{1, 2, 3})

	// Get all permutations
	permutations := set.Permutation()

	// Sort the set elements for consistent output
	elements := set.Slice()
	sort.Ints(elements)

	fmt.Printf("Set: %v\n", elements)
	fmt.Printf("Number of permutations: %d\n", len(permutations))
	fmt.Println("All permutations:")

	// Sort permutations for consistent output
	sort.Slice(permutations, func(i, j int) bool {
		for k := 0; k < len(permutations[i]); k++ {
			if k >= len(permutations[j]) {
				return false
			}
			if permutations[i][k] != permutations[j][k] {
				return permutations[i][k] < permutations[j][k]
			}
		}
		return len(permutations[i]) < len(permutations[j])
	})

	for i, perm := range permutations {
		fmt.Printf("  %d: %v\n", i+1, perm)
	}

	// Output:
	// Set: [1 2 3]
	// Number of permutations: 6
	// All permutations:
	//   1: [1 2 3]
	//   2: [1 3 2]
	//   3: [2 1 3]
	//   4: [2 3 1]
	//   5: [3 1 2]
	//   6: [3 2 1]
}
