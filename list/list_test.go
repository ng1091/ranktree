package list

import (
	"testing"
	"fmt"
)


func printNodes (l *List) {
	e := l.Head()
	for e != nil {
		fmt.Printf("%s ", e.Value)
		e = e.Next()
	}
	fmt.Println("")
}


func TestList(t *testing.T) {

	l := New()

	count := 0
	if l.Len() != count {
		t.Error("Len没通过")
	}

	if l.Head() != nil {
		t.Error("Head is not nil")
	}


	ea := l.PushFront("a")
	count++

	if l.Len() != count {
		t.Error("Len failed")
	}

	if ea.next != nil {
		t.Error("a.next is not nil")
	}

	if l.Head() != ea {
		t.Error("head error")
	}


	eb := l.PushFront("b")
	ec := l.PushFront("c")

	printNodes(l)


	l.InsertAfter("d", ec)

	printNodes(l)

	b1 := l.InsertAfter("b1", eb)

	printNodes(l)


	l.Remove(b1)

	printNodes(l)

	l.RemoveNext(ea)

	printNodes(l)

	l.RemoveNext(eb)
	printNodes(l)

	if l.Len() != 3 {
		t.Error("len != 3")
	}

	l.RemoveNext(l.Head())
	printNodes(l)

}

