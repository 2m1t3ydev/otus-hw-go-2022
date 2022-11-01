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
	mutex    sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (lru *lruCache) Set(key Key, value interface{}) bool {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	cacheItem := cacheItem{
		key:   key,
		value: value,
	}

	item, ok := lru.items[key]

	if ok {
		item.Value = cacheItem
		lru.queue.MoveToFront(item)
		return true
	}

	if lru.queue.Len() == lru.capacity {
		last := lru.queue.Back()
		lru.queue.Remove(last)
		delete(lru.items, key)
	}

	lru.items[key] = lru.queue.PushFront(cacheItem)

	return false
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	item, ok := lru.items[key]

	if ok {
		lru.queue.MoveToFront(item)
		return item.Value.(cacheItem).value, true
	}

	return nil, false
}

func (lru *lruCache) Clear() {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	lru.items = make(map[Key]*ListItem, lru.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
