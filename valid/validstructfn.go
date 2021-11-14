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

// To 验证输入的大小区间, 如果为字符串则是验证字符个数, 如果是数字的话就验证数字的大小
func To(errBuf *strings.Builder, validName, structName, filedName string, tv reflect.Value) {
	// parseTagTo 解析 validName: to 中 min, max
	parseTagTo := func(toVal string) (min int, max int, err error) {
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

	_, toVal := ParseValidNameKV(validName)
	min, max, err := parseTagTo(toVal)
	if err != nil {
		errBuf.WriteString(err.Error() + errEndFlag)
		return
	}

	var (
		unitStr                = "size" // 大小的单位
		isLessThan, isMoreThan bool
	)
	// fmt.Printf("min: %v, max: %v\n", min, max)
	switch tv.Kind() {
	case reflect.String:
		unitStr = "len"
		inLen := len([]rune(tv.String()))
		if min > 0 && inLen < min {
			isLessThan = true
		}

		if max > 0 && inLen > max {
			isMoreThan = true
		}
	case reflect.Float32, reflect.Float64:
		val := tv.Float()
		if min > 0 && val < float64(min) {
			isLessThan = true
		}

		if max > 0 && val > float64(max) {
			isMoreThan = true
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := tv.Int()
		if min > 0 && val < int64(min) {
			isLessThan = true
		}

		if max > 0 && val > int64(max) {
			isMoreThan = true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val := tv.Uint()
		if min > 0 && val < uint64(min) {
			isLessThan = true
		}

		if max > 0 && val > uint64(max) {
			isMoreThan = true
		}
	}

	if isLessThan {
		// 生成如: "TestOrder.AppName" is len less than 2
		errBuf.WriteString("\"" + structName + "." + filedName + "\" is " + unitStr + " less than " + fmt.Sprintf("%d", min) + errEndFlag)
	}

	if isMoreThan {
		// 生成如: "TestOrder.AppName" is len more than 30
		errBuf.WriteString("\"" + structName + "." + filedName + "\" is " + unitStr + " more than " + fmt.Sprintf("%d", max) + errEndFlag)
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