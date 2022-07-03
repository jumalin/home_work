package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	lock     sync.RWMutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) (ok bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	var item *ListItem
	item, ok = c.items[key]
	if ok {
		item.Value.(*cacheItem).value = value
		c.queue.MoveToFront(item)
		return
	}
	newCacheItem := &cacheItem{key, value}
	item = c.queue.PushFront(newCacheItem)
	c.items[key] = item
	if len(c.items) > c.capacity {
		back := c.queue.Back()
		c.queue.Remove(back)
		delete(c.items, back.Value.(*cacheItem).key)
	}
	return
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	c.queue.MoveToFront(item)
	return item.Value.(*cacheItem).value, ok
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
