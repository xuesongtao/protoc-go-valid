package valid

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// To 验证输入的大小区间, 如果为字符串则是验证字符个数, 如果是数字的话就验证数字的大小, 注: 左右都为闭区间
func To(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	_, toVal := ParseValidNameKV(validName)
	min, max, err := parseTagTo(toVal)
	if err != nil {
		errBuf.WriteString(err.Error())
		return
	}

	isLessThan, isMoreThan, valStr, unitStr := validInputSize(min, max, tv)
	if isLessThan {
		// 生成如: "TestOrder.AppName" input "xxx" len less than 2
		errBuf.WriteString(GetJoinValidErrStr(objName, filedName, valStr, unitStr, "less than or equal", strconv.Itoa(min)))
	}

	if isMoreThan {
		// 生成如: "TestOrder.AppName" input "xxx" len more than 30
		errBuf.WriteString(GetJoinValidErrStr(objName, filedName, valStr, unitStr, "more than or equal", strconv.Itoa(max)))
	}
}

// Ge 大于或等于验证, 如果为字符串则是验证字符个数, 如果是数字的话就验证数字的大小
func Ge(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	_, minStr := ParseValidNameKV(validName)
	min, _ := strconv.Atoi(minStr)
	isLessThan, _, valStr, unitStr := validInputSize(min, 0, tv)
	if isLessThan {
		// 生成如: "TestOrder.AppName" input "xxx" len less than 2
		errBuf.WriteString(GetJoinValidErrStr(objName, filedName, valStr, unitStr, "less than or equal", strconv.Itoa(min)))
	}
}

// Le 小于或等于验证, 如果为字符串则是验证字符个数, 如果是数字的话就验证数字的大小
func Le(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	_, maxStr := ParseValidNameKV(validName)
	max, _ := strconv.Atoi(maxStr)
	_, isMoreThan, valStr, unitStr := validInputSize(0, max, tv)
	if isMoreThan {
		// 生成如: "TestOrder.AppName" input "xxx" len more than 30
		errBuf.WriteString(GetJoinValidErrStr(objName, filedName, valStr, unitStr, "more than or equal", strconv.Itoa(max)))
	}
}

// OTo 验证输入的大小区间, 如果为字符串则是验证字符个数, 如果是数字的话就验证数字的大小, 注: 左右都为开区间
func OTo(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	_, toVal := ParseValidNameKV(validName)
	min, max, err := parseTagTo(toVal)
	if err != nil {
		errBuf.WriteString(err.Error())
		return
	}

	isLessThan, isMoreThan, valStr, unitStr := validInputSize(min, max, tv, false)
	if isLessThan {
		// 生成如: "TestOrder.AppName" input "xxx" len less than 2
		errBuf.WriteString(GetJoinValidErrStr(objName, filedName, valStr, unitStr, "less than", strconv.Itoa(min)))
	}

	if isMoreThan {
		// 生成如: "TestOrder.AppName" input "xxx" len more than 30
		errBuf.WriteString(GetJoinValidErrStr(objName, filedName, valStr, unitStr, "more than", strconv.Itoa(max)))
	}
}

// Gt 大于验证, 如果为字符串则是验证字符个数, 如果是数字的话就验证数字的大小
func Gt(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	_, minStr := ParseValidNameKV(validName)
	min, _ := strconv.Atoi(minStr)
	isLessThan, _, valStr, unitStr := validInputSize(min, 0, tv, false)
	if isLessThan {
		// 生成如: "TestOrder.AppName" input "xxx" len less than 2
		errBuf.WriteString(GetJoinValidErrStr(objName, filedName, valStr, unitStr, "less than", strconv.Itoa(min)))
	}
}

// Lt 小于验证, 如果为字符串则是验证字符个数, 如果是数字的话就验证数字的大小
func Lt(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	_, maxStr := ParseValidNameKV(validName)
	max, _ := strconv.Atoi(maxStr)
	_, isMoreThan, valStr, unitStr := validInputSize(0, max, tv, false)
	if isMoreThan {
		// 生成如: "TestOrder.AppName" input "xxx" len more than 30
		errBuf.WriteString(GetJoinValidErrStr(objName, filedName, valStr, unitStr, "more than", strconv.Itoa(max)))
	}
}

// Eq 等于验证, 如果为字符串则是验证字符个数, 如果是数字的话就验证数字的大小
func Eq(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	eqStr, uintStr, isEq := eq(validName, tv)
	if isEq {
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, filedName, fmt.Sprintf("%v", tv.Interface()), uintStr, "should equal", eqStr))
}

// NoEq 不等于验证, 如果为字符串则是验证字符个数, 如果是数字的话就验证数字的大小
func NoEq(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	eqStr, uintStr, isEq := eq(validName, tv)
	if !isEq {
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, filedName, fmt.Sprintf("%v", tv.Interface()), uintStr, "should no equal", eqStr))
}

// eq 相等
func eq(validName string, tv reflect.Value) (eqStr, uintStr string, isEq bool) {
	_, eqStr = ParseValidNameKV(validName)
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
func In(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	in(errBuf, validName, objName, filedName, tv, func(tvVal, v string) bool {
		return tvVal == v
	})
}

// Include 指定包含什么字符串(模糊匹配)
func Include(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	in(errBuf, validName, objName, filedName, tv, func(tvVal, v string) bool {
		return strings.Contains(tvVal, v)
	})
}

// in 是否包含
func in(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value, fn func(string, string) bool) {
	key, val := ParseValidNameKV(validName)
	// 取左括号的下标
	leftBracketIndex := strings.Index(val, "(")

	useErrMsg := inValErr.Error()
	if key == "include" {
		useErrMsg = includeErr.Error()
	}

	// 取右括号的下标
	rightBracketIndex := strings.Index(val, ")")
	if leftBracketIndex == -1 || rightBracketIndex == -1 {
		errBuf.WriteString(useErrMsg)
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
		errBuf.WriteString(GetJoinValidErrStr(objName, filedName, tvVal, "should "+key+" ("+inVals+")"))
	}
}

// Phone 验证手机号
func Phone(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	if err := checkFieldIsString(objName, filedName, tv); err != nil {
		errBuf.WriteString(err.Error())
		return
	}
	matched, _ := regexp.MatchString("^1[3,4,5,6,7,8,9]\\d{9}$", tv.String())
	if matched {
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, filedName, tv.String(), "is not phone"))
}

// Email 验证邮箱
func Email(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	if err := checkFieldIsString(objName, filedName, tv); err != nil {
		errBuf.WriteString(err.Error())
		return
	}
	matched, _ := regexp.MatchString("^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$", tv.String())
	if matched {
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, filedName, tv.String(), "is not email"))
}

// IDCard 验证身份证
func IDCard(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	if err := checkFieldIsString(objName, filedName, tv); err != nil {
		errBuf.WriteString(err.Error())
		return
	}
	matched, _ := regexp.MatchString("(^\\d{15}$)|(^\\d{18}$)|(^\\d{17}(\\d|X|x)$)", tv.String())
	if matched {
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, filedName, tv.String(), "is not idcard"))
}

// Date 验证日期
func Date(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	if err := checkFieldIsString(objName, filedName, tv); err != nil {
		errBuf.WriteString(err.Error())
		return
	}
	matched, _ := regexp.MatchString("^\\d{4}-\\d{2}-\\d{2}", tv.String())
	if matched {
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, filedName, tv.String(), "is not date, eg: 1996-09-28"))
}

// Datetime 验证时间
func Datetime(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	if err := checkFieldIsString(objName, filedName, tv); err != nil {
		errBuf.WriteString(err.Error())
		return
	}
	matched, _ := regexp.MatchString("^\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}", tv.String())
	if matched {
		return
	}
	errBuf.WriteString(GetJoinValidErrStr(objName, filedName, tv.String(), "is not datetime, eg: 1996-09-28 23:00:00"))
}

// Int 验证整数
func Int(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
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
	errBuf.WriteString(GetJoinValidErrStr(objName, filedName, valStr, "is not integer"))
}

// Float 验证浮动数
func Float(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
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
	errBuf.WriteString(GetJoinValidErrStr(objName, filedName, valStr, "is not float"))
}
