package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		// Write me
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}

func TestClear(t *testing.T) {
	c := &lruCache{
		queue:    NewList(),
		items:    make(map[Key]*ListItem),
		capacity: 5,
	}

	c.items["key1"] = &ListItem{Value: "item1"}
	c.items["key2"] = &ListItem{Value: "item2"}

	if len(c.items) == 0 {
		t.Fatal("Cache is empty before Clear is called")
	}

	c.Clear()

	if len(c.items) != 0 {
		t.Error("Cache was not cleared, items should be empty")
	}

	if c.queue == nil {
		t.Error("Cache queue was not reset")
	}
}

func TestLRUCache_Capacity(t *testing.T) {
	capacity := 3
	c := NewCache(capacity)

	c.Set("key1", 1)
	c.Set("key2", 2)
	c.Set("key3", 3)
	c.Set("key4", 4)

	_, exists := c.Get("key1")
	if exists {
		t.Errorf("Expected key1 to be evicted, but it still exists")
	}

	_, exists = c.Get("key2")
	if !exists {
		t.Errorf("Expected key2 to exist, but it's missing")
	}

	_, exists = c.Get("key3")
	if !exists {
		t.Errorf("Expected key3 to exist, but it's missing")
	}

	_, exists = c.Get("key4")
	if !exists {
		t.Errorf("Expected key4 to exist, but it's missing")
	}
}

func TestLRUCache_LRU(t *testing.T) {
	capacity := 3
	cache := NewCache(capacity)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	cache.Get("key2")
	cache.Get("key1")
	cache.Get("key3")
	cache.Set("key4", "value4")

	_, exists := cache.Get("key2")
	if exists {
		t.Errorf("Expected key2 to be evicted, but it still exists")
	}

	_, exists = cache.Get("key1")
	if !exists {
		t.Errorf("Expected key1 to exist, but it's missing")
	}

	_, exists = cache.Get("key3")
	if !exists {
		t.Errorf("Expected key3 to exist, but it's missing")
	}

	_, exists = cache.Get("key4")
	if !exists {
		t.Errorf("Expected key4 to exist, but it's missing")
	}
}
