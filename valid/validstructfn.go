package valid

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// ParseValidNameKV 解析 validName 中的 key 和 value,
// 如: "required", key 为 "required", value 为 ""
// 如: "to=1~2", key 为 "to", value 为 "1~2"
func ParseValidNameKV(validName string) (key, value string) {
	// 因为 validName 中的 k, v 通过 = 连接
	splitIndex := strings.Index(validName, "=")

	// 如果没有则代表 validName 不为 k=v 类型, 只有一个字段如: required
	if splitIndex == -1 {
		key = validName
		return
	}

	key = validName[:splitIndex]
	value = validName[splitIndex+1:]
	return
}

// validInputSize 验证输入的大小
func validInputSize(min, max int, tv reflect.Value, isHasEqual ...bool) (isLessThan, isMoreThan bool, unitStr string) {
	hasEqual := true // 标记对结果默认包含闭区间
	if len(isHasEqual) > 0 {
		hasEqual = isHasEqual[0]
	}
	unitStr = "size" // 大小的单位
	// fmt.Printf("min: %v, max: %v\n", min, max)
	switch tv.Kind() {
	case reflect.String:
		unitStr = "length"
		inLen := len([]rune(tv.String()))
		if hasEqual {
			if inLen < min {
				isLessThan = true
			}
			if inLen > max {
				isMoreThan = true
			}
			return
		}

		if inLen <= min {
			isLessThan = true
		}
		if inLen >= max {
			isMoreThan = true
		}
	case reflect.Float32, reflect.Float64:
		val := tv.Float()
		if hasEqual {
			if val < float64(min) {
				isLessThan = true
			}
			if val > float64(max) {
				isMoreThan = true
			}
			return
		} 

		if val <= float64(min) {
			isLessThan = true
		}
		if val >= float64(max) {
			isMoreThan = true
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := tv.Int()
		if hasEqual {
			if val < int64(min) {
				isLessThan = true
			}
			if val > int64(max) {
				isMoreThan = true
			}
			return
		} 

		if val <= int64(min) {
			isLessThan = true
		}
		if val >= int64(max) {
			isMoreThan = true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val := tv.Uint()
		if hasEqual {
			if val < uint64(min) {
				isLessThan = true
			}
			if val > uint64(max) {
				isMoreThan = true
			}
			return
		}
		if val < uint64(min) {
			isLessThan = true
		}
		if val > uint64(max) {
			isMoreThan = true
		}
	}
	return
}

// parseTagTo 解析 validName: to/oto 中 min, max
func parseTagTo(toVal string) (min int, max int, err error) {
	// 通过分割符来判断是否为区间
	toSlice := strings.Split(toVal, "~")
	l := len(toSlice)
	// fmt.Println("toSlice: ", toSlice)
	switch l {
	case 1:
		min, err = strconv.Atoi(toSlice[0])
	case 2:
		if min, err = strconv.Atoi(toSlice[0]); err != nil {
			return
		}

		if max, err = strconv.Atoi(toSlice[1]); err != nil {
			return
		}
	default:
		err = toValErr
	}
	return
}

// To 验证输入的大小区间, 如果为字符串则是验证字符个数, 如果是数字的话就验证数字的大小, 注: 左右都为闭区间
func To(errBuf *strings.Builder, validName, structName, filedName string, tv reflect.Value) {
	_, toVal := ParseValidNameKV(validName)
	min, max, err := parseTagTo(toVal)
	if err != nil {
		errBuf.WriteString(err.Error() + errEndFlag)
		return
	}

	isLessThan, isMoreThan, unitStr := validInputSize(min, max, tv)
	if isLessThan {
		// 生成如: "TestOrder.AppName" is len less than 2
		errBuf.WriteString("\"" + structName + "." + filedName + "\" is " + unitStr + " less than or equal " + fmt.Sprintf("%d", min) + errEndFlag)
	}

	if isMoreThan {
		// 生成如: "TestOrder.AppName" is len more than 30
		errBuf.WriteString("\"" + structName + "." + filedName + "\" is " + unitStr + " more than or equal " + fmt.Sprintf("%d", max) + errEndFlag)
	}
}

// Le 小于或等于验证, 如果为字符串则是验证字符个数, 如果是数字的话就验证数字的大小
func Le(errBuf *strings.Builder, validName, structName, filedName string, tv reflect.Value) {
	_, maxStr := ParseValidNameKV(validName)
	max, _ := strconv.Atoi(maxStr)
	_, isMoreThan, unitStr := validInputSize(0, max, tv)
	if isMoreThan {
		// 生成如: "TestOrder.AppName" is len more than 30
		errBuf.WriteString("\"" + structName + "." + filedName + "\" is " + unitStr + " more than or equal " + fmt.Sprintf("%d", max) + errEndFlag)
	}
}

// Ge 大于或等于验证, 如果为字符串则是验证字符个数, 如果是数字的话就验证数字的大小
func Ge(errBuf *strings.Builder, validName, structName, filedName string, tv reflect.Value) {
	_, minStr := ParseValidNameKV(validName)
	min, _ := strconv.Atoi(minStr)
	isLessThan, _, unitStr := validInputSize(min, 0, tv)
	if isLessThan {
		// 生成如: "TestOrder.AppName" is len less than 2
		errBuf.WriteString("\"" + structName + "." + filedName + "\" is " + unitStr + " less than or equal " + fmt.Sprintf("%d", min) + errEndFlag)
	}
}

// OTo 验证输入的大小区间, 如果为字符串则是验证字符个数, 如果是数字的话就验证数字的大小, 注: 左右都为开区间
func OTo(errBuf *strings.Builder, validName, structName, filedName string, tv reflect.Value) {
	_, toVal := ParseValidNameKV(validName)
	min, max, err := parseTagTo(toVal)
	if err != nil {
		errBuf.WriteString(err.Error() + errEndFlag)
		return
	}

	isLessThan, isMoreThan, unitStr := validInputSize(min, max, tv, false)
	if isLessThan {
		// 生成如: "TestOrder.AppName" is len less than 2
		errBuf.WriteString("\"" + structName + "." + filedName + "\" is " + unitStr + " less than " + fmt.Sprintf("%d", min) + errEndFlag)
	}

	if isMoreThan {
		// 生成如: "TestOrder.AppName" is len more than 30
		errBuf.WriteString("\"" + structName + "." + filedName + "\" is " + unitStr + " more than " + fmt.Sprintf("%d", max) + errEndFlag)
	}
}

// Lt 小于验证, 如果为字符串则是验证字符个数, 如果是数字的话就验证数字的大小
func Lt(errBuf *strings.Builder, validName, structName, filedName string, tv reflect.Value) {
	_, maxStr := ParseValidNameKV(validName)
	max, _ := strconv.Atoi(maxStr)
	_, isMoreThan, unitStr := validInputSize(0, max, tv, false)
	if isMoreThan {
		// 生成如: "TestOrder.AppName" is len more than 30
		errBuf.WriteString("\"" + structName + "." + filedName + "\" is " + unitStr + " more than " + fmt.Sprintf("%d", max) + errEndFlag)
	}
}

// Gt 大于验证, 如果为字符串则是验证字符个数, 如果是数字的话就验证数字的大小
func Gt(errBuf *strings.Builder, validName, structName, filedName string, tv reflect.Value) {
	_, minStr := ParseValidNameKV(validName)
	min, _ := strconv.Atoi(minStr)
	isLessThan, _, unitStr := validInputSize(min, 0, tv, false)
	if isLessThan {
		// 生成如: "TestOrder.AppName" is len less than 2
		errBuf.WriteString("\"" + structName + "." + filedName + "\" is " + unitStr + " less than " + fmt.Sprintf("%d", min) + errEndFlag)
	}
}

// In 指定输入选项
func In(errBuf *strings.Builder, validName, structName, filedName string, tv reflect.Value) {
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
		errBuf.WriteString("\"" + structName + "." + filedName + "\" should in (" + inVals + ")")
	}
}

// Phone 验证手机号
func Phone(errBuf *strings.Builder, validName, structName, filedName string, tv reflect.Value) {
	matched, _ := regexp.MatchString("^1[3,4,5,7,8,9]\\d{9}$", tv.String())
	if matched {
		return
	}
	errBuf.WriteString("\"" + structName + "." + filedName + "\" is not phone" + errEndFlag)
}

// Email 验证邮箱
func Email(errBuf *strings.Builder, validName, structName, filedName string, tv reflect.Value) {
	matched, _ := regexp.MatchString("^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$", tv.String())
	if matched {
		return
	}
	errBuf.WriteString("\"" + structName + "." + filedName + "\" is not email" + errEndFlag)
}

// IDCard 验证身份证
func IDCard(errBuf *strings.Builder, validName, structName, filedName string, tv reflect.Value) {
	matched, _ := regexp.MatchString("(^\\d{15}$)|(^\\d{18}$)|(^\\d{17}(\\d|X|x)$)", tv.String())
	if matched {
		return
	}
	errBuf.WriteString("\"" + structName + "." + filedName + "\" is not idcard" + errEndFlag)
}

// Date 验证日期
func Date(errBuf *strings.Builder, validName, structName, filedName string, tv reflect.Value) {
	matched, _ := regexp.MatchString("^\\d{4}-\\d{1,2}-\\d{1,2}", tv.String())
	if matched {
		return
	}
	errBuf.WriteString("\"" + structName + "." + filedName + "\" is not date, eg: 2021-11-15" + errEndFlag)
}

// Datetime 验证时间
func Datetime(errBuf *strings.Builder, validName, structName, filedName string, tv reflect.Value) {
	matched, _ := regexp.MatchString("^\\d{4}-\\d{1,2}-\\d{1,2} \\d{1,2}:\\d{1,2}:\\d{1,2}", tv.String())
	if matched {
		return
	}
	errBuf.WriteString("\"" + structName + "." + filedName + "\" is not datetime, eg: 2021-11-15 23:59:59" + errEndFlag)
}
