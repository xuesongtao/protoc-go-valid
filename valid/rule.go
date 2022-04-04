package valid

import "strings"

// RM 字段的自定义验证规则, key 为字段名, value 为验证规则
type RM map[string]string

func NewRule() RM {
	return make(map[string]string)
}

// Set 设置验证规则
// fieldName 多个字段通过逗号隔开
// rules 多个字段通过逗号隔开
func (r RM) Set(filedNames string, rules string) RM {
	for _, fieldName := range strings.Split(filedNames, ",") {
		// 如果存在的话就通过逗号隔开
		if _, ok := r[fieldName]; ok {
			r[fieldName] += "," + rules
			continue
		}
		r[fieldName] = rules
	}
	return r
}

// Get 获取验证规则
func (r RM) Get(fieldName string) string {
	if len(r) == 0 || fieldName == "" {
		return ""
	}
	return r[fieldName]
}
