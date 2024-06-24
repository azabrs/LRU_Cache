package lrucache

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func Test(t *testing.T) {
	t.Run("Put and Get test", func(t *testing.T) {
		cache, err := New(10)
		if err != nil{
			t.Error(err)
		}
		cache.Add("first key", "first value")
		val, ok := cache.Get("first key")
		if !ok || val != "first value"{
			t.Errorf("read value does not match the written value")
		}
	})
	t.Run("Wrong capacity", func(t *testing.T) {
		_, err1 := New(-2)
		err2 := fmt.Errorf("capacity must be positive")
		if !reflect.DeepEqual(err1, err2){
			t.Error(err1)
		}
	})
	t.Run("Overflow test", func(t *testing.T) {
		cache, err := New(2)
		if err != nil{
			t.Error(err)
		}
		cache.Add("first key", 1)
		cache.Add("second key", 2)
		cache.Add("third key", 3)
		if _, ok := cache.Get("first key"); ok{
			t.Errorf("There was no deletion")
		}
	})
	t.Run("re-writing test", func(t *testing.T) {
		cache, err := New(2)
		if err != nil{
			t.Error(err)
		}
		cache.Add("first key", 1)
		cache.Add("first key", 6)
		if val, ok := cache.Get("first key"); !ok || val != 6{
			t.Errorf("There was no re-writing")
		}
	})
	t.Run("invalid key test", func(t *testing.T) {
		cache, err := New(2)
		if err != nil{
			t.Error(err)
		}
		if val, ok := cache.Get("first key"); ok || val != nil{
			t.Errorf("Invalid value")
		}
	})

	t.Run("clear test", func(t *testing.T) {
		cache, err := New(2)
		if err != nil{
			t.Error(err)
		}
		cache.Add(1, "first value")
		cache.Add(2, "second value")
		cache.Clear()
		if val, ok := cache.Get(1); ok || val != nil{
			t.Errorf("cache was not cleared")
		}
		if val, ok := cache.Get(2); ok || val != nil{
			t.Errorf("cache was not cleared")
		}
	})

	t.Run("remove test", func(t *testing.T) {
		cache, err := New(5)
		if err != nil{
			t.Error(err)
		}
		cache.Add(1, "first value")
		cache.Remove(1)
		if val, ok := cache.Get(1); ok || val != nil{
			t.Errorf("item has not been deleted")
		}
	})

	t.Run("addTTL test", func(t *testing.T) {
		cache, err := New(5)
		if err != nil{
			t.Error(err)
		}
		cache.AddWithTTL(1, "first value", time.Second * 1)
		if val, ok := cache.Get(1); !ok || val == nil{
			t.Errorf("item has not been added")
		}
		time.Sleep(time.Second * 2)
		if val, ok := cache.Get(1); ok || val != nil{
			t.Errorf(" item was not deleted after the expiration date")
		}
	})

	t.Run("re-writing with ttl test", func(t *testing.T) {
		cache, err := New(5)
		if err != nil{
			t.Error(err)
		}
		cache.AddWithTTL(1, "first value", time.Second * 5)
		if val, ok := cache.Get(1); !ok || val == nil{
			t.Errorf("item has not been added")
		}
		cache.AddWithTTL(1, "first value", time.Second * 1)
		time.Sleep(time.Second * 2)
		if val, ok := cache.Get(1); ok || val != nil{
			t.Errorf(" item was not deleted")
		}
	})

	t.Run("element overflow test with ttl", func(t *testing.T) {
		cache, err := New(2)
		if err != nil{
			t.Error(err)
		}
		cache.AddWithTTL(1, "first value", time.Second * 10)
		cache.AddWithTTL(2, "second value", time.Second * 3)
		cache.AddWithTTL(3, "third valud", time.Second * 2)
		if val, ok := cache.Get(1); ok || val != nil{
			t.Errorf("item has not been deleted")
		}

	})
}
