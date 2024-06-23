package main

import (
	lrucache "LRU_Cache/LRU_Cache"
	"fmt"
	"log"
	"time"
)

func main(){
	cache, err := lrucache.New(3)
	if err != nil{
		log.Println(err)
	}
	fmt.Println(cache.Cap())
	cache.Add(6, "second")
	cache.AddWithTTL(5, "first", time.Second * 3)
	cache.Add(3, "third")
	cache.Add(2, "fourth")
	fmt.Println(cache.Get(5))
	time.Sleep(time.Second * 4)
	fmt.Println(cache.Get(5))
}