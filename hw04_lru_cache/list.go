package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	front *ListItem
	back  *ListItem
	len   int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	defer l.increaseLen()

	if l.len == 0 {
		l.front = &ListItem{
			Value: v,
			Next:  nil,
			Prev:  nil,
		}
		l.back = l.front
		return l.front
	}

	l.front.Prev = &ListItem{
		Value: v,
		Next:  l.front,
		Prev:  nil,
	}
	l.front = l.front.Prev
	return l.front
}

func (l *list) PushBack(v interface{}) *ListItem {
	defer l.increaseLen()

	if l.len == 0 {
		l.front = &ListItem{
			Value: v,
			Next:  nil,
			Prev:  nil,
		}
		l.back = l.front
		return l.back
	}

	l.back.Next = &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.back,
	}
	l.back = l.back.Next
	return l.back
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	l.len--

	if l.len == 0 {
		l.front = nil
		l.back = nil
		return
	}

	if i == l.back {
		l.back = i.Prev
		l.back.Next = nil
		return
	}

	if i == l.front {
		l.front = i.Next
		l.front.Prev = nil
		return
	}

	iPrev := i.Prev
	iPrev.Next = i.Next
	i.Next.Prev = iPrev
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil {
		return
	}

	if i == l.front {
		return
	}

	iPrev := i.Prev
	iNext := i.Next

	iPrev.Next = iNext
	if i != l.back {
		iNext.Prev = iPrev
	} else {
		l.back = iPrev
	}

	l.front.Prev = i
	i.Next = l.front
	i.Prev = nil
	l.front = i
}

func (l *list) increaseLen() {
	l.len++
}

func NewList() List {
	return &list{}
}
