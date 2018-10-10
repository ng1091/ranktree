// Package list implements a singly linked list.
package list

type Element struct {
	next *Element
	Value interface{}
}


func (e *Element) Next() *Element {
	if p := e.next; p != nil {
		return p
	}
	return nil
}


type List struct {
	root Element
	last *Element
	len int
}


func (l *List) Init() *List {
	l.root.next = nil
	l.last = nil
	l.len = 0
	return l
}


func New() *List {
	return new(List).Init()
}


func (l *List) Len() int  {
	return l.len
}


func (l *List) Root() *Element {
	return &l.root
}


func (l *List) Head() *Element {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}


func (l *List) Back() *Element {
	if l.len == 0 {
		return nil
	}
	return l.last
}


func (l *List) insert(e, at *Element) *Element {
	n := at.next
	e.next = n
	at.next = e
	l.len++

	if n == nil {
		l.last = e
	}
	return e
}


func (l *List) insertValue(v interface{}, at *Element) *Element {
	return l.insert(&Element{Value: v}, at)
}


func (l *List) remove(e *Element) *Element {
	h := &l.root
	for h.next != nil {
		if h.next == e {
			h.next = e.next
			e.next = nil
			l.len--
			if l.len > 0 && l.Back() == e {
				l.last = h
			}
			break
		}
		h = h.next
	}
	return e
}


func (l *List) removeNext(e *Element) *Element {
	n := e.next
	if n != nil {
		e.next = n.next
		n.next = nil
		l.len--
		if l.len > 0 && n == l.last {
			l.last = e
		}
	}
	return n
}


func (l *List) Remove(e *Element) *Element {
	return l.remove(e)
}


func (l *List) RemoveNext(e *Element) *Element {
	return l.removeNext(e)
}


func (l *List) PushFront(v interface{}) *Element {
	return l.insertValue(v, &l.root)
}


func (l *List) InsertAfter(v interface{}, mark *Element) *Element {
	return l.insertValue(v, mark)
}










