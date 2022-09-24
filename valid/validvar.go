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
	validFn map[string]CommonValidFn // 存放自定义的验证函数, 可以做到调用完就被清理
}

// NewVVar 单值校验
func NewVVar() *VVar {
	obj := syncValidVarPool.Get().(*VVar)
	if obj.errBuf == nil {
		obj.errBuf = new(strings.Builder)
	}
	// obj.errBuf.Grow(1 << 4)
	obj.ruleObj = NewRule()
	return obj
}

// free 释放
func (v *VVar) free() {
	v.errBuf.Reset()
	v.ruleObj = nil
	v.validFn = nil
	syncValidVarPool.Put(v)
}

// Valid 验证
// 对切片/数组/单个[int 系列, float系列, bool系列, string系列]进行验证
func (v *VVar) Valid(src interface{}) error {
	if src == nil {
		return errors.New("src is nil")
	}

	reflectValue := reflect.ValueOf(src)
	switch reflectValue.Kind() {
	case reflect.Ptr:
		if reflectValue.IsNil() {
			return errors.New("src \"" + reflectValue.Type().String() + "\" is nil")
		}
		// case reflect.Slice, reflect.Array:
		// 	for i := 0; i < reflectValue.Len(); i++ {
		// 		v.validate(reflectValue.Index(i))
		// 	}
		// 	return v.getError()
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
	if v.validFn == nil {
		v.validFn = make(map[string]CommonValidFn)
	}
	v.validFn[validName] = fn
	return v
}

// getValidFn 获取验证函数
func (v *VVar) getValidFn(validName string) (CommonValidFn, error) {
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
func (v *VVar) validate(value reflect.Value) *VVar {
	supportType := false
	tv := RemoveValuePtr(value)
	ty := tv.Type()

reValid:
	// 判断是否能进行验证
	switch ty.Kind() {
	case reflect.String:
		supportType = true
	case reflect.Slice, reflect.Array: // 再验证下里面的内容类型
		ty = ty.Elem()
		goto reValid
	default:
		if ReflectKindIsNum(ty.Kind()) {
			supportType = true
		}
	}
	if !supportType {
		v.errBuf.WriteString("src no support")
		return v
	}

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
