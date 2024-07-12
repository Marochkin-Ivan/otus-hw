package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v any) *ListItem
	PushBack(v any) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value any
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head *ListItem
	tail *ListItem
	len  int
}

func NewList() List {
	return new(list)
}

func NewListItem(v any, prev, next *ListItem) *ListItem {
	return &ListItem{
		v,
		next,
		prev,
	}
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v any) *ListItem {
	var item *ListItem
	if l.len == 0 { // пустой список
		item = NewListItem(v, nil, nil)
		l.head = item
		l.tail = item
	} else {
		item = NewListItem(v, nil, l.head)
		l.head.Prev = item
		l.head = item
	}

	l.len++
	return item
}

func (l *list) PushBack(v any) *ListItem {
	var item *ListItem
	if l.len == 0 {
		item = NewListItem(v, nil, nil)
		l.tail = item
		l.head = item
	} else {
		item = NewListItem(v, l.tail, nil)
		l.tail.Next = item
		l.tail = item
	}

	l.len++
	return item
}

func (l *list) Remove(i *ListItem) {
	if i.Prev == nil {
		l.head = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		l.tail = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev == nil { // элемент уже впереди
		return
	}

	l.Remove(i)
	l.PushFront(i.Value)
}
