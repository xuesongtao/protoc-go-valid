package valid

import "strings"

// RM 字段的自定义验证规则, key 为字段名, value 为验证规则
type RM map[string]string

func NewRule() RM {
	return make(map[string]string)
}

// Set 设置验证规则
// filedName 多个字段通过逗号隔开
// rules 多个字段通过逗号隔开
func (r RM) Set(filedNames string, rules string) RM {
	for _, filedName := range strings.Split(filedNames, ",") {
		// 如果存在的话就通过逗号隔开
		if _, ok := r[filedName]; ok {
			r[filedName] += "," + rules
			continue
		}
		r[filedName] = rules
	}
	return r
}

// Get 获取验证规则
func (r RM) Get(filedName string) string {
	if len(r) == 0 || filedName == "" {
		return ""
	}
	return r[filedName]
}

func (r RM) toLower(s string) string {
	strByte := []byte(s)
	l := len(strByte)
	for i := 0; i < l; i++ {
		strByte[i] |= ' '
	}
	return string(strByte)
}
