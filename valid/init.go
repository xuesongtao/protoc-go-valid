package valid

import (
	"errors"
	"reflect"
	"strings"
	"sync"
)

// err msg 单位
const (
	strUnitStr      = "strLength"
	numUnitStr      = "numSize"
	sliceLenUnitStr = "sliceLen"
)

// 验证 tag
const (
	Required    = "required"   // 必填
	Exist       = "exist"      // 有值才验证
	Either      = "either"     // 多个必须一个
	BothEq      = "botheq"     // 两者相等
	VTo         = "to"         // 两者之间, 闭区间
	VGe         = "ge"         // 大于或等于
	VLe         = "le"         // 小于或等于
	VOTo        = "oto"        // 两者之间, 开区间
	VGt         = "gt"         // 大于
	VLt         = "lt"         // 小于
	VEq         = "eq"         // 等于
	VNoEq       = "noeq"       // 不等于
	VIn         = "in"         // 指定输入选项
	VInclude    = "include"    // 指定输入包含选项
	VPhone      = "phone"      // 手机号
	VEmail      = "email"      // 邮箱
	VIDCard     = "idcard"     // 身份证号码
	VYear       = "year"       // 年
	VYear2Month = "year2month" // 年月
	VDate       = "date"       // 日
	VDatetime   = "datetime"   // 日期+时间点
	VInt        = "int"        // 整数
	VFloat      = "float"      // 浮动数
	VRe         = "re"         // 正则
)

// 验证函数
var validName2FuncMap = map[string]CommonValidFn{
	Required:    nil,
	Exist:       nil,
	Either:      nil,
	BothEq:      nil,
	VTo:         To,
	VGe:         Ge,
	VLe:         Le,
	VOTo:        OTo,
	VGt:         Gt,
	VLt:         Lt,
	VEq:         Eq,
	VNoEq:       NoEq,
	VIn:         In,
	VInclude:    Include,
	VPhone:      Phone,
	VEmail:      Email,
	VIDCard:     IDCard,
	VYear:       Year,
	VYear2Month: Year2Month,
	VDate:       Date,
	VDatetime:   Datetime,
	VInt:        Int,
	VFloat:      Float,
	VRe:         Re,
}

// 对象
var (
	syncValidPool = sync.Pool{New: func() interface{} { return new(VStruct) }}
)

// 标记
var (
	defaultTargetTag = "valid" // 默认的验证 tag
	ErrEndFlag       = "; "    // 错误结束符号(每个自定义 err 都需要将这个追加在后面, 用于分句)
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

	reErr = errors.New(defaultTargetTag + " \"re\" is not ok, eg: " +
		"type Test struct {\n" +
		"    Age string `valid:\"re='\\\\d+'\"`\n" +
		"}")
)

// CommonValidFn 通用验证函数, 主要用于回调
// 注: 在写 errBuf 的时候建议用 GetJoinValidErrStr 包裹下, 这样产生的结果易读.
//     否则需要再 errBuf.Writestring 最后要加上 ErrEndFlag 分割, 工具是通过 ErrEndFlag 进行分句
type CommonValidFn func(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value)

// Deprecated 此函数会修改全局变量, 会导致内存释放不了, 此推荐 ValidStructForMyValidFn
// SetCustomerValidFn 自定义验证函数
func SetCustomerValidFn(validName string, fn CommonValidFn) {
	validName2FuncMap[validName] = fn
}

// JoinTag2Val 生成 defaultTargetTag 的值
// tag 为验证点
// values[0] 会被解析为值
// values[1] 会被解析为自定义错误信息
// 如: JoinTag2Val(VRe, "\\d+", "必须为纯数字")
// => re='\\d+'|必须为纯数字
func JoinTag2Val(tag string, values ...string) string {
	l := len(values)
	if l == 0 {
		return tag
	}

	tagVal := tag
	if values[0] != "" {
		needAddEqual := values[0][0] != '=' // 判断第一个值得首字符是否为 "="

		// 处理 val 前缀
		switch tag {
		case Either, BothEq, VTo, VGe, VLe, VOTo, VGt, VLt, VEq, VNoEq:
			if needAddEqual {
				tagVal += "="
			}
		case VIn, VInclude:
			if needAddEqual {
				tagVal += "="
			}
			tagVal += "("
		case VRe:
			if needAddEqual {
				tagVal += "="
			}
			tagVal += "'"
		}

		// 处理 val
		tagVal += values[0]

		// 处理 val 后缀
		switch tag {
		case VIn, VInclude:
			tagVal += ")"
		case VRe:
			tagVal += "'"
		}
	}

	// 自定义说明
	if l >= 2 {
		tagVal += "|" + values[1]
	}
	return tagVal
}
