package valid

import (
	"errors"
	"reflect"
	"strings"
	"sync"
)

const (
	strUnitStr      = "strLength"
	numUnitStr      = "numSize"
	sliceLenUnitStr = "sliceLen"
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

	otoValErr = errors.New(defaultTargetTag + " \"to\" is not ok, eg: " +
		"type Test struct {\n" +
		"    Name string `valid:\"oto=1~10\"`\n" +
		"}")

	eitherValErr = errors.New(defaultTargetTag + " \"either\" is not ok, eg: " +
		"type Test struct {\n" +
		"    OrderNo string `valid:\"either=1\"`\n" +
		"    TradeNo sting `valid:\"either=1\"`\n" +
		"}, errMsg: \"OrderNo\" either \"TradeNo\" they shouldn't all be empty")

	bothEqValErr = errors.New(defaultTargetTag + " \"botheq\" is not ok, eg: " +
		"type Test struct {\n" +
		"    OrderNo string `valid:\"botheq=1\"`\n" +
		"    TradeNo sting `valid:\"botheq=1\"`\n" +
		"}, errMsg: \"OrderNo\" either \"TradeNo\" they shouldn't is no equal")

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
var validName2FuncMap = map[string]CommonValidFn{
	Required:     nil,
	Exist:        nil,
	Either:       nil,
	BothEq:       nil,
	"to":         To,
	"ge":         Ge,
	"le":         Le,
	"oto":        OTo,
	"gt":         Gt,
	"lt":         Lt,
	"eq":         Eq,
	"noeq":       NoEq,
	"in":         In,
	"include":    Include,
	"phone":      Phone,
	"email":      Email,
	"idcard":     IDCard,
	"year":       Year,
	"year2month": Year2Month,
	"date":       Date,
	"datetime":   Datetime,
	"int":        Int,
	"float":      Float,
}

type CommonValidFn func(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value)

// Deprecated 此函数会修改全局变量, 会导致内存释放不了, 此推荐 ValidStructForMyValidFn
// SetCustomerValidFn 自定义验证函数
func SetCustomerValidFn(validName string, fn CommonValidFn) {
	validName2FuncMap[validName] = fn
}
