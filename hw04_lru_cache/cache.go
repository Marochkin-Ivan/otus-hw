package hw04lrucache

import "encoding/json"

type Key string

type CacheElement struct {
	Key   Key
	Value any
}

type Cache interface {
	Set(key Key, value any) bool
	Get(key Key) (any, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value any) bool {
	v, _ := json.Marshal(CacheElement{key, value})

	if item, ok := c.items[key]; ok {
		item.Value = v
		c.queue.MoveToFront(item)
		return true
	}

	item := c.queue.PushFront(v)
	c.items[key] = item

	if c.queue.Len() > c.capacity {
		deleteItem := c.queue.Back()

		var el CacheElement
		_ = json.Unmarshal(deleteItem.Value.([]byte), &el)

		delete(c.items, el.Key)
	}

	return false
}

func (c *lruCache) Get(key Key) (any, bool) {
	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)

		var el CacheElement
		_ = json.Unmarshal(item.Value.([]byte), &el)

		return el.Value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
