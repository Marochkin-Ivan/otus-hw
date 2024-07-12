package hw04lrucache

type Key string

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

type cacheItem struct {
	key   Key
	value any
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func newCacheItem(key Key, value any) *cacheItem {
	return &cacheItem{
		key:   key,
		value: value,
	}
}

func (c *lruCache) Set(key Key, value any) bool {
	// если элемент присутствует в словаре, то обновить его значение и переместить элемент в начало очереди
	if item, exist := c.items[key]; exist {
		item.Value.(*cacheItem).value = value
		c.queue.MoveToFront(item)

		// возвращаемое значение - флаг, присутствовал ли элемент в кэше
		return true
	}

	// если элемента нет в словаре, то добавить в словарь и в начало очереди
	newItem := c.queue.PushFront(newCacheItem(key, value))
	c.items[key] = newItem

	// если размер очереди больше ёмкости кэша
	if c.queue.Len() > c.capacity {
		// необходимо удалить последний элемент из очереди и его значение из словаря
		deleteItem := c.queue.Back()
		c.queue.Remove(deleteItem)
		delete(c.items, deleteItem.Value.(*cacheItem).key)
	}

	// возвращаемое значение - флаг, присутствовал ли элемент в кэше
	return false
}

func (c *lruCache) Get(key Key) (any, bool) {
	// если элемент присутствует в словаре
	if item, exist := c.items[key]; exist {
		// переместить элемент в начало очереди
		c.queue.MoveToFront(item)

		// вернуть его значение и true
		return item.Value.(*cacheItem).value, true
	}

	// если элемента нет в словаре, то вернуть nil и false
	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
