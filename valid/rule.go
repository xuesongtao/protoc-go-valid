package valid

import "strings"

// RM 字段的自定义验证规则, key 为字段名, value 为验证规则
type RM map[string]string

func NewRule() RM {
	return make(RM, 4)
}

// Set 设置验证规则
// fieldName 多个字段通过逗号隔开
// rules 多个字段通过逗号隔开
func (r RM) Set(filedNames string, rules ...string) RM {
	for _, fieldName := range strings.Split(filedNames, ",") {
		// 如果存在的话就通过逗号隔开
		if _, ok := r[fieldName]; ok {
			r[fieldName] += "," + strings.Join(rules, ",")
			continue
		}
		r[fieldName] = strings.Join(rules, ",")
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

// Deprecated: 名字存在歧义, 因为已上线不能删除, 特此标记, 推荐使用 GenValidKV
// JoinTag2Val 生成 defaultTargetTag 的值
func JoinTag2Val(key string, values ...string) string {
	return GenValidKV(key, values...)
}

// GenValidKV 生成 defaultTargetTag 的值
// 说明: 函数名主要用于生成(如: `valid:"xxx"`) 中 "xxx" 的部分
// key 为验证规则
// values[0] 会被解析为值
// values[1] 会被解析为自定义错误信息
// 如1.: GenValidKV(VTo, "1~10", "需要在 1-10 的区间")
// => to=1~10|需要在 1-10 的区间
//
// 如2: GenValidKV(VRe, "'\\d+'", "必须为纯数字")
// => re='\\d+'|必须为纯数字
func GenValidKV(key string, values ...string) string {
	l := len(values)
	if l == 0 {
		return key
	}

	buf := newStrBuf()
	defer putStrBuf(buf)
	buf.Grow(1 << 4)
	buf.WriteString(key)
	if values[0] != "" {
		// 判断第一个值得首字符是否为 "="
		if values[0][0] != '=' {
			buf.WriteByte('=')
		}

		// 处理 val 前缀
		// 说明: 为了兼容老版本 in, include, re 特此只处理了3个
		switch key {
		case VIn, VInclude:
			buf.WriteString("(" + values[0] + ")")
		case VRe:
			if len(values[0]) > 1 && (values[0][0] == '\'' || values[0][1] == '\'') {
				buf.WriteString(values[0])
			} else {
				buf.WriteString("'" + values[0] + "'")
			}
		default:
			buf.WriteString(values[0])
		}
	}

	// 自定义说明
	if l >= 2 {
		buf.WriteString("|" + values[1])
	}
	return buf.String()
}
