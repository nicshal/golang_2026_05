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
	mu       sync.Mutex
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	i, ok := l.items[key]
	if ok {
		i.Value = value
		l.queue.MoveToFront(i)
		return true
	}

	if l.queue.Len() == l.capacity {
		i := l.queue.Back()
		for k, v := range l.items {
			if v == i {
				delete(l.items, k)
			}
		}
		l.queue.Remove(i)
	}

	l.items[key] = l.queue.PushFront(value)
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	i, ok := l.items[key]
	if ok {
		l.queue.MoveToFront(i)
		return i.Value, true
	}

	return nil, false
}

func (l *lruCache) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.items = make(map[Key]*ListItem, l.capacity)
	l.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		mu:       sync.Mutex{},
	}
}
