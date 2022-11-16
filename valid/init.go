package valid

import (
	"errors"
	"reflect"
	"regexp"
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

const (
	// 时间格式
	YearFmt int8 = 1 << iota
	MonthFmt
	DayFmt
	HourFmt
	MinFmt
	SecFmt

	DateFmt     = YearFmt | MonthFmt | DayFmt
	DateTimeFmt = DateFmt | HourFmt | MinFmt | SecFmt
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
	VIp         = "ip"         // ip
	VIpv4       = "ipv4"       // ipv4
	VIpv6       = "ipv6"       // ipv6
	VUnique     = "unique"     // 唯一验证
	VJson       = "json"       // json 格式验证
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
	VIp:         Ip,
	VIpv4:       Ipv4,
	VIpv6:       Ipv6,
	VUnique:     Unique,
	VJson:       Json,
}

// 对象
var (
	syncValidStructPool = sync.Pool{New: func() interface{} { return new(VStruct) }}
	// 如果需要释放内存可以通过调用 SetStructTypeCache, 如: SetStructTypeCache(NewLRU(2 << 8))
	// 考虑到性能, 用 sync.Map 缓存(缺点: 内存释放不到)
	// cacheStructType  CacheEr = new(sync.Map)
	cacheStructType  CacheEr = NewLRU(1 << 8)
	syncValidVarPool         = sync.Pool{New: func() interface{} { return new(VVar) }}
	syncBufPool              = sync.Pool{New: func() interface{} { return new(strings.Builder) }}
	timeReflectType          = reflect.TypeOf(time.Time{})
	once             sync.Once
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

// 正则
var (
	IncludeZhRe = regexp.MustCompile("[\u4e00-\u9fa5]")         // 中文
	PhoneRe     = regexp.MustCompile(`^1[3,4,5,6,7,8,9]\d{9}$`) // 手机号
	Ipv4Re      = regexp.MustCompile(`^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$`)
	EmailRe     = regexp.MustCompile(`^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`)
	IdCardRe    = regexp.MustCompile(`(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)`)
	IntRe       = regexp.MustCompile(`^\d+$`)
	FloatRe     = regexp.MustCompile(`^\d+.\d+$`)

	// Deprecated
	YearRe = regexp.MustCompile(`^\d{4}$`)
	// Deprecated
	Year2MonthRe = regexp.MustCompile(`^\d{4}-\d{2}$`)
	// Deprecated
	Year2MonthRe2 = regexp.MustCompile(`^\d{4}/\d{2}$`)
	// Deprecated
	DateRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	// Deprecated
	DateRe2 = regexp.MustCompile(`^\d{4}/\d{2}/\d{2}$`)
	// Deprecated
	DatetimeRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`)
	// Deprecated
	DatetimeRe2 = regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}$`)
)

// newStrBuf
func newStrBuf(size ...int) *strings.Builder {
	obj := syncBufPool.Get().(*strings.Builder)
	if len(size) > 0 {
		obj.Grow(size[0])
	}
	return obj
}

// putStrBuf
func putStrBuf(buf *strings.Builder) {
	if buf.Len() > 0 {
		buf.Reset()
	}
	syncBufPool.Put(buf)
}

// CommonValidFn 通用验证函数, 主要用于回调
// 注: 在写 errBuf 的时候建议用 GetJoinValidErrStr 包裹下, 这样产生的结果易读.
//    否则需要再 errBuf.WriteString 最后要加上 ErrEndFlag 分割, 工具是通过 ErrEndFlag 进行分句
type CommonValidFn func(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value)

// ValidName2ValidFnMap 自定义验证名对应自定义验证函数
type ValidName2ValidFnMap map[string]CommonValidFn

// SetCustomerValidFn 自定义验证函数
// 用于全局添加验证方法, 如果不想定义全局, 可根据验证对象分别调用 SetValidFn, 如: *VStruct.SetValidFn
func SetCustomerValidFn(validName string, fn CommonValidFn) {
	validName2FuncMap[validName] = fn
}

// SetStructTypeCache 设置 structType 缓存类型
func SetStructTypeCache(cacheEr CacheEr) {
	once.Do(func() {
		cacheStructType = cacheEr
	})
}

// GetOnlyExplainErr 获取所有的说明错误(不包含错误的字段信息)
// 使用场景: 在自定义错误信息时, 返回给非开发人员看的结果
func GetOnlyExplainErr(errMsg string) string {
	if errMsg == "" {
		return ""
	}
	buf := newStrBuf(1 << 8)
	defer putStrBuf(buf)
	zhLen := len(ExplainZh)
	enLen := len(ExplainEn)
	endLen := len(ErrEndFlag)
	splitLen := zhLen
	nullLen := 1 // err msg [说明: xxx] 里包含一个空需要处理
	for {
		s := strings.Index(errMsg, ExplainZh)
		e := strings.Index(errMsg, ErrEndFlag) // 未发现的话, 为最后一句错误
		if s == -1 || (e != -1 && s > e) {     // 说明为英文
			s = strings.Index(errMsg, ExplainEn)
			splitLen = enLen
		}
		if s == -1 { // 异常
			break
		}
		if e == -1 {
			buf.WriteString(errMsg[s+splitLen+nullLen:])
			break
		}
		buf.WriteString(errMsg[s+splitLen+nullLen : e])
		buf.WriteString(ErrEndFlag)
		errMsg = errMsg[e+endLen:]
	}
	return buf.String()
}

// GetTimeFmt 获取时间格式化
// splits 为分隔符
// splits[0] 为 [年月日] 的分割符, 默认为 "-"
// splits[1] 为 [年月日] 和 [时分秒] 的分割符, 默认为 " "
// splits[2] 为 [时分秒] 的分割符, 默认为 ":"
func GetTimeFmt(fmtType int8, splits ...string) (res string) {
	defaultDateSplit := "-"
	defaultDateTimeSplit := " "
	defaultTimeSplit := ":"
	l := len(splits)
	switch l {
	case 1:
		defaultDateSplit = splits[0]
	case 2:
		defaultDateSplit = splits[0]
		defaultDateTimeSplit = splits[1]
	case 3:
		defaultDateSplit = splits[0]
		defaultDateTimeSplit = splits[1]
		defaultTimeSplit = splits[2]
	}

	// 年月日
	if fmtType&YearFmt > 0 {
		res += "2006" + defaultDateSplit
	}
	if fmtType&MonthFmt > 0 {
		res += "01" + defaultDateSplit
	}
	if fmtType&DayFmt > 0 {
		res += "02" + defaultDateSplit
	}
	res = strings.TrimSuffix(res, defaultDateSplit)

	// 时分秒
	if res != "" {
		res += defaultDateTimeSplit
	}
	if fmtType&HourFmt > 0 {
		res += "15" + defaultTimeSplit
	}
	if fmtType&MinFmt > 0 {
		res += "04" + defaultTimeSplit
	}
	if fmtType&SecFmt > 0 {
		res += "05" + defaultTimeSplit
	}
	res = strings.TrimSuffix(res, defaultTimeSplit)
	return
}
