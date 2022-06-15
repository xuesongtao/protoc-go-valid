package valid

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// ParseValidNameKV 解析 validName 中的 key, value 和 cusMsg,
// 如: "required|必填", key 为 "required", value 为 "", cusMsg 为 "必填"
// 如: "to=1~2|大于等于 1 且小于等于 2", key 为 "to", value 为 "1~2", cusMsg 为 "大于等于 1 且小于等于 2"
func ParseValidNameKV(validName string) (key, value, cusMsg string) {
	// 因为 validName 中的 k, v 通过 = 连接
	splitIndex := strings.Index(validName, "=")

	// 如果没有则代表 validName 不为 k=v 类型, 只有一个字段如: required
	if splitIndex == -1 {
		// 需要确定下是否包含自定义 msg, 格式为: validName|xxx, 如: required|必填
		key = validName
		cusMsgIndex := strings.Index(validName, "|")
		if cusMsgIndex != -1 && len(validName)-1 > cusMsgIndex+1 {
			key = validName[:cusMsgIndex]
			cusMsg = validName[cusMsgIndex+1:]
			// 根据如果说明有中文就加前缀为: 说明; 否则为 Explain
			if match, _ := regexp.MatchString("[\u4e00-\u9fa5]", cusMsg); match {
				cusMsg = "说明: " + cusMsg
			} else {
				cusMsg = "Explain: " + cusMsg
			}
		}
		return
	}

	key = validName[:splitIndex]
	value = validName[splitIndex+1:]
	// 需要确定下是否包含自定义 msg, 格式为: validName|xxx, 如: "to=1~2|大于等于 1 且小于等于 2"
	cusMsgIndex := strings.Index(value, "|")
	if cusMsgIndex != -1 && len(value)-1 > cusMsgIndex+1 {
		// 根据如果说明有中文就加前缀为: 说明; 否则为 Explain
		cusMsg = value[cusMsgIndex+1:]
		if match, _ := regexp.MatchString("[\u4e00-\u9fa5]", cusMsg); match {
			cusMsg = "说明: " + cusMsg
		} else {
			cusMsg = "Explain: " + cusMsg
		}
		value = value[:cusMsgIndex]
	}
	return
}

// GetJoinValidErrStr 获取拼接验证的错误消息, 内容直接通过空格隔开, 最后会拼接 ErrEndFlag
func GetJoinValidErrStr(objName, fieldName, inputVal string, others ...string) (res string) {
	res = "\"" + objName + "." + fieldName + "\" input \"" + inputVal + "\" "
	if len(others) == 0 {
		res += ErrEndFlag
		return
	}

	lastIndex := len(others) - 1
	for i, content := range others {
		if i < lastIndex {
			res += content + " "
			continue
		}
		res += content + ErrEndFlag
	}
	return
}

// CheckFieldIsStr 验证字段类型是否为字符串
func CheckFieldIsStr(objName, fieldName string, tv reflect.Value) (err error) {
	switch tv.Kind() {
	case reflect.String:
	default:
		err = fmt.Errorf(GetJoinValidErrStr(objName, fieldName, tv.String()) + "must is string")
	}
	return
}

// IsExported 是可导出
func IsExported(fieldName string) bool {
	if fieldName == "" {
		return false
	}
	first := fieldName[0]
	return first >= 'A' && first <= 'Z'
}

// validInputSize 验证输入的大小
func validInputSize(min, max int, tv reflect.Value, isHasEqual ...bool) (isLessThan, isMoreThan bool, valStr, unitStr string) {
	hasEqual := true // 标记对结果默认包含闭区间
	if len(isHasEqual) > 0 {
		hasEqual = isHasEqual[0]
	}
	unitStr = numUnitStr // 大小的单位
	// fmt.Printf("min: %v, max: %v\n", min, max)
	switch tv.Kind() {
	case reflect.String:
		unitStr = strUnitStr
		valStr = tv.String()
		inLen := len([]rune(valStr))
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
		valStr = fmt.Sprintf("%v", val)
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
		valStr = fmt.Sprintf("%d", val)
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
		valStr = fmt.Sprintf("%d", val)
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
	case reflect.Slice:
		unitStr = sliceLenUnitStr
		l := tv.Len()
		valStr = fmt.Sprintf("%d", l)
		if hasEqual {
			if l < min {
				isLessThan = true
			}
			if l > max {
				isMoreThan = true
			}
			return
		}
		if l < min {
			isLessThan = true
		}
		if l > max {
			isMoreThan = true
		}
	}
	return
}

// parseTagTo 解析 validName: to/oto 中 min, max
func parseTagTo(toVal string, isHasEqual bool) (min int, max int, err error) {
	// 通过分割符来判断是否为区间
	toSlice := strings.Split(toVal, "~")
	l := len(toSlice)
	// fmt.Println("toSlice: ", toSlice)
	if l != 2 {
		if isHasEqual {
			err = toValErr
		} else {
			err = otoValErr
		}
		return
	}

	if min, err = strconv.Atoi(toSlice[0]); err != nil {
		return
	}

	if max, err = strconv.Atoi(toSlice[1]); err != nil {
		return
	}
	return
}

// RemoveTypePtr 移除多指针
func RemoveTypePtr(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// RemoveValuePtr 移除多指针
func RemoveValuePtr(t reflect.Value) reflect.Value {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
