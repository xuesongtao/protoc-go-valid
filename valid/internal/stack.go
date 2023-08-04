package internal

type stackByte struct {
	data []byte
}

func NewStackByte(size int) *stackByte {
	return &stackByte{
		data: make([]byte, 0, size),
	}
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
	s.data = nil
}
