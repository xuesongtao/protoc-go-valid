package valid

import (
	"container/list"
	"strings"
	"sync"
)

type LRUCache struct {
	rwMu    sync.RWMutex
	maxSize int
	nodeMap map[interface{}]*list.Element
	list    *list.List
}

func NewLRU(max int) *LRUCache {
	return &LRUCache{
		maxSize: max,
		nodeMap: make(map[interface{}]*list.Element),
		list:    list.New(),
	}
}

// Store
func (l *LRUCache) Store(key, value interface{}) {
	l.rwMu.Lock()
	defer l.rwMu.Unlock()
	node, ok := l.nodeMap[key]
	if ok {
		l.list.MoveToFront(node)
		return
	}
	
	front := l.list.PushFront(value)
	l.nodeMap[key] = front
	// 判断是否已满, 满了就删除最后一个
	if l.list.Len() >= l.maxSize {
		endNode := l.list.Back()
		delete(l.nodeMap, l.list.Remove(endNode))
	}
}

// Load
func (l *LRUCache) Load(key interface{}) (data interface{}, ok bool) {
	l.rwMu.RLock()
	defer l.rwMu.RUnlock()
	node, ok := l.nodeMap[key]
	if !ok {
		return
	}
	data = node.Value
	return
}

// Len
func (l *LRUCache) Len() int {
	l.rwMu.RLock()
	defer l.rwMu.RUnlock()
	return l.list.Len()
}

// Dump
func (l *LRUCache) Dump() string {
	head := l.list.Front()
	buf := new(strings.Builder)
	for head != nil {
		buf.WriteString(ToStr(head.Value))
		head = head.Next()
		if head != nil {
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}
