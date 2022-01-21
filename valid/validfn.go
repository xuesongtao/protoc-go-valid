package valid

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	strUnitStr = "length"
	numUnitStr = "size"
)

// To 验证输入的大小区间, 如果为字符串则是验证字符个数, 如果是数字的话就验证数字的大小, 注: 左右都为闭区间
func To(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	_, toVal := ParseValidNameKV(validName)
	min, max, err := parseTagTo(toVal)
	if err != nil {
		errBuf.WriteString(err.Error() + errEndFlag)
		return
	}

	isLessThan, isMoreThan, valStr, unitStr := validInputSize(min, max, tv)
	if isLessThan {
		// 生成如: "TestOrder.AppName" [xxx] len less than 2
		errBuf.WriteString("\"" + objName + "." + filedName + "\" " + "[" + valStr + "] " + unitStr + " less than or equal " + fmt.Sprintf("%d", min) + errEndFlag)
	}

	if isMoreThan {
		// 生成如: "TestOrder.AppName" [xxx] len more than 30
		errBuf.WriteString("\"" + objName + "." + filedName + "\" " + "[" + valStr + "] " + unitStr + " more than or equal " + fmt.Sprintf("%d", max) + errEndFlag)
	}
}

// Ge 大于或等于验证, 如果为字符串则是验证字符个数, 如果是数字的话就验证数字的大小
func Ge(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	_, minStr := ParseValidNameKV(validName)
	min, _ := strconv.Atoi(minStr)
	isLessThan, _, valStr, unitStr := validInputSize(min, 0, tv)
	if isLessThan {
		// 生成如: "TestOrder.AppName" [xxx] len less than 2
		errBuf.WriteString("\"" + objName + "." + filedName + "\" " + "[" + valStr + "] " + unitStr + " less than or equal " + fmt.Sprintf("%d", min) + errEndFlag)
	}
}

// Le 小于或等于验证, 如果为字符串则是验证字符个数, 如果是数字的话就验证数字的大小
func Le(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	_, maxStr := ParseValidNameKV(validName)
	max, _ := strconv.Atoi(maxStr)
	_, isMoreThan, valStr, unitStr := validInputSize(0, max, tv)
	if isMoreThan {
		// 生成如: "TestOrder.AppName" [xxx] len more than 30
		errBuf.WriteString("\"" + objName + "." + filedName + "\" " + "[" + valStr + "] " + unitStr + " more than or equal " + fmt.Sprintf("%d", max) + errEndFlag)
	}
}

// OTo 验证输入的大小区间, 如果为字符串则是验证字符个数, 如果是数字的话就验证数字的大小, 注: 左右都为开区间
func OTo(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	_, toVal := ParseValidNameKV(validName)
	min, max, err := parseTagTo(toVal)
	if err != nil {
		errBuf.WriteString(err.Error() + errEndFlag)
		return
	}

	isLessThan, isMoreThan, valStr, unitStr := validInputSize(min, max, tv, false)
	if isLessThan {
		// 生成如: "TestOrder.AppName" [xxx] len less than 2
		errBuf.WriteString("\"" + objName + "." + filedName + "\" " + "[" + valStr + "] " + unitStr + " less than " + fmt.Sprintf("%d", min) + errEndFlag)
	}

	if isMoreThan {
		// 生成如: "TestOrder.AppName" [xxx] len more than 30
		errBuf.WriteString("\"" + objName + "." + filedName + "\" " + "[" + valStr + "] " + unitStr + " more than " + fmt.Sprintf("%d", max) + errEndFlag)
	}
}

// Gt 大于验证, 如果为字符串则是验证字符个数, 如果是数字的话就验证数字的大小
func Gt(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	_, minStr := ParseValidNameKV(validName)
	min, _ := strconv.Atoi(minStr)
	isLessThan, _, valStr, unitStr := validInputSize(min, 0, tv, false)
	if isLessThan {
		// 生成如: "TestOrder.AppName" [xxx] len less than 2
		errBuf.WriteString("\"" + objName + "." + filedName + "\" " + "[" + valStr + "] " + unitStr + " less than " + fmt.Sprintf("%d", min) + errEndFlag)
	}
}

// Lt 小于验证, 如果为字符串则是验证字符个数, 如果是数字的话就验证数字的大小
func Lt(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	_, maxStr := ParseValidNameKV(validName)
	max, _ := strconv.Atoi(maxStr)
	_, isMoreThan, valStr, unitStr := validInputSize(0, max, tv, false)
	if isMoreThan {
		// 生成如: "TestOrder.AppName" [xxx] len more than 30
		errBuf.WriteString("\"" + objName + "." + filedName + "\" " + "[" + valStr + "] " + unitStr + " more than " + fmt.Sprintf("%d", max) + errEndFlag)
	}
}

// Eq 等于验证, 如果为字符串则是验证字符个数, 如果是数字的话就验证数字的大小
func Eq(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	_, eqStr := ParseValidNameKV(validName)
	eqInt, _ := strconv.Atoi(eqStr)
	isEq := true
	uintStr := numUnitStr
	valStr := ""
	switch tv.Kind() {
	case reflect.String:
		valStr = tv.String()
		uintStr = strUnitStr
		if len([]rune(valStr)) != eqInt {
			isEq = false
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := tv.Int()
		valStr = fmt.Sprintf("%d", val)
		if val != int64(eqInt) {
			isEq = false
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val := tv.Uint()
		valStr = fmt.Sprintf("%d", val)
		if val != uint64(eqInt) {
			isEq = false
		}
	case reflect.Float32, reflect.Float64:
		val := tv.Float()
		valStr = fmt.Sprintf("%v", val)
		if val != float64(eqInt) {
			isEq = false
		}
	default:
		valStr = fmt.Sprintf("%v", tv.Interface())
		isEq = false
	}

	if isEq {
		return
	}
	errBuf.WriteString("\"" + objName + "." + filedName + "\" " + "[" + valStr + "] " + uintStr + " should equal " + eqStr + errEndFlag)
}

// In 指定输入选项(精准匹配)
func In(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	_, val := ParseValidNameKV(validName)
	// 取左括号的下标
	leftBracketIndex := strings.Index(val, "(")

	// 取右括号的下标
	rightBracketIndex := strings.Index(val, ")")
	if leftBracketIndex == -1 || rightBracketIndex == -1 {
		errBuf.WriteString(inValErr.Error() + errEndFlag)
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
		tvVal = fmt.Sprintf("%v", tv.Interface())
	}

	for _, v := range strings.Split(inVals, "/") {
		if v == tvVal {
			isIn = true
			break
		}
	}

	if !isIn {
		errBuf.WriteString("\"" + objName + "." + filedName + "\" " + "[" + tvVal + "] " + "should in (" + inVals + ")" + errEndFlag)
	}
}

// Include 指定包含什么字符串(模糊匹配)
func Include(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	_, val := ParseValidNameKV(validName)
	// 取左括号的下标
	leftBracketIndex := strings.Index(val, "(")

	// 取右括号的下标
	rightBracketIndex := strings.Index(val, ")")
	if leftBracketIndex == -1 || rightBracketIndex == -1 {
		errBuf.WriteString(inValErr.Error() + errEndFlag)
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
		errBuf.WriteString(includeErr.Error() + errEndFlag)
		return
	}

	for _, v := range strings.Split(inVals, "/") {
		if strings.Contains(tvVal, v) {
			isIn = true
			break
		}
	}

	if !isIn {
		errBuf.WriteString("\"" + objName + "." + filedName + "\" " + "[" + tvVal + "] " + "should include (" + inVals + ")" + errEndFlag)
	}
}

// Phone 验证手机号
func Phone(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	if err := checkFieldIsString(validName, objName, filedName, tv); err != nil {
		errBuf.WriteString(err.Error() + errEndFlag)
		return
	}
	matched, _ := regexp.MatchString("^1[3,4,5,6,7,8,9]\\d{9}$", tv.String())
	if matched {
		return
	}
	errBuf.WriteString("\"" + objName + "." + filedName + " [" + tv.String() + "] " + "is not phone" + errEndFlag)
}

// Email 验证邮箱
func Email(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	if err := checkFieldIsString(validName, objName, filedName, tv); err != nil {
		errBuf.WriteString(err.Error() + errEndFlag)
		return
	}
	matched, _ := regexp.MatchString("^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$", tv.String())
	if matched {
		return
	}
	errBuf.WriteString("\"" + objName + "." + filedName + " [" + tv.String() + "] " + "is not email" + errEndFlag)
}

// IDCard 验证身份证
func IDCard(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	if err := checkFieldIsString(validName, objName, filedName, tv); err != nil {
		errBuf.WriteString(err.Error() + errEndFlag)
		return
	}
	matched, _ := regexp.MatchString("(^\\d{15}$)|(^\\d{18}$)|(^\\d{17}(\\d|X|x)$)", tv.String())
	if matched {
		return
	}
	errBuf.WriteString("\"" + objName + "." + filedName + " [" + tv.String() + "] " + "is not idcard" + errEndFlag)
}

// Date 验证日期
func Date(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	if err := checkFieldIsString(validName, objName, filedName, tv); err != nil {
		errBuf.WriteString(err.Error() + errEndFlag)
		return
	}
	matched, _ := regexp.MatchString("^\\d{4}-\\d{2}-\\d{2}", tv.String())
	if matched {
		return
	}
	errBuf.WriteString("\"" + objName + "." + filedName + " [" + tv.String() + "] " + "is not date, eg: 2021-11-15" + errEndFlag)
}

// Datetime 验证时间
func Datetime(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	if err := checkFieldIsString(validName, objName, filedName, tv); err != nil {
		errBuf.WriteString(err.Error() + errEndFlag)
		return
	}
	matched, _ := regexp.MatchString("^\\d{4}-\\d{1,2}-\\d{1,2} \\d{1,2}:\\d{1,2}:\\d{1,2}", tv.String())
	if matched {
		return
	}
	errBuf.WriteString("\"" + objName + "." + filedName + " [" + tv.String() + "] " + "is not datetime, eg: 2021-11-15 23:59:59" + errEndFlag)
}

// Int 验证整数
func Int(errBuf *strings.Builder, validName, objName, filedName string, tv reflect.Value) {
	matched := true
	valStr := ""
	switch tv.Kind() {
	case reflect.String:
		valStr = fmt.Sprintf("%s", tv.String())
		matched, _ = regexp.MatchString("^\\d+$", valStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valStr = fmt.Sprintf("%d", tv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		valStr = fmt.Sprintf("%d", tv.Uint())
	default:
		valStr = fmt.Sprintf("%v", tv.Interface())
		matched = false
	}

	if matched {
		return
	}
	errBuf.WriteString("\"" + objName + "." + filedName + " [" + valStr + "] " + "is not integer" + errEndFlag)
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
		valStr = fmt.Sprintf("%v", tv.Float())
	default:
		valStr = fmt.Sprintf("%v", tv.Interface())
		matched = false
	}

	if matched {
		return
	}
	errBuf.WriteString("\"" + objName + "." + filedName + " [" + valStr + "] " + "is not float" + errEndFlag)
}