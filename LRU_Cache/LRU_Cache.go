package lrucache

import (
	"container/list"
	"fmt"
)

type LRU_Cache struct{
	capacity int
	l *list.List
	h map[interface{}] *list.Element
}

type node struct{
	key interface{}
	value interface{}
}

func New(cap int) (LRU_Cache, error){
	if cap <= 0{
		return LRU_Cache{}, fmt.Errorf("capacity must be positive")
	}
	return LRU_Cache{
		capacity: cap,
		l: list.New(),
		h: make(map[interface{}] *list.Element),
	}, nil
}

func (c LRU_Cache)Cap() int{
	return c.capacity
}

func (c LRU_Cache)Add(key, value interface{}){
	if val, ok := c.h[key]; !ok{
		c.l.Remove(val)
	}
	node := c.l.PushBack(node{
		key: key,
		value: value,
	})
	c.h[key] = node
	if c.l.Len() > c.Cap(){
		c.l.Remove(c.l.Front())
	}
}


