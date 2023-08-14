package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mutex    sync.RWMutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type lruItem struct {
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

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	newElement := lruItem{
		key:   key,
		value: value,
	}

	if element, exists := cache.items[key]; exists {
		element.Value = newElement
		cache.queue.MoveToFront(element)
		return true
	}

	if cache.queue.Len() == cache.capacity {
		lastElement := cache.queue.Back()
		cache.queue.Remove(lastElement)
		delete(cache.items, lastElement.Value.(lruItem).key)
	}
	cache.queue.PushFront(newElement)
	cache.items[key] = cache.queue.Front()
	return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	element, exists := cache.items[key]
	if !exists {
		return nil, false
	}

	cache.queue.MoveToFront(element)
	return element.Value.(lruItem).value, true
}

func (cache *lruCache) Clear() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	element := cache.queue.Back()
	for element != nil {
		delete(cache.items, element.Value.(lruItem).key)
		element = element.Next
	}
	cache.queue = NewList()
}
