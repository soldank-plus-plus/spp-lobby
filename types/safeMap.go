package types

import "sync"

type SafeMap[Key comparable, Value any] struct {
	mutex     sync.Mutex
	unsafeMap map[Key]Value
}

func NewSafeMap[Key comparable, Value any]() *SafeMap[Key, Value] {
	return &SafeMap[Key, Value]{unsafeMap: make(map[Key]Value)}
}

func (self *SafeMap[Key, Value]) Get(key Key) (Value, bool) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	val, found := self.unsafeMap[key]
	return val, found
}

func (self *SafeMap[Key, Value]) Set(key Key, val Value) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.unsafeMap[key] = val
}

func (self *SafeMap[Key, Value]) Delete(key Key) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	delete(self.unsafeMap, key)
}
