package internal

import (
	"fmt"
	"sync"
	"time"
)

type LRUCache struct {
	rwMu     sync.RWMutex
	maxSize  int
	cacheMap map[interface{}]*cacheData
}

type cacheData struct {
	timestamp int64
	data      interface{}
}

func NewLRU(max int) *LRUCache {
	return &LRUCache{
		maxSize:  max,
		cacheMap: make(map[interface{}]*cacheData, max),
	}
}

func (l *LRUCache) getTime() int64 {
	return time.Now().Unix()
}

func (l *LRUCache) Load(key interface{}) (data interface{}, ok bool) {
	l.rwMu.RLock()
	defer l.rwMu.RUnlock()

	var res *cacheData
	res, ok = l.cacheMap[key]
	if !ok {
		return
	}

	// 如果存在就修改下使用时间, 不考虑数据竞争
	res.timestamp = l.getTime()
	data = res.data
	l.cacheMap[key] = res
	return
}

func (l *LRUCache) Store(key, value interface{}) {
	l.rwMu.Lock()
	defer l.rwMu.Unlock()

	if len(l.cacheMap) >= l.maxSize {
		var (
			i            = -1
			delKey       interface{}
			minTimestamp int64
		)
		for key, data := range l.cacheMap {
			i++
			if i == 0 {
				minTimestamp = data.timestamp
				delKey = key
				continue
			}

			if minTimestamp > data.timestamp {
				minTimestamp = data.timestamp
				delKey = key
			}
		}
		// fmt.Printf("delKey: %v, timestamp: %d\n", delKey, minTimestamp)
		delete(l.cacheMap, delKey)
	}

	l.cacheMap[key] = &cacheData{
		timestamp: l.getTime(),
		data:      value,
	}
}

func (l *LRUCache) Len() int {
	l.rwMu.RLock()
	defer l.rwMu.RUnlock()
	return len(l.cacheMap)
}

func (l *LRUCache) Dump() string {
	return fmt.Sprintf("%+v", l.cacheMap)
}
