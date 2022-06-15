package valid

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// To 验证输入的大小区间, 注: 左右都为闭区间
// 1. 如果为字符串则是验证字符个数
// 2. 如果是数字的话就验证数字的大小
// 3. 如果是切片的话就验证的长度
func To(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	_, toVal, cusMsg := ParseValidNameKV(validName)
	min, max, err := parseTagTo(toVal, true)
	if err != nil {
		errBuf.WriteString(err.Error() + ErrEndFlag)
		return
	}

	isLessThan, isMoreThan, valStr, unitStr := validInputSize(min, max, tv)
	if isLessThan {
		if cusMsg != "" {
			errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, cusMsg))
			return
		}
		// 生成如: "TestOrder.AppName" input "xxx" len less than 2
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, unitStr, "less than", strconv.Itoa(min)))
	}

	if isMoreThan {
		if cusMsg != "" {
			errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, cusMsg))
			return
		}
		// 生成如: "TestOrder.AppName" input "xxx" len more than 30
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, unitStr, "more than", strconv.Itoa(max)))
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
		// 生成如: "TestOrder.AppName" input "xxx" len less than 2
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, unitStr, "less than", strconv.Itoa(min)))
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
		// 生成如: "TestOrder.AppName" input "xxx" len more than 30
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, unitStr, "more than", strconv.Itoa(max)))
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
		errBuf.WriteString(err.Error() + ErrEndFlag)
		return
	}

	isLessThan, isMoreThan, valStr, unitStr := validInputSize(min, max, tv, false)
	if isLessThan {
		if cusMsg != "" {
			errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, cusMsg))
			return
		}
		// 生成如: "TestOrder.AppName" input "xxx" len less than 2
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, unitStr, "less than or equal", strconv.Itoa(min)))
	}

	if isMoreThan {
		if cusMsg != "" {
			errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, cusMsg))
			return
		}
		// 生成如: "TestOrder.AppName" input "xxx" len more than 30
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, unitStr, "more than or equal", strconv.Itoa(max)))
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
		// 生成如: "TestOrder.AppName" input "xxx" len less than 2
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, unitStr, "less than or equal", strconv.Itoa(min)))
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
		// 生成如: "TestOrder.AppName" input "xxx" len more than 30
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, unitStr, "more than or equal", strconv.Itoa(max)))
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
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, fmt.Sprintf("%v", tv.Interface()), uintStr, cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, fmt.Sprintf("%v", tv.Interface()), uintStr, "should equal", eqStr))
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
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, fmt.Sprintf("%v", tv.Interface()), cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, fmt.Sprintf("%v", tv.Interface()), uintStr, "should no equal", eqStr))
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

	useErrMsg := inValErr.Error()
	if key == "include" {
		useErrMsg = includeErr.Error()
	}

	// 取右括号的下标
	rightBracketIndex := strings.Index(val, ")")
	if leftBracketIndex == -1 || rightBracketIndex == -1 {
		errBuf.WriteString(useErrMsg + ErrEndFlag)
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
			errBuf.WriteString(useErrMsg)
			return
		}
		tvVal = fmt.Sprintf("%v", tv.Interface())
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
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tvVal, "should "+key+" ("+inVals+")"))
	}
}

// Phone 验证手机号
func Phone(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	if err := CheckFieldIsStr(objName, fieldName, tv); err != nil {
		errBuf.WriteString(err.Error())
		return
	}
	matched, _ := regexp.MatchString("^1[3,4,5,6,7,8,9]\\d{9}$", tv.String())
	if matched {
		return
	}

	_, _, cusMsg := ParseValidNameKV(validName)
	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), "is not phone"))
}

// Email 验证邮箱
func Email(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	if err := CheckFieldIsStr(objName, fieldName, tv); err != nil {
		errBuf.WriteString(err.Error())
		return
	}
	matched, _ := regexp.MatchString("^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$", tv.String())
	if matched {
		return
	}

	_, _, cusMsg := ParseValidNameKV(validName)
	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), "is not email"))
}

// IDCard 验证身份证
func IDCard(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	if err := CheckFieldIsStr(objName, fieldName, tv); err != nil {
		errBuf.WriteString(err.Error())
		return
	}
	matched, _ := regexp.MatchString("(^\\d{15}$)|(^\\d{18}$)|(^\\d{17}(\\d|X|x)$)", tv.String())
	if matched {
		return
	}

	_, _, cusMsg := ParseValidNameKV(validName)
	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), "is not idcard"))
}

// Year 验证年
func Year(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	if err := CheckFieldIsStr(objName, fieldName, tv); err != nil {
		errBuf.WriteString(err.Error())
		return
	}
	matched, _ := regexp.MatchString("^\\d{4}$", tv.String())
	if matched {
		return
	}

	_, _, cusMsg := ParseValidNameKV(validName)
	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), "is not year, eg: 1996"))
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
		defaultDateSplit = val
	}
	matched, _ := regexp.MatchString("^\\d{4}"+defaultDateSplit+"\\d{2}$", tv.String())
	if matched {
		return
	}

	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), "is not year2month, eg: 1996"+defaultDateSplit+"09"))
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
		defaultDateSplit = val
	}
	matched, _ := regexp.MatchString("^\\d{4}"+defaultDateSplit+"\\d{2}"+defaultDateSplit+"\\d{2}$", tv.String())
	if matched {
		return
	}

	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), "is not date, eg: 1996"+defaultDateSplit+"09"+defaultDateSplit+"28"))
}

// Datetime 验证时间
// 默认匹配 xxxx-xx-xx xx:xx:xx, 可以指定分割符
func Datetime(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	if err := CheckFieldIsStr(objName, fieldName, tv); err != nil {
		errBuf.WriteString(err.Error())
		return
	}
	defaultDateSplit := "-" // 默认时间拼接符号
	_, val, cusMsg := ParseValidNameKV(validName)
	if val != "" {
		defaultDateSplit = val
	}
	matched, _ := regexp.MatchString("^\\d{4}"+defaultDateSplit+"\\d{2}"+defaultDateSplit+"\\d{2} \\d{2}:\\d{2}:\\d{2}$", tv.String())
	if matched {
		return
	}

	if cusMsg != "" {
		errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), cusMsg))
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), "is not datetime, eg: 1996"+defaultDateSplit+"09"+defaultDateSplit+"28 23:00:00"))
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
		errBuf.WriteString(reErr.Error() + ErrEndFlag)
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
			errBuf.WriteString(reErr.Error() + ErrEndFlag)
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
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), "regex match is failed, regex: "+pattern))
}

// Int 验证整数
func Int(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	matched := true
	valStr := ""
	switch tv.Kind() {
	case reflect.String:
		valStr = tv.String()
		matched, _ = regexp.MatchString("^\\d+$", valStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	default:
		valStr = fmt.Sprintf("%v", tv.Interface())
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
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, "is not integer"))
}

// Float 验证浮动数
func Float(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
	matched := true
	valStr := ""
	switch tv.Kind() {
	case reflect.String:
		valStr = tv.String()
		matched, _ = regexp.MatchString("^\\d+.\\d+$", valStr)
	case reflect.Float32, reflect.Float64:
	default:
		valStr = fmt.Sprintf("%v", tv.Interface())
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
	errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, valStr, "is not float"))
}
