package valid

import (
	"container/list"
	"sync"
)

const (
	lruSize = 2 << 8 // lru 最大值
)

type LRUCache struct {
	rwMu             sync.RWMutex
	maxSize          int
	delMapCount      int // 记录 delete map 的次数, 当次数大于 2*lruSize 重建下 nodeMap, 防止 delete 没有释放内存
	nodeMap          map[interface{}]*list.Element
	list             *list.List
	deleteCallBackFn func(key, value interface{}) // 删除回调
}

func NewLRU(max ...int) *LRUCache {
	defaultMax := lruSize
	if len(max) > 0 {
		defaultMax = max[0]
	}

	return &LRUCache{
		maxSize: defaultMax,
		nodeMap: make(map[interface{}]*list.Element, defaultMax),
		list:    list.New(),
	}
}

func (l *LRUCache) SetDelCallBackFn(f func(key, value interface{})) {
	l.deleteCallBackFn = f
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
		l.delete(l.list.Back())
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

func (l *LRUCache) Delete(key interface{}) {
	l.rwMu.Lock()
	defer l.rwMu.Unlock()
	node, ok := l.nodeMap[key]
	if !ok {
		return
	}
	l.delete(node)
}

func (l *LRUCache) delete(node *list.Element) {
	var key interface{}
	for k, v := range l.nodeMap {
		if v == node {
			key = k
			break
		}
	}

	delete(l.nodeMap, key)
	l.list.Remove(node)
	if l.deleteCallBackFn != nil {
		l.deleteCallBackFn(key, node.Value)
	}

	// 重建 map
	if l.delMapCount > 2*l.maxSize {
		tmp := l.nodeMap
		l.nodeMap = make(map[interface{}]*list.Element, len(tmp))
		for k, v := range tmp {
			l.nodeMap[k] = v
		}
		l.delMapCount = 0
	} else {
		l.delMapCount++
	}
}

// Len 长度
// return -1 的话, 长度不正确
func (l *LRUCache) Len() int {
	l.rwMu.RLock()
	defer l.rwMu.RUnlock()
	if l.list.Len() != len(l.nodeMap) {
		return -1
	}
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
