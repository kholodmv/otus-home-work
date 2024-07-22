package hw04lrucache

type Key string

type Cache interface {
	// Set adds a value to the cache by key.
	Set(key Key, value interface{}) bool
	// Get returns value from cache by key.
	Get(key Key) (interface{}, bool)
	// Clear clears cache.
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

// NewCache creates a new LRU Cache with the specified capacity.
func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

// Set adds a value to the cache by key.
func (c *lruCache) Set(key Key, value interface{}) bool {
	if item, exists := c.items[key]; exists {
		c.queue.MoveToFront(item)
		item.Value = &ListItem{Value: cacheItem{key: key, value: value}}
		return true
	}

	newItem := c.queue.PushFront(&ListItem{Value: cacheItem{key: key, value: value}})
	c.items[key] = newItem

	if c.queue.Len() > c.capacity {
		delete(c.items, c.queue.Back().Value.(*ListItem).Value.(cacheItem).key)
		c.queue.Remove(c.queue.Back())
	}

	return false
}

// Get gets the value from the cache by key.
func (c *lruCache) Get(key Key) (interface{}, bool) {
	if item, exists := c.items[key]; exists {
		c.queue.MoveToFront(item)

		return item.Value.(*ListItem).Value.(cacheItem).value, true
	}
	return nil, false
}

// Clear clears the entire cache.
func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
