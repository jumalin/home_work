package hw04lrucache

import (
	"sync"
)

type Key string

var lock = sync.RWMutex{}

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
	lock.Lock()
	defer lock.Unlock()
	_, ok = c.items[key]
	if ok {
		elem := c.items[key].Value
		ci := elem.(cacheItem)
		ci.value = value
		return
	}
	c.queue.PushFront(cacheItem{
		key:   key,
		value: value,
	})

	c.items[key] = c.queue.Front()
	if c.queue.Len() > c.capacity {
		rv := c.queue.Back().Value
		c.items[rv.(cacheItem).key] = nil
		c.queue.Remove(c.queue.Back())
	}
	return
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	lock.Lock()
	defer lock.Unlock()
	val, ok := c.items[key]
	if !ok {
		return nil, false
	}
	c.queue.MoveToFront(c.items[key])

	c.items[key] = c.queue.Front()
	item := val.Value.(cacheItem)
	return item.value, ok
}

func (c *lruCache) Clear() {
	c = &lruCache{
		capacity: c.capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, c.capacity),
		//lock:     sync.RWMutex{},
	}
}
