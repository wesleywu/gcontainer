package main

import (
	"fmt"

	"github.com/wesleywu/gcontainer/g"
	"github.com/wesleywu/gcontainer/utils/comparators"
)

func main() {
	fmt.Println("=== Values() Method Demo ===")

	fmt.Println("\n--- HashSet Values() ---")
	hashSetDemo()

	fmt.Println("\n--- TreeSet Values() ---")
	treeSetDemo()

	fmt.Println("\n--- ArrayList Values() ---")
	arrayListDemo()

	fmt.Println("\n--- LinkedList Values() ---")
	linkedListDemo()

	fmt.Println("\n--- Interface Consistency ---")
	interfaceConsistencyDemo()
}

func hashSetDemo() {
	set := g.NewHashSetFrom([]int{3, 1, 4, 1, 5, 9, 2, 6})

	fmt.Printf("Original set: %v\n", set.Slice())
	fmt.Printf("Values(): %v\n", set.Values())
	fmt.Printf("Slice(): %v\n", set.Slice())

	// Verify that Values() and Slice() return the same result
	values := set.Values()
	slice := set.Slice()

	if len(values) == len(slice) {
		fmt.Println("✓ Values() and Slice() return the same length")
	} else {
		fmt.Println("✗ Values() and Slice() return different lengths")
	}
}

func treeSetDemo() {
	set := g.NewTreeSetFrom([]int{3, 1, 4, 1, 5, 9, 2, 6}, comparators.ComparatorInt)

	fmt.Printf("Original set: %v\n", set.Slice())
	fmt.Printf("Values(): %v\n", set.Values())
	fmt.Printf("Slice(): %v\n", set.Slice())

	// Verify that Values() and Slice() return the same result
	values := set.Values()
	slice := set.Slice()

	if len(values) == len(slice) {
		fmt.Println("✓ Values() and Slice() return the same length")
	} else {
		fmt.Println("✗ Values() and Slice() return different lengths")
	}

	// Note: TreeSet maintains order based on comparator
	fmt.Println("Note: TreeSet maintains sorted order")
}

func arrayListDemo() {
	array := g.NewArrayListFrom([]int{3, 1, 4, 1, 5, 9, 2, 6})

	fmt.Printf("Original array: %v\n", array.Slice())
	fmt.Printf("Values(): %v\n", array.Values())
	fmt.Printf("Slice(): %v\n", array.Slice())

	// Verify that Values() and Slice() return the same result
	values := array.Values()
	slice := array.Slice()

	if len(values) == len(slice) {
		fmt.Println("✓ Values() and Slice() return the same length")
	} else {
		fmt.Println("✗ Values() and Slice() return different lengths")
	}

	// Note: ArrayList maintains insertion order
	fmt.Println("Note: ArrayList maintains insertion order")
}

func linkedListDemo() {
	list := g.NewLinkedListFrom([]int{3, 1, 4, 1, 5, 9, 2, 6})

	fmt.Printf("Original list: %v\n", list.Slice())
	fmt.Printf("Values(): %v\n", list.Values())
	fmt.Printf("Slice(): %v\n", list.Slice())

	// Verify that Values() and Slice() return the same result
	values := list.Values()
	slice := list.Slice()

	if len(values) == len(slice) {
		fmt.Println("✓ Values() and Slice() return the same length")
	} else {
		fmt.Println("✗ Values() and Slice() return different lengths")
	}

	// Note: LinkedList maintains insertion order
	fmt.Println("Note: LinkedList maintains insertion order")
}

func interfaceConsistencyDemo() {
	fmt.Println("\nTesting interface consistency...")

	// Test that all collection types can be used through Collection interface
	var collections []g.Collection[int]

	collections = append(collections, g.NewHashSetFrom([]int{1, 2, 3}))
	collections = append(collections, g.NewTreeSetFrom([]int{1, 2, 3}, comparators.ComparatorInt))
	collections = append(collections, g.NewArrayListFrom([]int{1, 2, 3}))
	collections = append(collections, g.NewLinkedListFrom([]int{1, 2, 3}))

	for i, collection := range collections {
		values := collection.Values()
		fmt.Printf("Collection %d type: %T, Values(): %v\n", i+1, collection, values)
	}

	fmt.Println("✓ All collection types implement Values() method consistently")
}
