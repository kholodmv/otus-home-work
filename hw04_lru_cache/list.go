package hw04lrucache

type IList interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

// ListItem it is a single element in a doubly linked List.
type ListItem struct {
	Value interface{} // value
	Next  *ListItem   // next element
	Prev  *ListItem   // previous element
}

// List doubly linked List with efficient insertion and deletion.
type List struct {
	front *ListItem // first element of the List
	back  *ListItem // last element of the List
	len   int       // List length
}

// NewList creates a new empty List.
func NewList() *List {
	return &List{}
}

// Len returns the number of elements in the List.
func (l *List) Len() int {
	return l.len
}

// Front returns the first element of the List.
func (l *List) Front() *ListItem {
	return l.front
}

// Back returns the last element of the List.
func (l *List) Back() *ListItem {
	return l.back
}

// PushFront adds a value to the front of the List.
func (l *List) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.front == nil {
		l.front = item
		l.back = item
	} else {
		// Adding an element to the beginning of the List
		item.Next = l.front
		l.front.Prev = item
		l.front = item
	}

	l.len++
	return item
}

// PushBack adds a value to the back of the List.
func (l *List) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}

	if l.back == nil {
		l.front = item
		l.back = item
	} else {
		// Adding an element to the end of the List
		item.Prev = l.back
		l.back.Next = item
		l.back = item
	}
	l.len++
	return item
}

// Remove removes an item from the List.
func (l *List) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.front = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}
	l.len--
}

// MoveToFront moves an item to the front of the List.
func (l *List) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}
