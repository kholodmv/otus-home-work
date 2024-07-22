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

// ListItem it is a single element in a doubly linked list.
type ListItem struct {
	// Value - list item value
	Value interface{}
	// Next - next element
	Next *ListItem
	// Prev - previous element
	Prev *ListItem
}

// list doubly linked list with efficient insertion and deletion.
type list struct {
	// front - first element of the list
	front *ListItem
	// back - last element of the list
	back *ListItem
	// len - list length
	len int
}

// NewList creates a new empty list.
func NewList() List {
	return &list{}
}

// Len returns the number of elements in the list.
func (l *list) Len() int {
	return l.len
}

// Front returns the first element of the list.
func (l *list) Front() *ListItem {
	return l.front
}

// Back returns the last element of the list.
func (l *list) Back() *ListItem {
	return l.back
}

// PushFront adds a value to the front of the list.
func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.front == nil {
		l.front = item
		l.back = item
	} else {
		// Adding an element to the beginning of the list
		item.Next = l.front
		l.front.Prev = item
		l.front = item
	}

	l.len++
	return item
}

// PushBack adds a value to the back of the list.
func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}

	if l.back == nil {
		l.front = item
		l.back = item
	} else {
		// Adding an element to the end of the list
		item.Prev = l.back
		l.back.Next = item
		l.back = item
	}
	l.len++
	return item
}

// Remove removes an item from the list.
func (l *list) Remove(i *ListItem) {
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

// MoveToFront moves an item to the front of the list.
func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}
