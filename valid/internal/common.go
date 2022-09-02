package internal

import "unsafe"

// Byte2Str 
func Byte2Str(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}
