package valid

import (
	"errors"
	"reflect"
	"strings"
)

const (
	validVarFieldName = "validVar"
)

// VVar 验证单字段
type VVar struct {
	errBuf  *strings.Builder
	ruleObj RM
	vc      *validCommon // 组合验证
}

// NewVVar 单值校验
func NewVVar() *VVar {
	obj := syncValidVarPool.Get().(*VVar)
	obj.errBuf = newStrBuf()
	obj.vc = &validCommon{}
	obj.ruleObj = NewRule()
	return obj
}

// free 释放
func (v *VVar) free() {
	v.errBuf.Reset()
	v.ruleObj = nil
	v.vc = nil
	syncValidVarPool.Put(v)
}

// Valid 验证
// 对切片/数组/单个[int,float,bool,string,struct]进行验证
func (v *VVar) Valid(src interface{}) error {
	if src == nil {
		return errors.New("src is nil")
	}

	reflectValue := RemoveValuePtr(reflect.ValueOf(src))
	ty := reflectValue.Type()
	supportType := false

again:
	// 判断是否能进行验证
	switch kind := ty.Kind(); kind {
	case reflect.String:
		supportType = true
	case reflect.Ptr, reflect.Slice, reflect.Array: // 再验证下里面的内容类型
		ty = ty.Elem()
		goto again
	case reflect.Struct:
		return Struct(src)
	default:
		if ReflectKindIsNum(kind, true) {
			supportType = true
		}
	}
	if !supportType {
		return errors.New("src no support")
	}
	return v.validate(reflectValue).getError()
}

// SetRules 设置规则
func (v *VVar) SetRules(rules ...string) *VVar {
	v.ruleObj.Set(validVarFieldName, rules...)
	return v
}

// SetValidFn 自定义设置验证函数
func (v *VVar) SetValidFn(validName string, fn CommonValidFn) *VVar {
	v.vc.setValidFn(validName, fn)
	return v
}

// getValidFn 获取验证函数
func (v *VVar) getValidFn(validName string) (CommonValidFn, error) {
	return v.vc.getValidFn(validName)
}

// validate 验证执行体
func (v *VVar) validate(tv reflect.Value) *VVar {
	validNames := v.ruleObj.Get(validVarFieldName)
	if validNames == "" {
		v.errBuf.WriteString("you no set rule")
		return v
	}

	// 根据验证内容进行验证
	for _, validName := range ValidNamesSplit(validNames) {
		if validName == "" {
			continue
		}

		validKey, _, cusMsg := ParseValidNameKV(validName)
		fn, err := v.getValidFn(validKey)
		if err != nil {
			v.errBuf.WriteString(GetJoinFieldErr("", "", err))
			continue
		}

		// 开始验证
		if fn == nil {
			switch validKey {
			case Required:
				ok := true
				// 如果集合类型先判断下长度
				switch tv.Kind() {
				case reflect.Array, reflect.Slice:
					if tv.Len() == 0 {
						ok = false
					}
				}
				if ok && !tv.IsZero() {
					continue
				}
				if cusMsg != "" {
					v.errBuf.WriteString(GetJoinValidErrStr("", "", "", cusMsg))
					continue
				}
				v.errBuf.WriteString(GetJoinValidErrStr("", "", "", ExplainEn, "it is", Required))
			default:
				v.errBuf.WriteString(GetJoinFieldErr("", "", "valid \""+validName+"\" is no support"))
			}
			continue
		}
		// 拓展的验证方法
		if tv.IsZero() { // 空就直接跳过
			continue
		}
		fn(v.errBuf, validName, "", "", tv)
	}
	return v
}

// getError 获取 err
func (v *VVar) getError() error {
	defer v.free()

	if v.errBuf.Len() == 0 {
		return nil
	}

	// 这里需要去掉最后一个 ErrEndFlag
	return errors.New(strings.TrimSuffix(v.errBuf.String(), ErrEndFlag))
}
