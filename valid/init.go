package valid

import (
	"errors"
	"reflect"
	"strings"
	"sync"
)

// 对象
var (
	syncValidPool     = sync.Pool{New: func() interface{} { return new(vStruct) }}
	validName2FuncMap map[string]commonValidFn
)

// 标记
var (
	defaultTargetTag = "valid" // 默认的验证 tag
	errEndFlag       = " \n"   // 错误结束符号
)

// 错误
var (
	toValErr     = errors.New(defaultTargetTag + " \"to\" is not ok, eg: to=1/to=6~30")
	eitherValErr = errors.New(defaultTargetTag + " \"either\" is not ok, eg: " +
		"type Test struct {\n" +
		"    OrderNo string `either=1`\n" +
		"    TradeNo sting `either=1`\n" +
		"}, errMsg: \"OrderNo\" either \"TradeNo\" they shouldn't all be empty")
	inValErr = errors.New(defaultTargetTag + " \"in\" is not ok, eg: in=(1,2,3)")
)

type commonValidFn func(errBuf *strings.Builder, validName, structName, filedName string, tv reflect.Value)

func init() {
	validName2FuncMap = map[string]commonValidFn{
		"required": nil,
		"either":   nil,
		"to":       To,
		"in":       In,
		"phone":    Phone,
		"email":    Email,
		"idcard":   IDCard,
	}
}

// GetValidFn 获取验证函数
func GetValidFn(validName string) (commonValidFn, error) {
	f, ok := validName2FuncMap[validName]
	if !ok {
		return nil, errors.New("valid: \"" + validName + "\" is not exist, I am sorry")
	}
	return f, nil
}
