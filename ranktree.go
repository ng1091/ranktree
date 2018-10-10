// RankTree is a data structure used for Ranking.
//
// Features
//
// High performance
//
// Rank, Range, Pop, Update, Remove etc.
//
// Restrictions: the score must be non-negative integers.
//
package ranktree

import (
	"fmt"
	"log"
	"errors"
	"sort"

	"github.com/ng1091/ranktree/list"
)



// TreeNode is an element of a RankTree.
// If low < high, it's an internal node, if low = high, it's a leaf node.
type TreeNode struct {
	low    		int			// lower bound of the score range
	high   		int			// upper bound of the score range
	count  		int			// number of children
	left		*TreeNode	// left child
	right   	*TreeNode	// right child
	parent		*TreeNode

	element		*list.Element	// point to list.Element
	members 	[]string
}


// RankTree is a rank data structure based on binary tree.
type RankTree struct {
	root		*TreeNode
	nodeMap 	map[string]*TreeNode	// member to node
	list		*list.List				// singly linked list
	count   	int						// number of members

	minScore	int
	maxScore	int
	// usedMemory uint
}


// Rank result generator.
type rankResultGenerator struct {
	node *TreeNode
	index int
	reverse bool
}


// Rank result.
type RankWithScore struct {
	Member string
	Score int
}


// New Creates a RankTree.
// Low and high represents the score range.
func New(low int, high int) (*RankTree, error) {
	// check range
	if low < 0 || high < 0 {
		return nil, errors.New("low, high must be non-negative")
	}

	if low > high {
		return nil, errors.New("low less than high")
	}

	tree := new(RankTree)
	tree.root = new(TreeNode)
	tree.root.create(low, high, nil)
	tree.nodeMap = make(map[string]*TreeNode)
	tree.list = list.New()
	tree.minScore = low
	tree.maxScore = high
	return tree, nil
}


// Add adds a member to RankTree.
// If <member> exists, or <score> out of the range, false returned.
func (tree *RankTree) Add(member string, score int) bool {
	// member not in nodeMap
	if _, ok := tree.nodeMap[member]; ok == false {
		node := tree.find(score)

		if node != nil {
			// create a list element
			if node.element == nil {
				tree.createListElement(node)
			}

			tree.nodeMap[member] = node
			tree.count++

			node.members = append(node.members, member)
			if len(node.members) > 1 {
				sort.Strings(node.members)
			}

			node.incrementCount(1)
			return true
		}
	}
	return false
}


// Returns the rank of the member in RankTree.
// If member does not exist, -1 returned.
// Scores ordered from low to high.
// The rank is 0-based, which means that the member with the lowest score has rank 0.
// Use RevRank() to get the rank of an element with the scores ordered from high to low.
func (tree *RankTree) Rank(member string) int {
	if node, ok := tree.nodeMap[member]; ok == true {
		// offset in node.members
		offset := 0
		for k, v := range node.members {
			if v == member {
				offset = k
			}
		}

		sum := 0
		for node.parent != nil {
			thisNode := node
			node = node.parent
			if node.left != thisNode {
				sum += node.left.count
			}
		}
		return sum + offset
	}
	return -1
}



// Returns the rank of the member in RankTree
// If member does not exist, -1 returned.
// Scores ordered from high to low.
// The rank is 0-based, which means that the member with the highest score has rank 0.
// Use Rank() to get the rank of an element with the scores ordered from low to high.
func (tree *RankTree) RevRank(member string) int {
	if node, ok := tree.nodeMap[member]; ok == true {
		// offset in node.members
		offset := 0
		for k, v := range node.members {
			if v == member {
				offset = node.count - k - 1
			}
		}

		sum := 0
		for node.parent != nil {
			thisNode := node
			node = node.parent
			if node.right != thisNode {
				sum += node.right.count
			}
		}
		return sum + offset
	}
	return -1
}


// Returns the cardinality (number of members) of the RankTree.
func (tree *RankTree) Card() int {
	return tree.count
}


// Returns the score of member in the RankTree.
// If member does not exist in the RankTree, -1 is returned.
func (tree *RankTree) Score(member string) int {
	if node, ok := tree.nodeMap[member]; ok == true {
		return node.low
	}
	return -1
}


// Returns the number of members in the RankTree with a score between min and max.
func (tree *RankTree) Count(min, max int) int {
	if min < tree.minScore {
		min = tree.minScore
	}

	if max > tree.maxScore {
		max = tree.maxScore
	}

	if min > max || max < tree.minScore || min > tree.maxScore {
		return 0
	}

	leftNode := tree.find(min)
	rightNode := tree.find(max)
	leftCount := leftNode.countLeftArea()
	rightCount := rightNode.countLeftArea() + rightNode.count
	return rightCount - leftCount
}


// Removes <members> from the RankTree.
// Returns the number of members removed from the RankTree.
func (tree *RankTree) Remove(members ...string) (sum int) {
	for _, member := range members {
		sum += tree.remove(member)
	}
	return
}


// Removes <member> from the RankTree.
// If <member> exists, 1 is returned, otherwise 0 is returned.
func (tree *RankTree) remove(member string) int {
	if node, ok := tree.nodeMap[member]; ok == true {
		// remove member from node.members
		for i, v := range node.members {
			if v == member {
				node.members = append(node.members[:i], node.members[i+1:]...)
				break
			}
		}
		// remove list element
		if node.count == 1 {
			greaterNode := tree.findNextGreaterElement(node)
			tree.list.RemoveNext(greaterNode)
			node.element = nil
		}
		// remove map & count
		delete(tree.nodeMap, member)
		tree.count--
		node.incrementCount(-1)
		return 1
	}
	return 0
}


// Removes and returns a member with the highest score in the RankTree.
func (tree *RankTree) PopMax() (rank *RankWithScore) {
	if tree.count == 0 {
		return nil
	}

	e := tree.list.Head()
	if node, ok :=  e.Value.(*TreeNode); ok {
		rank = new(RankWithScore)
		member := node.members[len(node.members) - 1]
		tree.remove(member)
		rank.Member = member
		rank.Score = node.low
	}
	return
}


// Removes and returns a member with the lowest score in the RankTree.
func (tree *RankTree) PopMin() (rank *RankWithScore) {
	if tree.count == 0 {
		return nil
	}

	e := tree.list.Back()
	if node, ok :=  e.Value.(*TreeNode); ok {
		rank = new(RankWithScore)
		member := node.members[len(node.members) - 1]
		tree.remove(member)
		rank.Member = member
		rank.Score = node.low
	}
	return
}


// Removes and returns up to <n> members with the highest scores in the RankTree.
func (tree *RankTree) PopMaxN(n int) (ranks []RankWithScore) {
	if n < 0 {
		n = 0
	}

	if tree.count < n {
		n = tree.count
	}

	ranks = make([]RankWithScore, n)
	for i := 0; i < n; i++ {
		ranks[i] = *tree.PopMax()
	}
	return
}


// Removes and returns up to <n> members with the lowest scores in the RankTree.
func (tree *RankTree) PopMinN(n int) (ranks []RankWithScore) {
	if n < 0 {
		n = 0
	}

	if tree.count < n {
		n = tree.count
	}

	ranks = make([]RankWithScore, n)
	for i := 0; i < n; i++ {
		ranks[i] = *tree.PopMin()
	}
	return
}


// Increments the score of member in the RankTree.
// If <member> does not exist in the RankTree, it is added with <score>.
// Returns the new score of the member.
func (tree *RankTree) IncrementBy(member string, score int) int {
	currentScore := score
	if node, ok := tree.nodeMap[member]; ok == true {
		currentScore += node.low
		tree.remove(member)
	}
	if tree.Add(member, currentScore) {
		return currentScore
	} else {
		return -1
	}
}


// Updates the score of <member> in the RankTree.
// If <insert> is true, a new member is added when it does not exist in the RankTree.
// Returns a bool represents whether the update is successful or not.
func (tree *RankTree) UpdateScore(member string, score int, insert bool) bool {
	n := tree.remove(member)

	if n > 0 || insert {
		tree.Add(member, score)
		return true
	} else {
		return false
	}
}


// Sanitize indexes of rangeBasic(), rangeWithScore().
func (tree *RankTree) rangeSanitizeIndexes(start, end *int) (length int) {
	// Sanitize indexes
	if *start < 0 {
		if *start += tree.count; *start < 0 {
			*start = 0
		}
	}

	if *end < 0 {
		*end += tree.count
	}

	if *start > *end || *start >= tree.count {
		return 0
	}

	if *end >= tree.count {
		*end = tree.count - 1
	}

	return *end - *start + 1
}


// Basic Function of Range(), RevRange().
// Returns the specified range of members in the RankTree.
// If reverse is false, members are ordered from the lowest to the highest score.
// Otherwise, members are ordered from the highest to the lowest score.
// Lexicographical is used for members with equal score.
func (tree *RankTree) rangeBasic(start, end int, reverse bool) []string {
	// sanitize indexes
	rangeLen := tree.rangeSanitizeIndexes(&start, &end)
	if rangeLen == 0 {
		return make([]string, 0)
	}

	result := make([]string, rangeLen)
	var idx, skip int
	if reverse {
		skip = start
		idx = 0
	} else {
		skip = tree.count - end - 1
		idx = rangeLen - 1
	}

	// find first node
	node , index := tree.findFromRight(skip, reverse)
	gen := &rankResultGenerator{node, index, reverse}

	// collect result from linked list
	for i := 0; i < rangeLen; i++ {
		result[idx] = gen.Member()
		gen.Next()
		if reverse {
			idx++
		} else {
			idx--
		}
	}

	return result
}


// Basic Function of RangeWithScore(), RevRangeWithScore().
// Returns the specified range of members with tis score in the RankTree,
// If reverse is false, members are ordered from the lowest to the highest score.
// Otherwise, members are ordered from the highest to the lowest score.
// Lexicographical is used for members with equal score.
func (tree *RankTree) rangeWithScore(start, end int, reverse bool) []RankWithScore {
	// sanitize indexes
	rangeLen := tree.rangeSanitizeIndexes(&start, &end)
	if rangeLen == 0 {
		return make([]RankWithScore, 0)
	}

	result := make([]RankWithScore, rangeLen)
	var idx, skip int
	if reverse {
		skip = start
		idx = 0
	} else {
		skip = tree.count - end - 1
		idx = rangeLen - 1
	}

	// find first node
	node , index := tree.findFromRight(skip, reverse)
	gen := &rankResultGenerator{node, index, reverse}

	// collect result from linked list
	for i := 0; i < rangeLen; i++ {
		result[idx] = gen.RankWithScore()
		gen.Next()
		if reverse {
			idx++
		} else {
			idx--
		}
	}

	return result
}


// Returns the specified range of members in the RankTree.
// Members are ordered from the lowest to the highest score.
// Lexicographical is used for members with equal score.
func (tree *RankTree) Range(start, end int) []string {
	return tree.rangeBasic(start, end, false)
}


// Returns the specified range of members in the RankTree.
// Members are ordered from the highest to the lowest score.
// Lexicographical is used for members with equal score.
func (tree *RankTree) RevRange(start, end int) []string {
	return tree.rangeBasic(start, end, true)
}


// Returns the specified range of members with tis score in the RankTree,
// Members are ordered from the lowest to the highest score.
// Lexicographical is used for members with equal score.
func (tree *RankTree) RangeWithScore(start, end int) []RankWithScore {
	return tree.rangeWithScore(start, end, false)
}


// Returns the specified range of members with tis score in the RankTree,
// Members are ordered from the highest to the lowest score.
// Lexicographical is used for members with equal score.
func (tree *RankTree) RevRangeWithScore(start, end int) []RankWithScore {
	return tree.rangeWithScore(start, end, true)
}


// Basic Function of RangeByScore(), RevRangeByScore().
// Returns all the members in the RankTree with a score between min and max.
// If reverse is false, members are ordered from the lowest to the highest score.
// Otherwise, members are ordered from the highest to the lowest score.
// Lexicographical is used for members with equal score.
func (tree *RankTree) rangeByScoreBasic(min, max int, reverse bool) (ranks []RankWithScore) {
	if min < tree.minScore {
		min = tree.minScore
	}

	if max > tree.maxScore {
		max = tree.maxScore
	}

	if min > max {
		ranks = make([]RankWithScore, 0)
		return
	}


	node := tree.find(max)
	if node.count == 0 {
		node = tree.findNextSmallerNode(node)
	}

	if node != nil {
		length := tree.Count(min, max)
		ranks = make([]RankWithScore, length)
		var idx, index int
		if reverse == false {
			idx = length - 1
			index = node.count - 1
		}
		gen := &rankResultGenerator{node, index, reverse}


		for i := 0; i < length; i++ {
			ranks[idx] = gen.RankWithScore()
			gen.Next()
			if reverse {
				idx++
			} else {
				idx--
			}
		}
	}
	fmt.Println(ranks)
	return
}


// Returns all the members in the RankTree with a score between min and max.
// Members are ordered from the lowest to the highest score.
// Lexicographical is used for members with equal score.
func (tree *RankTree) RangeByScore(min, max int) (ranks []RankWithScore) {
	return tree.rangeByScoreBasic(min, max, false)
}


// Returns all the members in the RankTree with a score between min and max.
// Members are ordered from the highest to the lowest score.
// Lexicographical is used for members with equal score.
func (tree *RankTree) RevRangeByScore(min, max int) (ranks []RankWithScore) {
	return tree.rangeByScoreBasic(min, max, true)
}


// Adds a node element to the linked list.
func (tree *RankTree) createListElement(node *TreeNode) {
	// node must be leaf node
	if node.low != node.high {
		return
	}

	if tree.count == 0 {
		e := tree.list.PushFront(node)
		node.element = e
	} else if tree.count < 10 { // Linear Search Threshold = 10
		// Linear Search  O(N)
		var target *list.Element
		n := node.low
		e := tree.list.Head()
		for e != nil {
			v := e.Value.(*TreeNode)
			if v.low >= n {
				target = e
			} else {
				break
			}
			e = e.Next()
		}

		// node is the greatest
		if target == nil {
			e := tree.list.PushFront(node)
			node.element = e
		} else { // node is less than target
			e := tree.list.InsertAfter(node, target)
			node.element = e
		}
	} else {
		// Tree Search O(LogN)
		target := tree.findNextGreaterElement(node)
		e := tree.list.InsertAfter(node, target)
		node.element = e
	}
}


// Find a next greater node in the RankTree.
// If <node> is the greatest, nil is returned.
func (tree *RankTree) findNextGreaterNode(node *TreeNode) *TreeNode {
	count := node.count
	p := node.parent
	var target *TreeNode = nil

	for p != nil {
		if p.left == node && p.count > count {
			target = p
			break
		}
		count = p.count
		node = p
		p = p.parent
	}


	// node is the greatest
	if target == nil {
		return nil
	} else {
		// find next greater
		n := p.right
		for n.low != n.high {
			if n.left.count > 0 {
				n = n.left
			} else if n.right.count > 0 {
				n = n.right
			} else {
				log.Fatal("findNextGreaterElement left,right=0")
			}
		}
		return n
	}
}


// Find a next smaller node in the RankTree.
// If <node> is the smallest, nil is returned.
func (tree *RankTree) findNextSmallerNode(node *TreeNode) *TreeNode {
	count := node.count
	p := node.parent
	var target *TreeNode

	for p != nil {
		if p.right == node && p.count > count {
			target = p
			break
		}
		count = p.count
		node = p
		p = p.parent
	}

	// node is the smallest
	if target == nil {
		return nil
	} else {
		n := p.left
		for n.low != n.high {
			if n.right.count > 0 {
				n = n.right
			} else if n.left.count > 0 {
				n = n.left
			} else {
				log.Fatal("findNextSmallerNode left,right=0")
			}
		}
		return n
	}
}


// Find a next greater element of the list.
// If the element of <node> is the greatest node, the root element of the list is returned.
func (tree *RankTree) findNextGreaterElement(node *TreeNode) *list.Element {
	n := tree.findNextGreaterNode(node)

	if n == nil {
		return tree.list.Root()
	}
	return n.element
}


// Find a node with <score>.
// If the node does not exist or <score> is out of the range, nil is returned.
func (tree *RankTree) find(score int) *TreeNode {
	if tree.root == nil || score < tree.minScore || score > tree.maxScore {
		return nil
	}

	node := tree.root
	for  {
		low, high := node.low, node.high
		if low == high {
			if low == score {
				return node
			}
			return nil
		} else {
			mid := (low + high) / 2
			if score <= mid {
				node = node.left
			} else {
				node = node.right
			}
		}
	}
}


 // Basic Function for rangeBasic(), rangeWithScore().
 // Find a leaf node right to left, skip <count> node(s).
 // If not found, (nil, 0) is returned.
 // <index> represents the index of member in node.members.
func (tree *RankTree) findFromRight(count int, reverse bool) (node *TreeNode, index int) {
	if count < 0 || count >= tree.count {
		return nil, 0
	}

	node = tree.root
	skip := count

	for node.low < node.high {
		if node.right.count > 0 && skip < node.right.count {
			node = node.right
		} else {
			if skip > 0 { // && node.right.count > 0
				skip -= node.right.count
			}
			node = node.left
		}
	}

	if reverse { // index start from left to right
		index = skip
	} else {  // index start from right to left
		index = node.count - skip - 1
	}

	return node, index
}


// Increases the count of the node and the parents by <delta>.
func (node *TreeNode) incrementCount(delta int) {
	for {
		node.count += delta
		node = node.parent
		if node == nil {
			break
		}
	}
}


// Returns count of the left area.
func (node *TreeNode) countLeftArea() (sum int) {
	for node.parent != nil {
		thisNode := node
		node = node.parent
		if node.left != thisNode {
			sum += node.left.count
		}
	}
	return
}


// Returns count of the right area.
func (node *TreeNode) countRightArea() (sum int) {
	for node.parent != nil {
		thisNode := node
		node = node.parent
		if node.right != thisNode {
			sum += node.right.count
		}
	}
	return
}


// Creates children nodes for <parent>.
func (node *TreeNode) create(low int, high int, parent *TreeNode) {
	node.low = low
	node.high = high
	node.parent = parent

	if low < high {
		mid := (low + high) / 2
		node.left = new(TreeNode)
		node.right = new(TreeNode)
		node.left.create(low, mid, node)
		node.right.create(mid + 1, high, node)
	}
}


// Next result.
func (r *rankResultGenerator) Next() {
	if r.reverse {
		if r.index < r.node.count - 1 {
			r.index++
		} else {
			e := r.node.element
			if e = e.Next(); e != nil {
				r.node = e.Value.(*TreeNode)
				r.index = 0
			}
		}
	} else {
		if r.index > 0 {
			r.index--
		} else {
			e := r.node.element
			if e = e.Next(); e != nil {
				r.node = e.Value.(*TreeNode)
				r.index = r.node.count - 1
			}
		}
	}
}


// Returns member of the result node.
func (r *rankResultGenerator) Member() string {
	return r.node.members[r.index]
}

// Returns rank and score of the result node.
func (r *rankResultGenerator) RankWithScore() RankWithScore {
	return RankWithScore{
		Member: r.node.members[r.index],
		Score: r.node.low }
}



