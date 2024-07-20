package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty List", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}

func TestRemove(t *testing.T) {
	l := &List{}
	item1 := &ListItem{Value: "first element"}
	item2 := &ListItem{Value: "second element"}
	item3 := &ListItem{Value: "third element"}

	l.front = item1
	l.back = item3
	l.len = 3

	item1.Next = item2
	item2.Prev = item1
	item2.Next = item3
	item3.Prev = item2

	l.Remove(item2)
	if l.len != 2 {
		t.Errorf("expected length 2, got %d", l.len)
	}

	if item1.Next != item3 {
		t.Errorf("first element's Next should be third element")
	}

	if item3.Prev != item1 {
		t.Errorf("third element's Prev should be first element")
	}

	l.Remove(item1)
	if l.len != 1 {
		t.Errorf("expected length 1, got %d", l.len)
	}

	if l.front != item3 {
		t.Errorf("front should be third element")
	}

	if item3.Prev != nil {
		t.Errorf("third element's Prev should be nil")
	}

	l.Remove(item3)
	if l.len != 0 {
		t.Errorf("expected length 0, got %d", l.len)
	}

	if l.front != nil {
		t.Errorf("front should be nil")
	}

	if l.back != nil {
		t.Errorf("back should be nil")
	}
}
