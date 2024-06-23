package lrucache

import (
	"container/list"
	"fmt"
)

type LRU_Cache struct{
	capacity int
	l *list.List
	h map[*list.Element] interface{}
}

func New(cap int) (LRU_Cache, error){
	if cap <= 0{
		return LRU_Cache{}, fmt.Errorf("Capacity must be positive")
	}
	return LRU_Cache{
		capacity: cap,
		l: list.New(),
		h: make(map[*list.Element]interface{}),
	}, nil
}
