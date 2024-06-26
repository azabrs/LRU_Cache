package lrucache

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)


type LRU_Cache struct{
	capacity int
	l *list.List
	h map[interface{}] *list.Element
	m sync.RWMutex
}

type NodeValue struct{
	key interface{}
	value interface{}
	intCh chan struct{} //interrupt channel. It is used to interrupt the waiting of elements with ttl
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

func (c *LRU_Cache)Cap() int{
	return c.capacity
}

func (c *LRU_Cache)Add(key, value interface{}){
	c.m.Lock()
	defer c.m.Unlock()
	if node, ok := c.h[key]; ok{ //Check whether an element with such key exists
		
		nodeVal := c.l.Remove(node)
		close(nodeVal.(NodeValue).intCh)
	}
	node := c.l.PushBack(NodeValue{
		key: key,
		value: value,
		intCh: make(chan struct{}),
	})
	c.h[key] = node
	if c.l.Len() > c.Cap(){//checking the cache for overflow
		firstNode := c.l.Front() //push out the element that has not been used the longest
		nodeVal := c.l.Remove(firstNode)
		close(nodeVal.(NodeValue).intCh)
		delete(c.h, nodeVal.(NodeValue).key)
	}
}

func (c *LRU_Cache)Get(key interface{}) (value interface{}, ok bool){
	c.m.Lock()
	defer c.m.Unlock()
	node, ok := c.h[key]
	if !ok{
		return nil, false
	}
	c.l.MoveToBack(node) //Moving the item to the end of the doubly linked list
	
	return node.Value.(NodeValue).value, true
}

func (c *LRU_Cache)Clear(){
	c.m.Lock()
	defer c.m.Unlock()
	for k := range c.h {
		close(c.h[k].Value.(NodeValue).intCh)
		delete(c.h, k)
		
	}
	c.l.Init()
	
}

func (c *LRU_Cache)Remove(key interface{}){
	c.m.Lock()
	defer c.m.Unlock()
	node, ok := c.h[key]
	if !ok{
		return
	}
	nodeVal := c.l.Remove(node)
	close(nodeVal.(NodeValue).intCh)
	delete(c.h, key)
	
}

func (c *LRU_Cache)AddWithTTL(key, value interface{}, ttl time.Duration){
	c.m.Lock()
	defer c.m.Unlock()
	if node, ok := c.h[key]; ok{ //Check whether an element with such key exists
		node := c.l.Remove(node)
		close(node.(NodeValue).intCh)
	}
	node := c.l.PushBack(NodeValue{
		key: key,
		value: value,
		intCh: make(chan struct{}),
	})
	c.h[key] = node
	if c.l.Len() > c.Cap(){
		firstNode := c.l.Front()
		nodeVal := c.l.Remove(firstNode)
		close(nodeVal.(NodeValue).intCh)
		delete(c.h, nodeVal.(NodeValue).key)
	}
	go func(ttl time.Duration, key interface{}){
		
		select{
		case <-time.After(ttl):
			c.m.Lock()
			defer c.m.Unlock()
			c.l.Remove(node)
			delete(c.h, node.Value.(NodeValue).key)
		case <- node.Value.(NodeValue).intCh:
			return
		}
		
	}(ttl, key)
}