package ranktree

import (
	"log"
	"fmt"
)

func Example() {
	// create
	tree, err := New(0, 10000)
	if err != nil {
		log.Fatal(err)
	}

	// add
	tree.Add("Alice", 123)
	tree.Add("Bob", 1234)
	tree.Add("Charles", 12)

	// get rank
	n := tree.RevRank("Bob")
	fmt.Printf("Bob is No.%d\n", n + 1)

	// get range
	result := tree.RevRange(0, -1)
	for i, name := range result {
		fmt.Printf("No.%d: %s\n", i + 1, name)
	}

	// Output:
	// Bob is No.1
	//
	// No.1: Bob
	// No.2: Alice
	// No.3: Charles
}


