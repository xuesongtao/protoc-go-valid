package internal

import (
	"sync"

	"gitee.com/xuesongtao/protoc-go-valid/log"
)

var (
	cache = sync.Pool{
		New: func() interface{} {
			return &stackByte{}
		},
	}
)

type stackByte struct {
	data []byte
}

func NewStackByte(size int) *stackByte {
	obj := cache.Get().(*stackByte)
	if obj.data == nil || cap(obj.data) > 1<<8 {
		obj.data = make([]byte, 0, size)
	} else {
		obj.data = obj.data[:0]
	}
	return obj
}

func (s *stackByte) Append(b byte) {
	s.data = append(s.data, b)
}

func (s *stackByte) IsEmpty() bool {
	return len(s.data) == 0
}

func (s *stackByte) Pop() byte {
	if s.IsEmpty() {
		return byte(' ')
	}
	lastIndex := len(s.data) - 1
	b := s.data[lastIndex]
	if lastIndex >= 1 {
		s.data = append(s.data[:0], s.data[:lastIndex-1]...)
	} else {
		s.data = s.data[:0]
	}
	return b
}

func (s *stackByte) LastVal() byte {
	if s.IsEmpty() {
		return byte(' ')
	}
	return s.data[len(s.data)-1]
}

func (s *stackByte) IsEqualLastVal(b byte) bool {
	return s.LastVal() == b
}

func (s *stackByte) Reset() {
	cache.Put(s)
}

func (s *stackByte) Dump() {
	log.Infof("%+v", s.data)
}
