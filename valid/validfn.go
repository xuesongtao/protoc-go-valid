package valid

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gitee.com/xuesongtao/protoc-go-valid/valid/internal"
)

// 注: 拓展规则
// 	1. 如果为 key value 的需要通过 "=" 连接
//	2. 如果 value 中有包裹情况, 只能通过 "''" 进行包裹; 已上线其他字符包裹了一直保留兼容, 如: in=(xxx); include=(xxx) 等

// To 验证输入的大小区间, 注: 左右都为闭区间
// 1. 如果为字符串则是验证字符个数
// 2. 如果是数字的话就验证数字的大小
// 3. 如果是切片的话就验证的长度
func To(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	_, toVal, cusMsg := ParseValidNameKV(validName)
	min, max, err := parseTagTo(toVal, true)
	if err != nil {
		errBuf.WriteString(GetJoinFieldErr(objName, fieldName, err))
		return
	}

	isLessThan, isMoreThan, valStr, unitStr := validInputSize(min, max, tv)
	if isLessThan {
		if cusMsg != "" {
			errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, cusMsg))
			return
		}
		// 生成如: "TestOrder.AppName" input "xxx", Explain: it is less than 2 length
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, ExplainEn, "it is less than", ToStr(min), unitStr))
	}

	if isMoreThan {
		if cusMsg != "" {
			errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, cusMsg))
			return
		}
		// 生成如: "TestOrder.AppName" input "xxx", Explain: it is more than 30 length
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, ExplainEn, "it is more than", ToStr(max), unitStr))
	}
}

// Ge 大于或等于验证
// 1. 如果为字符串则是验证字符个数
// 2. 如果是数字的话就验证数字的大小
// 3. 如果是切片的话就验证的长度
func Ge(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	_, minStr, cusMsg := ParseValidNameKV(validName)
	min, _ := strconv.Atoi(minStr)
	isLessThan, _, valStr, unitStr := validInputSize(min, 0, tv)
	if isLessThan {
		if cusMsg != "" {
			errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, cusMsg))
			return
		}
		// 生成如: "TestOrder.AppName" input "xxx", Explain: it is less than 2 length
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, ExplainEn, "it is less than", ToStr(min), unitStr))
	}
}

// Le 小于或等于验证
// 1. 如果为字符串则是验证字符个数
// 2. 如果是数字的话就验证数字的大小
// 3. 如果是切片的话就验证的长度
func Le(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	_, maxStr, cusMsg := ParseValidNameKV(validName)
	max, _ := strconv.Atoi(maxStr)
	_, isMoreThan, valStr, unitStr := validInputSize(0, max, tv)
	if isMoreThan {
		if cusMsg != "" {
			errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, cusMsg))
			return
		}
		// 生成如: "TestOrder.AppName" input "xxx", Explain: it is len more than 30 length
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, ExplainEn, "it is more than", ToStr(max), unitStr))
	}
}

// OTo 验证输入的大小区间, 注: 左右都为开区间
// 1. 如果为字符串则是验证字符个数
// 2. 如果是数字的话就验证数字的大小
// 3. 如果是切片的话就验证的长度
func OTo(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	_, toVal, cusMsg := ParseValidNameKV(validName)
	min, max, err := parseTagTo(toVal, false)
	if err != nil {
		errBuf.WriteString(GetJoinFieldErr(objName, fieldName, err))
		return
	}

	isLessThan, isMoreThan, valStr, unitStr := validInputSize(min, max, tv, false)
	if isLessThan {
		if cusMsg != "" {
			errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, cusMsg))
			return
		}
		// 生成如: "TestOrder.AppName" input "xxx", Explain: it is less than 2 length
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, ExplainEn, "it is less than or equal", ToStr(min), unitStr))
	}

	if isMoreThan {
		if cusMsg != "" {
			errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, cusMsg))
			return
		}
		// 生成如: "TestOrder.AppName" input "xxx", Explain: it is more than 30 length
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, ExplainEn, "it is more than or equal", ToStr(max), unitStr))
	}
}

// Gt 大于验证
// 1. 如果为字符串则是验证字符个数
// 2. 如果是数字的话就验证数字的大小
// 3. 如果是切片的话就验证的长度
func Gt(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	_, minStr, cusMsg := ParseValidNameKV(validName)
	min, _ := strconv.Atoi(minStr)
	isLessThan, _, valStr, unitStr := validInputSize(min, 0, tv, false)

	if isLessThan {
		if cusMsg != "" {
			errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, cusMsg))
			return
		}
		// 生成如: "TestOrder.AppName" input "xxx", Explain: it is less than or equal 2 length
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, ExplainEn, "it is less than or equal", ToStr(min), unitStr))
	}
}

// Lt 小于验证, 如果为字符串则是验证字符个数, 如果是数字的话就验证数字的大小
func Lt(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	_, maxStr, cusMsg := ParseValidNameKV(validName)
	max, _ := strconv.Atoi(maxStr)
	_, isMoreThan, valStr, unitStr := validInputSize(0, max, tv, false)
	if isMoreThan {
		if cusMsg != "" {
			errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, cusMsg))
			return
		}
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, ExplainEn, "it is more than or equal", ToStr(max), unitStr))
	}
}

// Eq 等于验证
// 1. 如果为字符串则是验证字符个数
// 2. 如果是数字的话就验证数字的大小
// 3. 如果是切片的话就验证的长度
func Eq(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	eqStr, uintStr, cusMsg, isEq := eq(validName, tv)
	if isEq {
		return
	}

	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, ToStr(tv.Interface()), cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, ToStr(tv.Interface()), ExplainEn, "it should equal", eqStr, uintStr))
}

// NoEq 不等于验证
// 1. 如果为字符串则是验证字符个数
// 2. 如果是数字的话就验证数字的大小
// 3. 如果是切片的话就验证的长度
func NoEq(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	eqStr, uintStr, cusMsg, isEq := eq(validName, tv)
	if !isEq {
		return
	}

	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, ToStr(tv.Interface()), cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, ToStr(tv.Interface()), ExplainEn, "it is not equal", eqStr, uintStr))
}

// eq 相等
func eq(validName string, tv reflect.Value) (eqStr, uintStr, cusMsg string, isEq bool) {
	_, eqStr, cusMsg = ParseValidNameKV(validName)
	eqInt, _ := strconv.Atoi(eqStr)
	isEq = true
	uintStr = numUnitStr
	switch tv.Kind() {
	case reflect.String:
		uintStr = strUnitStr
		if len([]rune(tv.String())) != eqInt {
			isEq = false
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if tv.Int() != int64(eqInt) {
			isEq = false
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if tv.Uint() != uint64(eqInt) {
			isEq = false
		}
	case reflect.Float32, reflect.Float64:
		if tv.Float() != float64(eqInt) {
			isEq = false
		}
	default:
		isEq = false
	}
	return
}

// In 指定输入选项(精准匹配)
func In(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	in(errBuf, validName, objName, fieldName, tv, func(tvVal, v string) bool {
		return tvVal == v
	})
}

// Include 指定包含什么字符串(模糊匹配)
func Include(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	in(errBuf, validName, objName, fieldName, tv, func(tvVal, v string) bool {
		return strings.Contains(tvVal, v)
	})
}

// in 是否包含
func in(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value, fn func(string, string) bool) {
	key, val, cusMsg := ParseValidNameKV(validName)
	// 取左括号的下标
	leftBracketIndex := strings.Index(val, "(")

	useErrMsg := inValErr
	if key == "include" {
		useErrMsg = includeErr
	}

	// 取右括号的下标
	rightBracketIndex := strings.LastIndex(val, ")")
	if leftBracketIndex == -1 || rightBracketIndex == -1 {
		errBuf.WriteString(GetJoinFieldErr(objName, fieldName, useErrMsg))
		return
	}

	var (
		isIn   bool   // 默认不在输入范围选项中
		tvVal  string // 获取输入内容的值, 用于判断
		inVals = val[leftBracketIndex+1 : rightBracketIndex]
	)

	switch tv.Kind() {
	case reflect.String:
		tvVal = tv.String()
	default:
		// include 必须为字符串才验证, 其他就不处理
		if key == "include" {
			errBuf.WriteString(GetJoinFieldErr(objName, fieldName, useErrMsg))
			return
		}
		tvVal = ToStr(tv.Interface())
	}

	for _, v := range strings.Split(inVals, "/") {
		if fn(tvVal, v) {
			isIn = true
			break
		}
	}

	if !isIn {
		if cusMsg != "" {
			errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tvVal, cusMsg))
			return
		}
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tvVal, ExplainEn, "it should "+key+" ("+inVals+")"))
	}
}

// Phone 验证手机号
func Phone(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	if err := CheckFieldIsStr(objName, fieldName, tv); err != nil {
		errBuf.WriteString(err.Error())
		return
	}
	matched := PhoneRe.MatchString(tv.String())
	if matched {
		return
	}

	_, _, cusMsg := ParseValidNameKV(validName)
	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), ExplainEn, "it is not phone"))
}

// Ip ip 验证
func Ip(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	if err := CheckFieldIsStr(objName, fieldName, tv); err != nil {
		errBuf.WriteString(err.Error())
		return
	}
	ip := net.ParseIP(tv.String())
	if ip != nil {
		return
	}
	_, _, cusMsg := ParseValidNameKV(validName)
	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), ExplainEn, "it is not ip"))
}

// Ipv4 ipv4 验证
func Ipv4(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	if err := CheckFieldIsStr(objName, fieldName, tv); err != nil {
		errBuf.WriteString(err.Error())
		return
	}
	ip := net.ParseIP(tv.String())
	if ip != nil && ip.To4() != nil {
		return
	}
	_, _, cusMsg := ParseValidNameKV(validName)
	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), ExplainEn, "it is not ipv4"))
}

// Ipv6 ipv6 验证
func Ipv6(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	if err := CheckFieldIsStr(objName, fieldName, tv); err != nil {
		errBuf.WriteString(err.Error())
		return
	}
	ip := net.ParseIP(tv.String())
	if ip != nil && ip.To4() == nil {
		return
	}
	_, _, cusMsg := ParseValidNameKV(validName)
	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), ExplainEn, "it is not ipv6"))
}

// Email 验证邮箱
func Email(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	if err := CheckFieldIsStr(objName, fieldName, tv); err != nil {
		errBuf.WriteString(err.Error())
		return
	}
	matched := EmailRe.MatchString(tv.String())
	if matched {
		return
	}

	_, _, cusMsg := ParseValidNameKV(validName)
	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), ExplainEn, "it is not email"))
}

// IDCard 验证身份证
func IDCard(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	if err := CheckFieldIsStr(objName, fieldName, tv); err != nil {
		errBuf.WriteString(err.Error())
		return
	}
	matched := IdCardRe.MatchString(tv.String())
	if matched {
		return
	}

	_, _, cusMsg := ParseValidNameKV(validName)
	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), ExplainEn, "it is not idcard"))
}

// Year 验证年
func Year(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	if err := CheckFieldIsStr(objName, fieldName, tv); err != nil {
		errBuf.WriteString(err.Error())
		return
	}
	_, err := time.Parse(GetTimeFmt(YearFmt), tv.String())
	if err == nil {
		return
	}

	_, _, cusMsg := ParseValidNameKV(validName)
	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), ExplainEn, "it is not year, eg: 1996"))
}

// Year2Month 验证年月
// 默认匹配 xxxx-xx, 可以指定分割符
func Year2Month(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	if err := CheckFieldIsStr(objName, fieldName, tv); err != nil {
		errBuf.WriteString(err.Error())
		return
	}
	defaultDateSplit := "-" // 默认时间拼接符号
	_, val, cusMsg := ParseValidNameKV(validName)
	if val != "" {
		defaultDateSplit = strings.Trim(val, "'")
	}
	_, err := time.Parse(GetTimeFmt(YearFmt|MonthFmt, defaultDateSplit), tv.String())
	if err == nil {
		return
	}

	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), ExplainEn, "it is not year2month, eg: 1996"+defaultDateSplit+"09"))
}

// Date 验证日期
// 默认匹配 xxxx-xx-xx, 可以指定分割符
func Date(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	if err := CheckFieldIsStr(objName, fieldName, tv); err != nil {
		errBuf.WriteString(err.Error())
		return
	}
	defaultDateSplit := "-" // 默认时间拼接符号
	_, val, cusMsg := ParseValidNameKV(validName)
	if val != "" {
		defaultDateSplit = strings.Trim(val, "'")
	}
	_, err := time.Parse(GetTimeFmt(DateFmt, defaultDateSplit), tv.String())
	if err == nil {
		return
	}

	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), ExplainEn, "it is not date, eg: 1996"+defaultDateSplit+"09"+defaultDateSplit+"28"))
}

// Datetime 验证时间
// 默认匹配 xxxx-xx-xx xx:xx:xx, 可以指定分割符
func Datetime(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	if err := CheckFieldIsStr(objName, fieldName, tv); err != nil {
		errBuf.WriteString(err.Error())
		return
	}
	_, val, cusMsg := ParseValidNameKV(validName)
	defaultSplit := []string{"-", " ", ":"}
	if val != "" {
		for i, split := range strings.Split(strings.Trim(val, "'"), ",") {
			defaultSplit[i] = split
		}
	}
	_, err := time.Parse(GetTimeFmt(DateTimeFmt, defaultSplit...), tv.String())
	if err == nil {
		return
	}

	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(
		objName, fieldName, tv.String(), ExplainEn,
		fmt.Sprintf("it is not datetime, eg: 1996%s09%s28%s23%s00%s00", defaultSplit[0], defaultSplit[0], defaultSplit[1], defaultSplit[2], defaultSplit[2]),
	))
}

// Re 正则表达式
// 使用格式如: re='\\d+'|匹配错误
func Re(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	if err := CheckFieldIsStr(objName, fieldName, tv); err != nil {
		errBuf.WriteString(err.Error())
		return
	}

	// 解析正则, 使用格式: re='\\d+'|匹配错误
	splitIndex := strings.Index(validName, "'")
	if splitIndex == -1 {
		errBuf.WriteString(GetJoinFieldErr(objName, fieldName, reErr))
		return
	}

	l := len(validName)
	b := make([]byte, 0, l)
	i := splitIndex + 1
	for ; i < l; i++ {
		v := validName[i]
		b = append(b, v)
		// 寻找结束 "'", 同时需要跳过里面有转义的单引号("\'")
		next := i + 1
		if next > l-1 {
			errBuf.WriteString(GetJoinFieldErr(objName, fieldName, reErr))
			return
		}

		if v != '\\' && validName[next] == '\'' {
			break
		}
	}

	pattern := string(b)
	newValidName := validName[:splitIndex] + validName[i+1:] // 重新解析下自定义消息, 这里已经排除正则部分, 处理结果为: re='|xxxx
	// fmt.Printf("pattern: %s, newValidName: %s\n", pattern, newValidName)
	_, _, cusMsg := ParseValidNameKV(newValidName)
	matched, _ := regexp.MatchString(pattern, tv.String())
	if matched {
		return
	}

	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), ExplainEn, "regex match is failed, pattern: "+pattern))
}

// Int 验证整数
func Int(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	matched := true
	valStr := ""
	switch kind := tv.Kind(); kind {
	case reflect.String:
		valStr = tv.String()
		matched = IntRe.MatchString(valStr)
	default:
		if !ReflectKindIsNum(kind) {
			valStr = ToStr(tv.Interface())
			matched = false
		}
	}

	if matched {
		return
	}

	_, _, cusMsg := ParseValidNameKV(validName)
	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, ExplainEn, "it is not integer"))
}

// Ints 验证是否为多个数字
// 1. 如果输入为 string, 默认按逗号拼接进行处理
// 2. 如果为 slice/array, 会将每个值进行匹配判断
func Ints(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	is := true
	valStr := ""
	errSuffix := ""
	_, split, cusMsg := ParseValidNameKV(validName)
	if split == "" {
		split = ","
	}
	switch kind := tv.Kind(); kind {
	case reflect.String:
		valStr = tv.String()
		for _, v := range strings.Split(valStr, split) {
			is = IntRe.MatchString(v)
			if !is {
				break
			}
		}
		errSuffix = "it is not separated by \"" + split + "\"" + " num"
	case reflect.Array, reflect.Slice:
		var tmpIs bool
		l := tv.Len()
		valStr = "["
		for i := 0; i < l; i++ {
			v := ToStr(tv.Index(i).Interface())
			tmpIs = IntRe.MatchString(v)
			if !tmpIs {
				is = tmpIs
			}
			if valStr == "[" {
				valStr += v
			} else {
				valStr += ", " + v
			}
		}
		valStr += "]"
		errSuffix = "slice/array element is not all num"
	default:
		if ReflectKindIsNum(kind) {
			return
		}
		errBuf.WriteString(GetJoinFieldErr(objName, fieldName, intsErr))
		return
	}

	if is {
		return
	}

	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, ExplainEn, errSuffix))
}

// Float 验证浮动数
func Float(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	matched := true
	valStr := ""
	switch tv.Kind() {
	case reflect.String:
		valStr = tv.String()
		matched = FloatRe.MatchString(valStr)
	case reflect.Float32, reflect.Float64:
	default:
		valStr = ToStr(tv.Interface())
		matched = false
	}

	if matched {
		return
	}

	_, _, cusMsg := ParseValidNameKV(validName)
	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, ExplainEn, "it is not float"))
}

// Unique 对集合字段进行唯一验证
// 1. 对以逗号隔开的字符串进行唯一验证
// 2. 对切片/数组元素[int 系列, float系列, bool系列, string系列]进行唯一验证
func Unique(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	var (
		ok        bool
		inVal     string
		uniqueMap map[string]struct{}
	)
	switch tv.Kind() {
	case reflect.String:
		inVal = tv.String()
		inValArr := strings.Split(inVal, ",")
		uniqueMap = make(map[string]struct{}, len(inValArr))
		for _, v := range inValArr {
			uniqueMap[v] = struct{}{}
		}
		ok = len(inValArr) == len(uniqueMap)
	case reflect.Slice, reflect.Array:
		inVal = "["
		l := tv.Len()
		uniqueMap = make(map[string]struct{}, l)
		for i := 0; i < l; i++ {
			val := ToStr(tv.Index(i).Interface())
			uniqueMap[val] = struct{}{}

			if inVal == "[" {
				inVal += val
			} else {
				inVal += "," + val
			}
		}
		inVal += "]"
		ok = l == len(uniqueMap)
	default:
		errBuf.WriteString(GetJoinFieldErr(objName, fieldName, uniqueErr))
		return
	}

	if ok {
		return
	}

	_, _, cusMsg := ParseValidNameKV(validName)
	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, inVal, cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, inVal, ExplainEn, "they're not unique"))
}

// Json 验证是否为 json
func Json(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	if err := CheckFieldIsStr(objName, fieldName, tv); err != nil {
		errBuf.WriteString(err.Error())
		return
	}
	valStr := tv.String()
	if json.Valid(internal.UnsafeStr2Bytes(valStr)) {
		return
	}
	if len(valStr) > 2<<7 {
		valStr = "more than 256 byte(it is ignore)"
	}
	valStr = StrEscape(valStr)
	_, _, cusMsg := ParseValidNameKV(validName)
	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, ExplainEn, "it is not json"))
}

// Prefix 验证字符串包含前缀
func Prefix(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	if err := CheckFieldIsStr(objName, fieldName, tv); err != nil {
		errBuf.WriteString(err.Error())
		return
	}
	_, prefix, cusMsg := ParseValidNameKV(validName)
	if strings.HasPrefix(tv.String(), prefix) {
		return
	}
	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), ExplainEn, "prefix is not ok"))
}

// Suffix 验证字符串包含后缀
func Suffix(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	if err := CheckFieldIsStr(objName, fieldName, tv); err != nil {
		errBuf.WriteString(err.Error())
		return
	}
	_, suffix, cusMsg := ParseValidNameKV(validName)
	if strings.HasSuffix(tv.String(), suffix) {
		return
	}
	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), ExplainEn, "suffix is not ok"))
}
