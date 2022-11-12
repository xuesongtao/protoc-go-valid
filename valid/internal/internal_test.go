package internal

import (
	"testing"
)

func TestStack(t *testing.T) {
	a := "required|必填,phone|'手机号码必填,同时正确',re='\\d+{1,2}'"
	stack := NewStackByte(2)
	defer stack.Reset()
	tmp := make([]byte, 0, 5)
	res := make([]string, 0, 3)
	isParseSingleQuotes := false
	for i := 0; i < len(a); i++ {
		v := a[i]
		if !isParseSingleQuotes && v != ',' {
			tmp = append(tmp, v)
		} else if isParseSingleQuotes { // 如果是单引号就不处理
			tmp = append(tmp, v)
		}
		if !isParseSingleQuotes && v == '\'' {
			stack.Append(v)
			isParseSingleQuotes = true
			continue
		}
		if isParseSingleQuotes && stack.IsEqualLastVal(v) {
			stack.Pop()
			isParseSingleQuotes = false
		}

		if v == ',' && stack.IsEmpty() {
			res = append(res, string(tmp))
			tmp = tmp[:0]
		}
	}
	// stack.Dump()
	if len(tmp) > 0 {
		res = append(res, string(tmp))
	}
	t.Log(res)
}

func TestByte2Str(t *testing.T) {
	t.Log(UnsafeBytes2Str([]byte("hello 你好 a")))
}

func BenchmarkByte2Str1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = UnsafeBytes2Str([]byte("hello 你好"))
	}
}

func BenchmarkByte2Str2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = string([]byte("hello 你好"))
	}
}
