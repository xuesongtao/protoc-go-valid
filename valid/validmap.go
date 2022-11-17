package valid

import (
	"errors"
	"reflect"
	"strings"
)

type VMap struct {
	ruleObj RM
	errBuf  *strings.Builder
	vc      *validCommon
}

// NewVMap 验证 map
func NewVMap() *VMap {
	return &VMap{
		errBuf: newStrBuf(),
		vc:     new(validCommon),
	}
}

// SetRule 设置规则
func (v *VMap) SetRule(ruleObj RM) *VMap {
	v.ruleObj = ruleObj
	return v
}

// SetValidFn 自定义设置验证函数
func (v *VMap) SetValidFn(validName string, fn CommonValidFn) *VMap {
	v.vc.setValidFn(validName, fn)
	return v
}

// getValidFn 获取验证函数
func (v *VMap) getValidFn(validName string) (CommonValidFn, error) {
	return v.vc.getValidFn(validName)
}

// Valid 验证
// 支持 key 为 [string]
func (v *VMap) Valid(src interface{}) error {
	if src == nil {
		return errors.New("src is nil")
	}

	if len(v.ruleObj) == 0 {
		return errors.New("have no set rules")
	}

	tv := RemoveValuePtr(reflect.ValueOf(src))
	switch tv.Kind() {
	case reflect.Array, reflect.Slice:
		l := tv.Len()
		for i := 0; i < l; i++ {
			v.validate("["+ToStr(i)+"]", tv.Index(i))
		}
		return v.getError()
	}
	return v.validate("", tv).getError()
}

// validate 验证执行体
func (v *VMap) validate(prefix string, tv reflect.Value) *VMap {
	if tv.Type().Key().Kind() != reflect.String {
		v.errBuf.WriteString(GetJoinFieldErr("", prefix, "map key must string"))
		return v
	}

	if tv.Kind() != reflect.Map {
		v.errBuf.WriteString(GetJoinFieldErr("", prefix, "val must map"))
		return v
	}

	mapIter := tv.MapRange()
	for mapIter.Next() {
		key := mapIter.Key().String()
		val := mapIter.Value()
		validNames := v.ruleObj.Get(key)
		if validNames == "" {
			continue
		}

		// 根据验证内容进行验证
		for _, validName := range ValidNamesSplit(validNames) {
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
					if !val.IsZero() { // 验证必填
						continue
					}
					if cusMsg != "" {
						v.errBuf.WriteString(GetJoinValidErrStr("", v.getKey(prefix, key), "", cusMsg))
						continue
					}
					v.errBuf.WriteString(GetJoinValidErrStr("", v.getKey(prefix, key), "", ExplainEn, "it is", Required))
				case Either, BothEq:
					v.vc.initValid2FieldsMap(&name2Value{
						validName:  validName,
						fieldName:  key,
						cusMsg:     cusMsg,
						reflectVal: reflect.ValueOf(val),
					})
				default:
					v.errBuf.WriteString(GetJoinFieldErr("", v.getKey(prefix, key), "valid \""+validName+"\" is no support"))
				}
				continue
			}
			// 拓展的验证方法
			if val.IsZero() { // 空就直接跳过
				continue
			}
			fn(v.errBuf, validName, "", v.getKey(prefix, key), val)
		}
	}
	return v
}

// getKey 获取 key
func (v *VMap) getKey(prefix, key string) string {
	if prefix == "" && key == "" {
		return ""
	}
	if key == "" {
		return prefix + "map"
	}
	return prefix + "map[" + key + "]"
}

// getError 获取 err
func (v *VMap) getError() error {
	defer putStrBuf(v.errBuf)
	v.vc.valid(v.errBuf)
	if v.errBuf.Len() == 0 {
		return nil
	}
	// 这里需要去掉最后一个 ErrEndFlag
	return errors.New(strings.TrimSuffix(v.errBuf.String(), ErrEndFlag))
}
