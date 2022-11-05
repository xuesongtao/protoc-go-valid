package valid

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"gitee.com/xuesongtao/protoc-go-valid/valid/internal"
)

// name2Value
type name2Value struct {
	objName    string
	fieldName  string
	cusMsg     string
	val        string
	reflectVal reflect.Value
}

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
			if match := IncludeZhRe.MatchString(cusMsg); match {
				cusMsg = ExplainZh + " " + cusMsg
			} else {
				cusMsg = ExplainEn + " " + cusMsg
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
		if match := IncludeZhRe.MatchString(cusMsg); match {
			cusMsg = ExplainZh + " " + cusMsg
		} else {
			cusMsg = ExplainEn + " " + cusMsg
		}
		value = value[:cusMsgIndex]
	}
	return
}

// GetJoinFieldErr 拼接字段错误
func GetJoinFieldErr(objName, fieldName string, err interface{}) string {
	res := ""
	if objName != "" && fieldName != "" {
		res += "\"" + objName + "." + fieldName + "\" "
	}
	switch v := err.(type) {
	case string:
		res += v
	case error:
		res += v.Error()
	}
	return res + ErrEndFlag
}

// GetJoinValidErrStr 获取拼接验证的错误消息, 内容直接通过空格隔开, 最后会拼接 ErrEndFlag
func GetJoinValidErrStr(objName, fieldName, inputVal string, others ...string) string {
	res := new(strings.Builder)
	if objName != "" && fieldName != "" {
		res.WriteString("\"" + objName + "." + fieldName + "\" ")
	} else if objName == "" && fieldName != "" {
		res.WriteString("\"" + fieldName + "\" ")
	}
	res.WriteString("input \"" + inputVal + "\"")
	if len(others) == 0 {
		res.WriteString(ErrEndFlag)
		return res.String()
	}

	res.WriteString(", ")
	// 判断下是否需要注入: ExplainEn
	if !strings.Contains(others[0], ExplainEn) && !strings.Contains(others[0], ExplainZh) {
		res.WriteString(ExplainEn + " ")
	}
	lastIndex := len(others) - 1
	for i, content := range others {
		if i < lastIndex {
			res.WriteString(content + " ")
			continue
		}
		res.WriteString(content + ErrEndFlag)
	}
	return res.String()
}

// CheckFieldIsStr 验证字段类型是否为字符串
func CheckFieldIsStr(objName, fieldName string, tv reflect.Value) (err error) {
	switch tv.Kind() {
	case reflect.String:
	default:
		err = fmt.Errorf(GetJoinValidErrStr(objName, fieldName, tv.String(), ExplainEn, "it must is string"))
	}
	return
}

// ReflectKindIsNum 值是否为数字
func ReflectKindIsNum(kind reflect.Kind) (is bool) {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		is = true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		is = true
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
		valStr = ToStr(val)
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
		valStr = ToStr(val)
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
		valStr = ToStr(val)
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
		valStr = ToStr(l)
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

// ToStr 将内容转为 string
func ToStr(src interface{}) string {
	if src == nil {
		return ""
	}

	switch value := src.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.Itoa(int(value))
	case int16:
		return strconv.Itoa(int(value))
	case int32:
		return strconv.Itoa(int(value))
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case string:
		return value
	case []byte:
		return string(value)
	default:
		return fmt.Sprintf("%v", value)
	}
}

// ValidNamesSplit 验证点进行分割
func ValidNamesSplit(s string, sep ...byte) []string {
	if s == "" {
		return nil
	}

	defaultSep := byte(',')
	if len(sep) > 0 {
		defaultSep = sep[0]
	}

	// 判断是否包含单单引号, 包含单引号需要排除里面的逗号
	// fast path
	if strings.IndexByte(s, '\'') == -1 {
		return strings.Split(s, string(defaultSep))
	}

	stack := internal.NewStackByte(2)
	defer stack.Reset()

	l := len(s)
	tmp := make([]byte, 0, 5)
	res := make([]string, 0, 3)
	isParseSingleQuotes := false // 标记是否正在处理单引号里的内容
	for i := 0; i < l; i++ {
		v := s[i]
		// 非单引号
		if !isParseSingleQuotes && v != defaultSep {
			tmp = append(tmp, v)
		} else if isParseSingleQuotes { // 单引号中的所有内容都添加
			tmp = append(tmp, v)
		}

		// 判断是否为 单引号
		if !isParseSingleQuotes && v == '\'' {
			stack.Append(v)
			isParseSingleQuotes = true
			continue
		}

		// 判断单引号结束
		if isParseSingleQuotes && stack.IsEqualLastVal(v) {
			stack.Pop()
			isParseSingleQuotes = false
			continue
		}

		if v == defaultSep && stack.IsEmpty() {
			res = append(res, string(tmp))
			tmp = tmp[:0]
		}
	}

	if len(tmp) > 0 {
		res = append(res, string(tmp))
	}
	return res
}
