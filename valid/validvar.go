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
}

// NewVVar 单值校验
func NewVVar() *VVar {
	obj := syncValidVarPool.Get().(*VVar)
	if obj.errBuf == nil {
		obj.errBuf = new(strings.Builder)
	}
	obj.errBuf.Grow(1 << 6)
	obj.ruleObj = make(RM)
	return obj
}

// free 释放
func (v *VVar) free() {
	v.errBuf.Reset()
	v.ruleObj = nil
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
	for _, validName := range strings.Split(validNames, ",") {
		if validName == "" {
			continue
		}

		validKey, _, cusMsg := ParseValidNameKV(validName)
		fn, ok := validName2FuncMap[validKey]
		if !ok {
			v.errBuf.WriteString(GetJoinFieldErr("", "", "valid \""+validKey+"\" is not exist"))
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
				if ok || !tv.IsZero() {
					continue
				}
				if cusMsg != "" {
					v.errBuf.WriteString(GetJoinValidErrStr("", "", "", cusMsg))
					continue
				}
				v.errBuf.WriteString(GetJoinValidErrStr("", "", "", ExplainEn, "it is", Required))
			default:
				v.errBuf.WriteString(GetJoinFieldErr("", "", "valid \""+validName+"\" is no support"))
				continue
			}

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

// =========================== 常用方法进行封装 =======================================

// Var 验证变量
func Var(src interface{}, rules ...string) error {
	return NewVVar().SetRules(rules...).Valid(src)
}
