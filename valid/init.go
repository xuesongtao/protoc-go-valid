package valid

import (
	"errors"
	"reflect"
	"strings"
	"sync"
)

const (
	strUnitStr = "length"
	numUnitStr = "size"
)

// 对象
var (
	syncValidPool = sync.Pool{New: func() interface{} { return new(vStruct) }}
)

// 标记
var (
	defaultTargetTag = "valid" // 默认的验证 tag
	errEndFlag       = "; "    // 错误结束符号
)

// 错误
var (
	toValErr = errors.New(defaultTargetTag + " \"to\" is not ok, eg: " +
		"type Test struct {\n" +
		"    Name string `valid:\"to=1~10\"`\n" +
		"}")

	eitherValErr = errors.New(defaultTargetTag + " \"either\" is not ok, eg: " +
		"type Test struct {\n" +
		"    OrderNo string `valid:\"either=1\"`\n" +
		"    TradeNo sting `valid:\"either=1\"`\n" +
		"}, errMsg: \"OrderNo\" either \"TradeNo\" they shouldn't all be empty")

	inValErr = errors.New(defaultTargetTag + " \"in\" is not ok, eg: " +
		"type Test struct {\n" +
		"   hobby int `valid:\"in=(1/2/3)\"`\n" +
		"}")

	includeErr = errors.New(defaultTargetTag + " \"include\" is not ok, filed type must is string, eg: " +
		"type Test struct {\n" +
		"    Name string `valid:\"include=(ab/cd)\"`\n" +
		"}")
)

// 验证函数
var validName2FuncMap = map[string]commonValidFn{
	"required": nil,
	"either":   nil,
	"to":       To,
	"ge":       Ge,
	"le":       Le,
	"oto":      OTo,
	"gt":       Gt,
	"lt":       Lt,
	"eq":       Eq,
	"noeq":     NoEq,
	"in":       In,
	"include":  Include,
	"phone":    Phone,
	"email":    Email,
	"idcard":   IDCard,
	"date":     Date,
	"datetime": Datetime,
	"int":      Int,
	"float":    Float,
}

type commonValidFn func(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value)

// SetCustomerValidFn 自定义验证函数
func SetCustomerValidFn(validName string, fn commonValidFn) {
	validName2FuncMap[validName] = fn
}

// GetValidFn 获取验证函数
func GetValidFn(validName string) (commonValidFn, error) {
	f, ok := validName2FuncMap[validName]
	if !ok {
		return nil, errors.New("valid: \"" + validName + "\" is not exist, You can call SetCustomerValidFn")
	}
	return f, nil
}
