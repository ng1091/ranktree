package ranktree

import (
	"testing"
	"container/list"
)



func checkRankTree(t *testing.T, tree *RankTree, low, high, count int) {
	if n := tree.count; n != count {
		t.Errorf("tree.count=%d, want %d", n, count)
	}

	if n := len(tree.nodeMap); n != count {
		t.Errorf("len(tree.nodeMap)=%d, want %d", n, count)
	}


	// check leaf node

	checkValue := low
	stack := list.New()
	n := tree.root
	calcCount := 0

	for stack.Len() > 0 || n != nil {
		for n != nil {
			stack.PushBack(n)
			n = n.left
		}

		if stack.Len() > 0 {
			e := stack.Back()
			stack.Remove(e)
			var ok bool
			if n, ok = e.Value.(*TreeNode); ok {
				if n.low == n.high {
					calcCount += n.count
					if n.low != checkValue {
						t.Errorf("%p, node.low = %d, want %d", n, n.low, checkValue)
					}
					//t.Logf("%p, node.low = %d", n, checkValue)
					checkValue++
				}
				n = n.right
			}  else {
				t.Errorf("node %p is not *TreeNode, %#v", n, n)
			}
		}
	}


	// check count
	if n := tree.count; n != count {
		t.Errorf("tree.count = %d, want %d", n, count)
	}

	if calcCount != count {
		t.Errorf("calcCount = %d, want %d", calcCount, count)
	}
}


func TestNew(t *testing.T) {
	// one node tree
	tree, err := New(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	checkRankTree(t, tree, 0, 0, 0)


	// two node tree
	tree, err = New(1, 2)
	if err != nil {
		t.Error(err)
	}
	checkRankTree(t, tree, 1, 2, 0)


	// many node tree
	tree, err = New(0, 256)
	if err != nil {
		t.Error(err)
	}
	checkRankTree(t, tree, 0, 256, 0)
}


func TestRankTree_Add(t *testing.T) {
	tree, err := New(1, 8)
	if err != nil {
		t.Fatal(err)
	}

	if tree.Add("a", 1) != true {
		t.Error("tree.Add() return false, want true")
	}

	if tree.Add("a", 2) != false {
		t.Error("tree.Add() return true, want false")
	}

	if tree.Add("b", 8) != true {
		t.Error("tree.Add() return false, want true")
	}

	if tree.Add("c", 9) != false {
		t.Error("tree.Add() return true, want false")
	}

	checkRankTree(t, tree, 1, 8, 2)

}


func checkListNode(t *testing.T, tree *RankTree, list []*TreeNode) {
	l := tree.list.Root()

	if n := tree.list.Len(); n != len(list) {
		t.Errorf("tree.list.Len() = %d, want %d", n, len(list))
	}

	for i := len(list) - 1; i > 0; i-- {
		l = l.Next()
		if l == nil {
			t.Error("l.Next() = nil")
		}

		if n, ok := l.Value.(*TreeNode); ok {
			if list[i] != n {
				t.Errorf("tree.list[%d] = %d (%p), want %d (%p)", len(list) - i - 1, n.low, n, list[i].low, list[i])
			}
		} else {
			t.Errorf("%p, l.Value.(*TreeNode) failed", l.Value)
		}
	}
}


func TestList(t *testing.T) {
	tree, err := New(1, 8)
	if err != nil {
		t.Fatal(err)
	}

	tree.Add("a", 1)
	a := tree.find(1)
	checkListNode(t, tree, []*TreeNode{a})

	tree.Add("b", 2)
	b := tree.find(2)
	checkListNode(t, tree, []*TreeNode{a, b})

	tree.Add("c", 3)
	c := tree.find(3)
	checkListNode(t, tree, []*TreeNode{a, b, c})

	tree.Add("d", 5)
	d := tree.find(5)
	checkListNode(t, tree, []*TreeNode{a, b, c, d})

	tree.Add("e", 8)
	e := tree.find(8)
	checkListNode(t, tree, []*TreeNode{a, b, c, d, e})
}


func TestList2(t *testing.T) {
	tree, err := New(1, 8)
	if err != nil {
		t.Fatal(err)
	}

	tree.Add("a", 7)
	a := tree.find(7)
	checkListNode(t, tree, []*TreeNode{a})

	tree.Add("b", 8)
	b := tree.find(8)
	checkListNode(t, tree, []*TreeNode{a, b})


	tree.Add("c", 5)
	c := tree.find(5)
	checkListNode(t, tree, []*TreeNode{c, a, b})


	tree.Add("d", 1)
	d := tree.find(1)
	checkListNode(t, tree, []*TreeNode{d, c, a, b})

	tree.Add("f", 6)
	e := tree.find(6)
	checkListNode(t, tree, []*TreeNode{d, c, e, a, b})

	tree.Add("g", 5)
	g := tree.find(5)
	checkListNode(t, tree, []*TreeNode{d, c, e, a, b})
	checkListNode(t, tree, []*TreeNode{d, g, e, a, b})
}



func checkRank(t *testing.T, rank []string, s []string) {
	if len(rank) != len(s) {
		t.Errorf("len(rank) = %d, want %d", len(rank), len(s))
	}

	for i, v := range rank {
		if v != s[i] {
			t.Errorf("rank[%d] = %s, want %s", i, v, s[i])
		}
	}
}


func checkRankWithScore(t *testing.T, rank []RankWithScore, member []string, score []int) {
	r := make([]string, len(rank))
	for i, v := range rank {
		r[i] = v.Member
		if v.Score != score[i] {
			t.Errorf("rank[%d].Score = %d, want %d", i, v.Score, score[i])
		}
	}
	checkRank(t, r, member)
}


func TestRankTree_Range(t *testing.T) {
	tree, err := New(1, 8)
	if err != nil {
		t.Fatal(err)
	}

	tree.Add("a", 1)
	tree.Add("b", 2)
	tree.Add("c", 3)
	tree.Add("d", 4)
	tree.Add("e", 5)

	checkRank(t, tree.Range(0, 0), []string{"a"})
	checkRank(t, tree.Range(0, 1), []string{"a", "b"})
	checkRank(t, tree.Range(1, 3), []string{"b", "c", "d"})
	checkRank(t, tree.Range(4, 4), []string{"e"})

	checkRank(t, tree.Range(0, -1), []string{"a", "b", "c", "d", "e"})
	checkRank(t, tree.Range(-3, -1), []string{"c", "d", "e"})
	checkRank(t, tree.Range(2, -2), []string{"c", "d"})

	checkRank(t, tree.Range(-10, 8), []string{"a", "b", "c", "d", "e"})
	checkRank(t, tree.Range(-2, 0), []string{})
	checkRank(t, tree.Range(6, 8), []string{})

	tree.Add("b2", 2) // a b b2 c d e
	checkRank(t, tree.Range(1, 3), []string{"b", "b2", "c"})

	tree.Add("b1", 2) // a b b1 b2 c d e
	checkRank(t, tree.Range(1, 3), []string{"b", "b1", "b2"})

	tree.Add("c2", 3) // a b b1 b2 c c2 d e
	checkRank(t, tree.Range(1, 5), []string{"b", "b1", "b2", "c", "c2"})

	tree.Add("1", 1) // 1 a b b1 b2 c c2 d e
	checkRank(t, tree.Range(0, 1), []string{"1", "a"})

	tree.Add("f", 5) // 1 a b b1 b2 c c2 d e f
	checkRank(t, tree.Range(-2, -1), []string{"e", "f"})
	checkRank(t, tree.Range(0, -1), []string{"1", "a", "b", "b1", "b2", "c", "c2", "d", "e", "f"})
}


func TestRankTree_RevRange(t *testing.T) {
	tree, err := New(1, 8)
	if err != nil {
		t.Fatal(err)
	}

	tree.Add("a", 5)
	tree.Add("b", 4)
	tree.Add("c", 3)
	tree.Add("d", 2)
	tree.Add("e", 1)

	checkRank(t, tree.RevRange(0, 0), []string{"a"})
	checkRank(t, tree.RevRange(0, 1), []string{"a", "b"})
	checkRank(t, tree.RevRange(1, 3), []string{"b", "c", "d"})
	checkRank(t, tree.RevRange(4, 4), []string{"e"})

	checkRank(t, tree.RevRange(0, -1), []string{"a", "b", "c", "d", "e"})
	checkRank(t, tree.RevRange(-3, -1), []string{"c", "d", "e"})
	checkRank(t, tree.RevRange(2, -2), []string{"c", "d"})

	checkRank(t, tree.RevRange(-10, 8), []string{"a", "b", "c", "d", "e"})
	checkRank(t, tree.RevRange(-2, 0), []string{})
	checkRank(t, tree.RevRange(6, 8), []string{})


	tree.Add("b2", 4) // a b b2 c d e
	checkRank(t, tree.RevRange(1, 3), []string{"b", "b2", "c"})

	tree.Add("b1", 4) // a b b1 b2 c d e
	checkRank(t, tree.RevRange(1, 3), []string{"b", "b1", "b2"})

	tree.Add("c2", 3) // a b b1 b2 c c2 d e
	checkRank(t, tree.RevRange(1, 5), []string{"b", "b1", "b2", "c", "c2"})

	tree.Add("1", 5) // 1 a b b1 b2 c c2 d e
	checkRank(t, tree.RevRange(0, 1), []string{"1", "a"})

	tree.Add("f", 1) // 1 a b b1 b2 c c2 d e f
	checkRank(t, tree.RevRange(-2, -1), []string{"e", "f"})
	checkRank(t, tree.RevRange(0, -1), []string{"1", "a", "b", "b1", "b2", "c", "c2", "d", "e", "f"})
}


func TestRankTree_RangeWithScore(t *testing.T) {
	tree, err := New(1, 8)
	if err != nil {
		t.Fatal(err)
	}

	tree.Add("a", 1)
	tree.Add("b", 2)
	tree.Add("c", 3)
	tree.Add("d", 4)
	tree.Add("e", 5)

	checkRankWithScore(t, tree.RangeWithScore(0, 0), []string{"a"}, []int{1})
	checkRankWithScore(t, tree.RangeWithScore(0, 1), []string{"a", "b"}, []int{1, 2})
	checkRankWithScore(t, tree.RangeWithScore(1, 3), []string{"b", "c", "d"}, []int{2, 3, 4})
	checkRankWithScore(t, tree.RangeWithScore(4, 4), []string{"e"}, []int{5})
	
	checkRankWithScore(t, tree.RangeWithScore(0, -1), []string{"a", "b", "c", "d", "e"}, []int{1, 2, 3, 4, 5})
	checkRankWithScore(t, tree.RangeWithScore(-3, -1), []string{"c", "d", "e"}, []int{3, 4, 5})
	checkRankWithScore(t, tree.RangeWithScore(2, -2), []string{"c", "d"}, []int{3, 4})

	checkRankWithScore(t, tree.RangeWithScore(-10, 8), []string{"a", "b", "c", "d", "e"}, []int{1, 2, 3, 4, 5})
	checkRankWithScore(t, tree.RangeWithScore(-2, 0), []string{}, []int{})
	checkRankWithScore(t, tree.RangeWithScore(6, 8), []string{}, []int{})


	tree.Add("b2", 2) // a b b2 c d e
	checkRankWithScore(t, tree.RangeWithScore(1, 3), []string{"b", "b2", "c"}, []int{2, 2, 3})

	tree.Add("b1", 2) // a b b1 b2 c d e
	checkRankWithScore(t, tree.RangeWithScore(1, 3), []string{"b", "b1", "b2"}, []int{2, 2, 2})

	tree.Add("c2", 3) // a b b1 b2 c c2 d e
	checkRankWithScore(t, tree.RangeWithScore(1, 5), []string{"b", "b1", "b2", "c", "c2"}, []int{2, 2, 2, 3, 3})

	tree.Add("1", 1) // 1 a b b1 b2 c c2 d e
	checkRankWithScore(t, tree.RangeWithScore(0, 1), []string{"1", "a"}, []int{1, 1})

	tree.Add("f", 5) // 1 a b b1 b2 c c2 d e f
	checkRankWithScore(t, tree.RangeWithScore(-2, -1), []string{"e", "f"}, []int{5, 5})
	checkRankWithScore(t, tree.RangeWithScore(0, -1), []string{"1", "a", "b", "b1", "b2", "c", "c2", "d", "e", "f"}, []int{1, 1, 2, 2, 2, 3, 3, 4, 5, 5})
}


func TestRankTree_RevRangeWithScore(t *testing.T) {
	tree, err := New(1, 8)
	if err != nil {
		t.Fatal(err)
	}

	tree.Add("a", 5)
	tree.Add("b", 4)
	tree.Add("c", 3)
	tree.Add("d", 2)
	tree.Add("e", 1)

	checkRankWithScore(t, tree.RevRangeWithScore(0, 0), []string{"a"}, []int{5})
	checkRankWithScore(t, tree.RevRangeWithScore(0, 1), []string{"a", "b"}, []int{5, 4})
	checkRankWithScore(t, tree.RevRangeWithScore(1, 3), []string{"b", "c", "d"}, []int{4, 3, 2})
	checkRankWithScore(t, tree.RevRangeWithScore(4, 4), []string{"e"}, []int{1})
	
	checkRankWithScore(t, tree.RevRangeWithScore(0, -1), []string{"a", "b", "c", "d", "e"}, []int{5, 4, 3, 2, 1})
	checkRankWithScore(t, tree.RevRangeWithScore(-3, -1), []string{"c", "d", "e"}, []int{3, 2, 1})
	checkRankWithScore(t, tree.RevRangeWithScore(2, -2), []string{"c", "d"}, []int{3, 2})
	
	checkRankWithScore(t, tree.RevRangeWithScore(-10, 8), []string{"a", "b", "c", "d", "e"}, []int{5, 4, 3, 2, 1})
	checkRankWithScore(t, tree.RevRangeWithScore(-2, 0), []string{}, []int{})
	checkRankWithScore(t, tree.RevRangeWithScore(6, 8), []string{}, []int{})


	tree.Add("b2", 4) // a b b2 c d e
	checkRankWithScore(t, tree.RevRangeWithScore(1, 3), []string{"b", "b2", "c"}, []int{4, 4, 3})

	tree.Add("b1", 4) // a b b1 b2 c d e
	checkRankWithScore(t, tree.RevRangeWithScore(1, 3), []string{"b", "b1", "b2"}, []int{4, 4, 4})

	tree.Add("c2", 3) // a b b1 b2 c c2 d e
	checkRankWithScore(t, tree.RevRangeWithScore(1, 5), []string{"b", "b1", "b2", "c", "c2"}, []int{4, 4, 4, 3, 3})

	tree.Add("1", 5) // 1 a b b1 b2 c c2 d e
	checkRankWithScore(t, tree.RevRangeWithScore(0, 1), []string{"1", "a"}, []int{5, 5})

	tree.Add("f", 1) // 1 a b b1 b2 c c2 d e f
	checkRankWithScore(t, tree.RevRangeWithScore(-2, -1), []string{"e", "f"}, []int{1, 1})
	checkRankWithScore(t, tree.RevRangeWithScore(0, -1), []string{"1", "a", "b", "b1", "b2", "c", "c2", "d", "e", "f"}, []int{5, 5, 4, 4, 4, 3, 3, 2, 1, 1})
}


func TestRankTree_Rank(t *testing.T) {
	tree, err := New(1, 8)
	if err != nil {
		t.Fatal(err)
	}

	tree.Add("a", 1)
	tree.Add("b", 1)
	tree.Add("c", 3)
	tree.Add("d", 5)
	tree.Add("e", 5)

	if n := tree.Rank("a"); n != 0 {
		t.Errorf("tree.Rank(\"%s\") = %d, want 0", "a", n)
	}

	if n := tree.Rank("b"); n != 1 {
		t.Errorf("tree.Rank(\"%s\") = %d, want 1", "b", n)
	}

	if n := tree.Rank("c"); n != 2 {
		t.Errorf("tree.Rank(\"%s\") = %d, want 2", "c", n)
	}

	if n := tree.Rank("d"); n != 3 {
		t.Errorf("tree.Rank(\"%s\") = %d, want 3", "d", n)
	}

	if n := tree.Rank("e"); n != 4 {
		t.Errorf("tree.Rank(\"%s\") = %d, want 4", "e", n)
	}

	if n := tree.Rank("f"); n != -1 {
		t.Errorf("tree.Rank(\"%s\") = %d, want -1", "f", n)
	}
}


func TestRankTree_RevRank(t *testing.T) {
	tree, err := New(1, 8)
	if err != nil {
		t.Fatal(err)
	}

	tree.Add("a", 1)
	tree.Add("b", 1)
	tree.Add("c", 3)
	tree.Add("d", 5)
	tree.Add("e", 5)


	if n := tree.RevRank("e"); n != 0 {
		t.Errorf("tree.Rank(\"%s\") = %d, want 0", "e", n)
	}

	if n := tree.RevRank("d"); n != 1 {
		t.Errorf("tree.Rank(\"%s\") = %d, want 1", "d", n)
	}

	if n := tree.RevRank("c"); n != 2 {
		t.Errorf("tree.Rank(\"%s\") = %d, want 2", "c", n)
	}

	if n := tree.RevRank("b"); n != 3 {
		t.Errorf("tree.Rank(\"%s\") = %d, want 3", "b", n)
	}

	if n := tree.RevRank("a"); n != 4 {
		t.Errorf("tree.Rank(\"%s\") = %d, want 4", "a", n)
	}

	if n := tree.RevRank("f"); n != -1 {
		t.Errorf("tree.Rank(\"%s\") = %d, want -1", "f", n)
	}
}


func TestRankTree_Count(t *testing.T) {
	tree, err := New(1, 8)
	if err != nil {
		t.Fatal(err)
	}

	tree.Add("a", 1)
	tree.Add("b", 2)
	tree.Add("c", 3)
	tree.Add("d", 4)
	tree.Add("e", 6)
	tree.Add("f", 8)

	if n := tree.Count(0, 9); n != 6 {
		t.Errorf("tree.Count = %d, want %d", n, 6)
	}

	if n := tree.Count(2, 5); n != 3 {
		t.Errorf("tree.Count = %d, want %d", n, 3)
	}

	if n := tree.Count(5, 8); n != 2 {
		t.Errorf("tree.Count = %d, want %d", n, 2)
	}

	if n := tree.Count(5, 7); n != 1 {
		t.Errorf("tree.Count = %d, want %d", n, 1)
	}

	if n := tree.Count(6, 6); n != 1 {
		t.Errorf("tree.Count = %d, want %d", n, 1)
	}

	if n := tree.Count(2, 4); n != 3 {
		t.Errorf("tree.Count = %d, want %d", n, 3)
	}

	if n := tree.Count(2, 1); n != 0 {
		t.Errorf("tree.Count = %d, want %d", n, 0)
	}

	if n := tree.Count(-3, -2); n != 0 {
		t.Errorf("tree.Count = %d, want %d", n, 0)
	}

	if n := tree.Count(10, 12); n != 0 {
		t.Errorf("tree.Count = %d, want %d", n, 0)
	}
}


func TestRankTree_Card(t *testing.T) {
	tree, err := New(1, 8)
	if err != nil {
		t.Fatal(err)
	}

	if n := tree.Card(); n != 0 {
		t.Errorf("tree.Card() = %d, want %d", n, 0)
	}

	tree.Add("a", 1)

	if n := tree.Card(); n != 1 {
		t.Errorf("tree.Card() = %d, want %d", n, 1)
	}

	tree.Add("b", 1)

	if n := tree.Card(); n != 2 {
		t.Errorf("tree.Card() = %d, want %d", n, 2)
	}

	tree.Add("c", 2)

	if n := tree.Card(); n != 3 {
		t.Errorf("tree.Card() = %d, want %d", n, 3)
	}

	tree.Remove("a")
	if n := tree.Card(); n != 2 {
		t.Errorf("tree.Card() = %d, want %d", n, 2)
	}

	tree.PopMax()
	if n := tree.Card(); n != 1 {
		t.Errorf("tree.Card() = %d, want %d", n, 1)
	}

	tree.PopMin()
	if n := tree.Card(); n != 0 {
		t.Errorf("tree.Card() = %d, want %d", n, 0)
	}
}


func TestRankTree_Score(t *testing.T) {
	tree, err := New(1, 256)
	if err != nil {
		t.Fatal(err)
	}

	tree.Add("a", 256)
	tree.Add("b", 256)
	tree.Add("c", 100)
	tree.Add("d", 1)
	tree.Add("e", 1)
	tree.Add("f", 1)
	tree.Remove("d")

	if n := tree.Score("a"); n != 256 {
		t.Errorf("tree.Score() = %d, want %d", n, 256)
	}

	if n := tree.Score("b"); n != 256 {
		t.Errorf("tree.Score() = %d, want %d", n, 256)
	}

	if n := tree.Score("c"); n != 100 {
		t.Errorf("tree.Score() = %d, want %d", n, 100)
	}

	if n := tree.Score("d"); n != -1 {
		t.Errorf("tree.Score() = %d, want %d", n, -1)
	}

	if n := tree.Score("e"); n != 1 {
		t.Errorf("tree.Score() = %d, want %d", n, -1)
	}

	if n := tree.Score("f"); n != 1 {
		t.Errorf("tree.Score() = %d, want %d", n, -1)
	}
}


func TestRankTree_RangeByScore(t *testing.T) {
	tree, err := New(1, 8)
	if err != nil {
		t.Fatal(err)
	}

	tree.Add("a", 1)
	tree.Add("b", 2)
	tree.Add("b2", 2)
	tree.Add("c", 4)
	tree.Add("d", 5)
	tree.Add("e", 6)
	tree.Add("e2", 6)
	tree.Add("f", 8)

	checkRankWithScore(t, tree.RangeByScore(0, 0), []string{}, []int{})
	checkRankWithScore(t, tree.RangeByScore(1, 1), []string{"a"}, []int{1})
	checkRankWithScore(t, tree.RangeByScore(-2, -2), []string{}, []int{})
	checkRankWithScore(t, tree.RangeByScore(8, 8), []string{"f"}, []int{8})
	checkRankWithScore(t, tree.RangeByScore(6, 6), []string{"e", "e2"}, []int{6, 6})
	checkRankWithScore(t, tree.RangeByScore(2, 3), []string{"b", "b2"}, []int{2, 2})
	checkRankWithScore(t, tree.RangeByScore(3, 7), []string{"c", "d", "e", "e2"}, []int{4, 5, 6, 6})
	checkRankWithScore(t, tree.RangeByScore(1, 8), []string{"a", "b", "b2", "c", "d", "e", "e2", "f"}, []int{1, 2, 2, 4, 5, 6, 6, 8})
}


func TestRankTree_RevRangeByScore(t *testing.T) {
	tree, err := New(1, 8)
	if err != nil {
		t.Fatal(err)
	}

	tree.Add("a", 1)
	tree.Add("b", 2)
	tree.Add("b2", 2)
	tree.Add("c", 4)
	tree.Add("d", 5)
	tree.Add("e", 6)
	tree.Add("e2", 6)
	tree.Add("f", 8)


	checkRankWithScore(t, tree.RevRangeByScore(0, 0), []string{}, []int{})
	checkRankWithScore(t, tree.RevRangeByScore(1, 1), []string{"a"}, []int{1})
	checkRankWithScore(t, tree.RevRangeByScore(- 2, -2), []string{}, []int{})
	checkRankWithScore(t, tree.RevRangeByScore(8, 8), []string{"f"}, []int{8})
	checkRankWithScore(t, tree.RevRangeByScore(6, 6), []string{"e", "e2"}, []int{6, 6})
	checkRankWithScore(t, tree.RevRangeByScore(2, 3), []string{"b", "b2"}, []int{2, 2})
	checkRankWithScore(t, tree.RevRangeByScore(3, 7), []string{"e", "e2", "d", "c"}, []int{6, 6, 5, 4})
	checkRankWithScore(t, tree.RevRangeByScore(1, 8), []string{"f", "e", "e2", "d", "c", "b", "b2", "a"}, []int{8, 6, 6, 5, 4, 2, 2, 1})
}

func TestRankTree_PopMax(t *testing.T) {

}



func TestRankTree_Remove(t *testing.T) {
	tree, err := New(1, 8)
	if err != nil {
		t.Fatal(err)
	}

	tree.Add("a", 1)

	if n := tree.Remove("c"); n != 0 {
		t.Errorf("tree.Remove() = %d, want %d", n, 0)
	}
}





