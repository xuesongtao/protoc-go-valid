package internal

import "unsafe"

// Bytes2Str 
func Bytes2Str(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}
