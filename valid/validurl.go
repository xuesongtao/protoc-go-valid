package valid

import (
	"errors"
	"net/url"
	"reflect"
	"strings"
)

// VUrl 验证 url
type VUrl struct {
	errBuf          *strings.Builder
	ruleObj         RM
	valid2FieldsMap map[string][]*name2Value // 已存在的, 用于辅助 either, bothexist, botheq tag
	validFn         map[string]CommonValidFn // 存放自定义的验证函数, 可以做到调用完就被清理
}

// NewVUrl
func NewVUrl() *VUrl {
	obj := new(VUrl)
	obj.errBuf = new(strings.Builder)
	obj.errBuf.Grow(1 << 4)
	return obj
}

// SetRule 添加验证规则
func (v *VUrl) SetRule(ruleObj RM) *VUrl {
	v.ruleObj = ruleObj
	return v
}

// Valid 验证
func (v *VUrl) Valid(src interface{}) error {
	if src == nil {
		return errors.New("src is nil")
	}

	var srcStr *string
	switch v := src.(type) {
	case string:
		srcStr = &v
	case *string:
		srcStr = v
	default:
		return errors.New("src must is string/*string")
	}
	return v.validate(srcStr).getError()
}

// SetValidFn 自定义设置验证函数
func (v *VUrl) SetValidFn(validName string, fn CommonValidFn) *VUrl {
	if v.validFn == nil {
		v.validFn = make(map[string]CommonValidFn)
	}
	v.validFn[validName] = fn
	return v
}

// getValidFn 获取验证函数
func (v *VUrl) getValidFn(validName string) (CommonValidFn, error) {
	// 先从本地找, 如果本地没有就从全局里找
	fn, ok := v.validFn[validName]
	if ok {
		return fn, nil
	}

	fn, ok = validName2FuncMap[validName]
	if !ok {
		return nil, errors.New("valid \"" + validName + "\" is not exist, You can call SetValidFn")
	}
	return fn, nil
}

// validate 验证执行体
func (v *VUrl) validate(value *string) *VUrl {
	// 解码处理
	decUrl, err := url.QueryUnescape(*value)
	if err != nil {
		v.errBuf.WriteString("url unescape is failed, err: " + err.Error())
		return v
	}
	queryIndex := strings.Index(decUrl, "?")
	urlQuerys := ""
	if queryIndex != -1 {
		urlQuerys = decUrl[queryIndex+1:]
	}
	if urlQuerys == "" {
		return v
	}

	var key, val string
	for _, query := range strings.Split(urlQuerys, "&") {
		key = ""
		val = ""
		key2val := strings.Split(query, "=")
		l := len(key2val)
		if l > 0 {
			key = key2val[0]
			if l > 1 {
				val = key2val[1]
			}
		}

		validNames := v.ruleObj.Get(key)
		if validNames == "" {
			continue
		}
		// 根据验证内容进行验证
		for _, validName := range strings.Split(validNames, ",") {
			if validName == "" {
				continue
			}

			validKey, _, cusMsg := ParseValidNameKV(validName)
			fn, err := v.getValidFn(validKey)
			if err != nil {
				v.errBuf.WriteString(GetJoinFieldErr("", key, err))
				continue
			}

			// 开始验证
			if fn == nil {
				switch validKey {
				case Required:
					v.required(key, cusMsg, val)
				case Either, BothEq:
					v.initValid2FieldsMap(validName, key, cusMsg, val)
				default:
					v.errBuf.WriteString(GetJoinFieldErr("", key, "valid \""+validName+"\" is no support"))
				}
				continue

			}
			// 拓展的验证方法
			if val == "" { // 空就直接跳过
				continue
			}
			fn(v.errBuf, validName, "", key, reflect.ValueOf(val))
		}
	}
	return v
}

// required 验证 required
func (v *VUrl) required(fieldName, cusMsg, val string) {
	if val != "" { // 验证必填
		return
	}
	if cusMsg != "" {
		v.errBuf.WriteString(GetJoinValidErrStr("", fieldName, "", cusMsg))
		return
	}
	v.errBuf.WriteString(GetJoinValidErrStr("", fieldName, "", ExplainEn, "it is", Required))
}

// initValid2FieldsMap 为验证 either/bothexist/botheq 进行准备
func (v *VUrl) initValid2FieldsMap(validName, fieldName, cusMsg, val string) {
	if v.valid2FieldsMap == nil {
		v.valid2FieldsMap = make(map[string][]*name2Value, 5)
	}

	if _, ok := v.valid2FieldsMap[validName]; !ok {
		v.valid2FieldsMap[validName] = make([]*name2Value, 0, 2)
	}
	v.valid2FieldsMap[validName] = append(v.valid2FieldsMap[validName], &name2Value{fieldName: fieldName, cusMsg: cusMsg, val: val})
}

// either 判断两者不能都为空
func (v *VUrl) either(fieldInfos []*name2Value) {
	l := len(fieldInfos)
	if l == 1 { // 如果只有 1 个就没有必要向下执行了
		info := fieldInfos[0]
		v.errBuf.WriteString(GetJoinFieldErr(info.structName, info.fieldName, eitherValErr))
		return
	}
	isZeroLen := 0
	fieldInfoStr := "" // 拼接空的 structName, fliedName
	for _, fieldInfo := range fieldInfos {
		fieldInfoStr += "\"" + fieldInfo.fieldName + "\", "
		if fieldInfo.val == "" {
			isZeroLen++
		}
	}

	// 判断下是否全部为空
	if l == isZeroLen {
		fieldInfoStr = strings.TrimSuffix(fieldInfoStr, ", ")
		v.errBuf.WriteString(fieldInfoStr + " " + ExplainEn + " they shouldn't all be empty" + ErrEndFlag)
	}
}

// bothEq 判断两者相等
func (v *VUrl) bothEq(fieldInfos []*name2Value) {
	l := len(fieldInfos)
	if l == 1 { // 如果只有 1 个就没有必要向下执行了
		info := fieldInfos[0]
		v.errBuf.WriteString(GetJoinFieldErr(info.structName, info.fieldName, bothEqValErr))
		return
	}

	var (
		tmp          string
		fieldInfoStr string // 拼接空的 structName, fliedName
		eq           = true
	)
	for i, fieldInfo := range fieldInfos {
		fieldInfoStr += "\"" + fieldInfo.fieldName + "\", "
		if !eq { // 避免多次比较
			continue
		}

		if i == 0 {
			tmp = fieldInfo.val
			continue
		}

		if tmp != fieldInfo.val {
			eq = false
		}
	}

	if !eq {
		fieldInfoStr = strings.TrimSuffix(fieldInfoStr, ", ")
		v.errBuf.WriteString(fieldInfoStr + " " + ExplainEn + " they should be equal" + ErrEndFlag)
	}
}

// againValid 再一次验证
func (v *VUrl) againValid() {
	// 判断下是否有值, 有就说明有 either 验证
	if len(v.valid2FieldsMap) == 0 {
		return
	}

	for validName, fieldInfos := range v.valid2FieldsMap {
		validKey, _, _ := ParseValidNameKV(validName)
		switch validKey {
		case Either:
			v.either(fieldInfos)
		case BothEq:
			v.bothEq(fieldInfos)
		}
	}
}

// getError 获取 err
func (v *VUrl) getError() error {
	v.againValid()

	if v.errBuf.Len() == 0 {
		return nil
	}

	// 这里需要去掉最后一个 ErrEndFlag
	return errors.New(strings.TrimSuffix(v.errBuf.String(), ErrEndFlag))
}
