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
	if node, ok := c.h[key]; !ok{
		c.l.Remove(node)
	}
	node := c.l.PushBack(value)
	c.h[key] = node
	if c.l.Len() > c.Cap(){
		c.l.Remove(c.l.Front())
	}
}

func (c LRU_Cache)Get(key interface{}) (value interface{}, ok bool){
	node, ok := c.h[key]
	if !ok{
		return nil, false
	}
	c.l.MoveToBack(node)
	
	return node.Value.(int), true
}

func (c LRU_Cache)Clear(){
	for k := range c.h {
		delete(c.h, k)
	}
	c.l.Init()
}

func (c LRU_Cache)Remove(key interface{}) error{
	node, ok := c.h[key]
	if !ok{
		return fmt.Errorf("element does not exist")
	}
	c.l.Remove(node)
	delete(c.h, key)
	return nil
}