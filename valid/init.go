package valid

import (
	"errors"
	"reflect"
	"strings"
	"sync"
	"time"
)

const (
	// err msg 单位
	strUnitStr      = "str-length"
	numUnitStr      = "num-size"
	sliceLenUnitStr = "slice-len"

	// err msg 前缀
	ExplainEn = "explain:"
	ExplainZh = "说明:"
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
	VInts       = "ints"       // 多个数字验证
	VFloat      = "float"      // 浮动数
	VRe         = "re"         // 正则
	VIpv4       = "ipv4"       // ipv4
	VUnique     = "unique"     // 唯一验证
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
	VInts:       Ints,
	VFloat:      Float,
	VRe:         Re,
	VIpv4:       Ipv4,
	VUnique:     Unique,
}

// 对象
var (
	syncValidPool   = sync.Pool{New: func() interface{} { return new(VStruct) }}
	timeReflectType = reflect.TypeOf(time.Time{})
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

	intsErr = errors.New(defaultTargetTag + " \"ints\" is not ok, eg: " +
		"type Test struct {\n" +
		"    Hobby1 string `valid:\"ints\"`\n" + // 默认按 "," 进行分割对字符串进行判断是否为整数
		"    Hobby2 string `valid:\"ints=-\"`\n" + // 按 "-" 进行分割对字符串进行判断是否为整数
		"    Hobby3 []string `valid:\"ints\"`\n" + // 遍历切片中的元素是否为整数
		"}")

	uniqueErr = errors.New(defaultTargetTag + " \"unique\" is not ok, eg: " +
		"type Test struct {\n" +
		"    Hobby1 string `valid:\"unique\"`\n" + // 按 "," 进行分割对字符串进行判断是否唯一
		"    Hobby2 []string `valid:\"ints\"`\n" + // 遍历切片中的元素是否唯一
		"    Hobby3 []int `valid:\"ints\"`\n" + // 遍历切片中的元素是否唯一
		"}")
)

// CommonValidFn 通用验证函数, 主要用于回调
// 注: 在写 errBuf 的时候建议用 GetJoinValidErrStr 包裹下, 这样产生的结果易读.
//     否则需要再 errBuf.Writestring 最后要加上 ErrEndFlag 分割, 工具是通过 ErrEndFlag 进行分句
type CommonValidFn func(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value)

// SetCustomerValidFn 自定义验证函数
// Deprecated 此函数会修改全局变量, 会导致内存释放不了, 此推荐 *VStruct.SetValidFn
func SetCustomerValidFn(validName string, fn CommonValidFn) {
	validName2FuncMap[validName] = fn
}
