package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Add a value to the cache by key
	Get(key Key) (interface{}, bool)     // Get value from cache by key
	Clear()                              // Clear cache
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

// NewCache creates a new LRU Cache with the specified capacity.
func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),                         // Initialize an empty doubly linked list.
		items:    make(map[Key]*ListItem, capacity), // Initialize an empty dictionary to store keys and references to list elements.
	}
}

// Set adds a value to the cache by key.
func (c *lruCache) Set(key Key, value interface{}) bool {
	if item, exists := c.items[key]; exists {
		c.queue.MoveToFront(item)
		item.Value = value
		return true
	}

	newItem := c.queue.PushFront(value)
	c.items[key] = newItem

	if c.queue.Len() > c.capacity {
		lastItem := c.queue.Back()
		if lastItem != nil {
			delete(c.items, lastItem.Value.(Key))
			c.queue.Remove(lastItem)
		}
	}
	return false
}

// Get gets the value from the cache by key.
func (c *lruCache) Get(key Key) (interface{}, bool) {
	if item, exists := c.items[key]; exists {
		c.queue.MoveToFront(item)
		return item.Value, true
	}
	return nil, false
}

// Clear clears the entire cache.
func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
