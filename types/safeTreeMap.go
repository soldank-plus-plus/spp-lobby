package types

import (
	"sync"

	"github.com/igrmk/treemap/v2"
	"golang.org/x/exp/constraints"
)

type SafeTreeMap[Key constraints.Ordered, Value any] struct {
	Mutex         sync.Mutex
	UnsafeTreeMap *treemap.TreeMap[Key, Value]
}

func NewSafeTreeMap[Key constraints.Ordered, Value any]() *SafeTreeMap[Key, Value] {
	return &SafeTreeMap[Key, Value]{
		UnsafeTreeMap: treemap.New[Key, Value](),
	}
}

func (self *SafeTreeMap[Key, Value]) Get(key Key) (Value, bool) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	val, found := self.UnsafeTreeMap.Get(key)
	return val, found
}

func (self *SafeTreeMap[Key, Value]) Set(key Key, val Value) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	self.UnsafeTreeMap.Set(key, val)
}

func (self *SafeTreeMap[Key, Value]) Del(key Key) {
	self.Mutex.Lock()
	defer self.Mutex.Unlock()

	self.UnsafeTreeMap.Del(key)
}
