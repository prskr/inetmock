package mock

import (
	"reflect"
	"runtime"
	"sync"
)

func NewCallCounter() *CallCounter {
	return &CallCounter{
		lock:  sync.RWMutex{},
		index: make(map[string]int),
	}
}

type CallCounter struct {
	lock  sync.RWMutex
	index map[string]int
}

func (cc *CallCounter) Inc(i interface{}) {
	cc.lock.Lock()
	defer cc.lock.Unlock()
	funcName := funcName(i)
	if currentCount, ok := cc.index[funcName]; ok {
		cc.index[funcName] = currentCount + 1
	} else {
		cc.index[funcName] = 1
	}
}

func (cc *CallCounter) GetValueFor(i interface{}) int {
	return cc.getByName(funcName(i))
}

func (cc *CallCounter) getByName(name string) int {
	cc.lock.RLock()
	defer cc.lock.RUnlock()
	if val, ok := cc.index[name]; ok {
		return val
	}
	return 0
}

func funcName(i interface{}) string {
	reflect.ValueOf(i)
	funcEntry := runtime.FuncForPC(reflect.ValueOf(i).Pointer())
	return funcEntry.Name()
}
