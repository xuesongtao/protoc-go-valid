package internal

import "unsafe"

// UnsafeBytes2Str
func UnsafeBytes2Str(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

// UnsafeStr2Bytes
func UnsafeStr2Bytes(str string) []byte {
	return *(*[]byte)(unsafe.Pointer(&str))
}
