package pokecache

import (
	"testing"
	"time"
)

func TestCreateCache(t *testing.T) {
	cache := NewCache(time.Second * 5)
	if cache.cacheMap == nil {
		t.Error("Cache is nil")
	}
}

func TestAddGetCache(t *testing.T) {
	cache := NewCache(time.Second * 5)
	if cache.cacheMap == nil {
		t.Error("Cache is nil")
	}
	data := []byte{'a','b','c'}
	cache.Add("entry_one", data)
	out, ok := cache.Get("entry_one")
	if !ok {
		t.Error("Cache didn't find entry")
	}
	for i, v := range out {
		if v != data[i] {
			t.Errorf("Entry doesn't match proper data: \nExpected: %v\nActual: %v",data,out)
		}
	}
}

func TestReapCache(t *testing.T) {
	cache := NewCache(time.Millisecond*20)
	if cache.cacheMap == nil {
		t.Error("Cache is nil")
	}
	cache.Add("entry1", []byte("abc"))
	time.Sleep(time.Millisecond*25)
	_, ok := cache.Get("entry1")
	if ok {
		t.Error("Cache failed to reap entry")
	}
}