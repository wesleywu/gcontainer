// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/Agogf/gf.

package g_test

import (
	"fmt"

	"github.com/wesleywu/gcontainer/g"
	"github.com/wesleywu/gcontainer/utils/comparators"
	"github.com/wesleywu/gcontainer/utils/gconv"
)

func ExampleNewAVLTree() {
	avlTree := g.NewAVLTree[string, string](comparators.ComparatorString)
	for i := 0; i < 6; i++ {
		avlTree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	fmt.Println(avlTree)

	// Output:
	// │       ┌── key5
	// │   ┌── key4
	// └── key3
	//     │   ┌── key2
	//     └── key1
	//         └── key0
}

func ExampleNewAVLTreeFrom() {
	avlTree := g.NewAVLTree[string, string](comparators.ComparatorString)
	for i := 0; i < 6; i++ {
		avlTree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	otherAvlTree := g.NewAVLTreeFrom(comparators.ComparatorString, avlTree.Map())
	fmt.Println(otherAvlTree)

	// May Output:
	// │   ┌── key5
	// │   │   └── key4
	// └── key3
	//     │   ┌── key2
	//     └── key1
	//         └── key0
}

func ExampleNewBTree() {
	bTree := g.NewBTree[string, string](3, comparators.ComparatorString)
	for i := 0; i < 6; i++ {
		bTree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}
	fmt.Println(bTree.Map())

	// Output:
	// map[key0:val0 key1:val1 key2:val2 key3:val3 key4:val4 key5:val5]
}

func ExampleNewBTreeFrom() {
	bTree := g.NewBTree[string, string](3, comparators.ComparatorString)
	for i := 0; i < 6; i++ {
		bTree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}

	otherBTree := g.NewBTreeFrom(3, comparators.ComparatorString, bTree.Map())
	fmt.Println(otherBTree.Map())

	// Output:
	// map[key0:val0 key1:val1 key2:val2 key3:val3 key4:val4 key5:val5]
}

func ExampleNewTreeMap() {
	tree := g.NewTreeMap[string, string](comparators.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}
	fmt.Println(tree)
	// Output:
	// │           ┌── key5
	// │       ┌── key4
	// │   ┌── key3
	// │   │   └── key2
	// └── key1
	//     └── key0
}

func ExampleNewTreeMapFrom() {
	tree := g.NewTreeMap[string, string](comparators.ComparatorString)
	for i := 0; i < 6; i++ {
		tree.Put("key"+gconv.String(i), "val"+gconv.String(i))
	}
	otherTree := g.NewTreeMapFrom[string, string](comparators.ComparatorString, tree.Map())
	fmt.Println(otherTree)
	// May Output:
	// │           ┌── key5
	// │       ┌── key4
	// │   ┌── key3
	// │   │   └── key2
	// └── key1
	//     └── key0
}
