package valid

import (
	"fmt"
	"reflect"
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

// checkfieldIsString 验证字段类型是否为字符串
func checkFieldIsString(validName, objName, filedName string, tv reflect.Value) (err error) {
	switch tv.Kind() {
	case reflect.String:
	default:
		err = fmt.Errorf("\"" + objName + "." + filedName + " [" + tv.String() + "] " + "must is string")
	}
	return
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
			if min > 0 && inLen < min {
				isLessThan = true
			}
			if max > 0 && inLen > max {
				isMoreThan = true
			}
			return
		}

		if min > 0 && inLen <= min {
			isLessThan = true
		}
		if max > 0 && inLen >= max {
			isMoreThan = true
		}
	case reflect.Float32, reflect.Float64:
		val := tv.Float()
		valStr = fmt.Sprintf("%v", val)
		if hasEqual {
			if min > 0 && val < float64(min) {
				isLessThan = true
			}
			if max > 0 && val > float64(max) {
				isMoreThan = true
			}
			return
		}

		if min > 0 && val <= float64(min) {
			isLessThan = true
		}
		if max > 0 && val >= float64(max) {
			isMoreThan = true
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := tv.Int()
		valStr = fmt.Sprintf("%d", val)
		if hasEqual {
			if min > 0 && val < int64(min) {
				isLessThan = true
			}
			if max > 0 && val > int64(max) {
				isMoreThan = true
			}
			return
		}

		if min > 0 && val <= int64(min) {
			isLessThan = true
		}
		if max > 0 && val >= int64(max) {
			isMoreThan = true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val := tv.Uint()
		valStr = fmt.Sprintf("%d", val)
		if hasEqual {
			if min > 0 && val < uint64(min) {
				isLessThan = true
			}
			if max > 0 && val > uint64(max) {
				isMoreThan = true
			}
			return
		}
		if min > 0 && val < uint64(min) {
			isLessThan = true
		}
		if max > 0 && val > uint64(max) {
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
