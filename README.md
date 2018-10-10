#  RankTree

RankTree is an in-memory data structure, written in [Go](http://golang.org/). used as a ranking database. It is similar to [Redis](https://redis.io/) [Sorted Sets](https://redis.io/commands#sorted_set), but based on Complete Binary Tree. 



## Features

* In-memory database
* High Performance
* Redis-like Commands



Restrictions:  the value of score must be non-negative integers. 



## Installation

Install:

```
go get -u https://github.com/ng1091/ranktree
```

Import:

```go
import "github.com/ng1091/ranktree"
```



## Usage

```go
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
```

Output:

```
 Output:
 Bob is No.1

 No.1: Bob
 No.2: Alice
 No.3: Charles
```



## Documents

### Commands

```
    New(low int, high int) (*RankTree, error)
    Add(member string, score int) bool
    Card() int
    Count(min, max int) int
    IncrementBy(member string, score int) int
    PopMax() (rank *RankWithScore)
    PopMaxN(n int) (ranks []RankWithScore)
    PopMin() (rank *RankWithScore)
    PopMinN(n int) (ranks []RankWithScore)
    Range(start, end int) []string
    RangeByScore(min, max int) (ranks []RankWithScore)
    RangeWithScore(start, end int) []RankWithScore
    Rank(member string) int
    Remove(members ...string) (sum int)
    RevRange(start, end int) []string
    RevRangeByScore(min, max int) (ranks []RankWithScore)
    RevRangeWithScore(start, end int) []RankWithScore
    RevRank(member string) int
    Score(member string) int
    UpdateScore(member string, score int, insert bool) bool
```



**Please check  [GoDoc - ranktree](https://www.godoc.org/github.com/ng1091/ranktree) for more details.**

