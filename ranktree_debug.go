package ranktree

import "fmt"

func (node *TreeNode) print() {
	if node.low < node.high {
		fmt.Printf("NODE (%d, %d) %d\n", node.low, node.high, node.count)
		node.left.print()
		node.right.print()
	} else {
		fmt.Printf("LEAF (%d) %d %p\n",  node.low, node.count, node)
	}
}


func (node *TreeNode) printBackward() {
	for node != nil {
		fmt.Println(node.low, node.high)
		node = node.element.Next().Value.(*TreeNode)
	}
}

/*
func (tree *RankTree) Print() {
	tree.root.print()
}


func (node *TreeNode) Print() {
	if node == nil {
		fmt.Println("nil")
		return
	}
	fmt.Printf("NODE (%d, %d) %d\n", node.low, node.high, node.count)
	fmt.Println(node.members)
}
*/


