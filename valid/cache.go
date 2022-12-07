package valid

import (
	"container/list"
	"sync"
)

const (
	lruSize = 2 << 8 // lru 最大值
)

type LRUCache struct {
	rwMu        sync.RWMutex
	maxSize     int
	delMapCount int // 记录 delete map 的次数, 当次数大于 2*lruSize 重建下 nodeMap, 防止 delete 没有释放内存
	nodeMap     map[interface{}]*list.Element
	list        *list.List
}

func NewLRU(max int) *LRUCache {
	return &LRUCache{
		maxSize: max,
		nodeMap: make(map[interface{}]*list.Element),
		list:    list.New(),
	}
}

func (l *LRUCache) Store(key, value interface{}) {
	l.rwMu.Lock()
	defer l.rwMu.Unlock()
	node, ok := l.nodeMap[key]
	if ok {
		l.list.MoveToFront(node)
		return
	}

	// 不存在
	head := l.list.PushFront(value)
	l.nodeMap[key] = head
	// 判断是否已满, 满了就删除最后一个
	if l.list.Len() > l.maxSize {
		delete(l.nodeMap, l.list.Remove(l.list.Back()))
		l.delMapCount++

		// 重建 map
		if l.delMapCount > 2*lruSize {
			tmp := l.nodeMap
			l.nodeMap = make(map[interface{}]*list.Element, len(tmp))
			for k, v := range tmp {
				l.nodeMap[k] = v
			}
			l.delMapCount = 0
		}
	}
}

func (l *LRUCache) Load(key interface{}) (data interface{}, ok bool) {
	l.rwMu.Lock()
	defer l.rwMu.Unlock()
	node, ok := l.nodeMap[key]
	if !ok {
		return
	}
	data = node.Value
	l.list.MoveToFront(node)
	return
}

func (l *LRUCache) Len() int {
	l.rwMu.RLock()
	defer l.rwMu.RUnlock()
	return l.list.Len()
}

func (l *LRUCache) Dump() string {
	head := l.list.Front()
	buf := newStrBuf()
	defer putStrBuf(buf)
	for head != nil {
		buf.WriteString(ToStr(head.Value))
		head = head.Next()
		if head != nil {
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}
