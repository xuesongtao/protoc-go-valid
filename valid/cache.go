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
	// 判断下是否存在
	l.rwMu.Lock()
	defer l.rwMu.Unlock()
	if l.list.Len() >= l.maxSize {
		// 删除最后一个 node
		endNode := l.list.Back()
		// 删除 map
		delete(l.nodeMap, l.list.Remove(endNode))
	}

	node, ok := l.nodeMap[key]
	if !ok {
		// 不存在就新建一个放到 头部
		front := l.list.PushFront(value)
		l.nodeMap[key] = front
		return
	}
	// 将节点移到 头部
	l.list.MoveToFront(node)
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
