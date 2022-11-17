package valid

import (
	"errors"
	"reflect"
	"strings"
)

// CacheEr 缓存接口
type CacheEr interface {
	Load(key interface{}) (interface{}, bool)
	Store(key, value interface{})
}

// name2Value
type name2Value struct {
	validName  string
	objName    string
	fieldName  string
	cusMsg     string
	reflectVal reflect.Value
}

// validCommon 验证体抽象类
type validCommon struct {
	validFn         map[string]CommonValidFn // 存放自定义的验证函数, 可以做到调用完就被清理
	valid2FieldsMap map[string][]*name2Value // 已存在的, 用于辅助 either, bothexist, botheq tag
}

// setValidFn 自定义设置验证函数
func (v *validCommon) setValidFn(validName string, fn CommonValidFn) {
	if v.validFn == nil {
		v.validFn = make(map[string]CommonValidFn)
	}
	v.validFn[validName] = fn
}

// getValidFn 获取验证函数
func (v *validCommon) getValidFn(validName string) (CommonValidFn, error) {
	// 先从本地找, 如果本地没有就从全局里找
	fn, ok := v.validFn[validName]
	if ok {
		return fn, nil
	}

	fn, ok = validName2FnMap[validName]
	if !ok {
		return nil, errors.New("valid \"" + validName + "\" is not exist, You can call SetValidFn")
	}
	return fn, nil
}

// initValid2FieldsMap 为验证 either/bothexist/botheq 进行准备
func (v *validCommon) initValid2FieldsMap(data *name2Value) {
	if data == nil {
		return
	}
	if v.valid2FieldsMap == nil {
		v.valid2FieldsMap = make(map[string][]*name2Value, 5)
	}
	if _, ok := v.valid2FieldsMap[data.validName]; !ok {
		v.valid2FieldsMap[data.validName] = make([]*name2Value, 0, 2)
	}
	v.valid2FieldsMap[data.validName] = append(v.valid2FieldsMap[data.validName], data)
}

// either 判断两者不能都为空
func (v *validCommon) either(errBuf *strings.Builder, fieldInfos []*name2Value) {
	l := len(fieldInfos)
	if l == 1 { // 如果只有 1 个就没有必要向下执行了
		info := fieldInfos[0]
		errBuf.WriteString(GetJoinFieldErr(info.objName, info.fieldName, eitherValErr))
		return
	}
	isZeroLen := 0
	fieldInfoBuf := newStrBuf(1 << 6) // 拼接空的 structName, fliedName
	defer putStrBuf(fieldInfoBuf)
	for _, fieldInfo := range fieldInfos {
		if fieldInfo.objName != "" {
			fieldInfoBuf.WriteString("\"" + fieldInfo.objName + ".")
		} else {
			fieldInfoBuf.WriteByte('"')
		}
		fieldInfoBuf.WriteString(fieldInfo.fieldName + "\", ")
		if fieldInfo.reflectVal.IsZero() {
			isZeroLen++
		}
	}

	// 判断下是否全部为空
	if l == isZeroLen {
		errBuf.WriteString(strings.TrimSuffix(fieldInfoBuf.String(), ", ") + " " + ExplainEn + " they shouldn't all be empty" + ErrEndFlag)
	}
}

// bothEq 判断两者相等
func (v *validCommon) bothEq(errBuf *strings.Builder, fieldInfos []*name2Value) {
	l := len(fieldInfos)
	if l == 1 { // 如果只有 1 个就没有必要向下执行了
		info := fieldInfos[0]
		errBuf.WriteString(GetJoinFieldErr(info.objName, info.fieldName, bothEqValErr))
		return
	}

	var (
		tmp          interface{}
		fieldInfoBuf = newStrBuf(1 << 8) // 拼接空的 structName, fliedName
		eq           = true
	)
	defer putStrBuf(fieldInfoBuf)
	for i, fieldInfo := range fieldInfos {
		if fieldInfo.objName != "" {
			fieldInfoBuf.WriteString("\"" + fieldInfo.objName + ".")
		} else {
			fieldInfoBuf.WriteByte('"')
		}
		fieldInfoBuf.WriteString(fieldInfo.fieldName + "\", ")
		if !eq { // 避免多次比较
			continue
		}

		if i == 0 {
			tmp = fieldInfo.reflectVal.Interface()
			continue
		}

		if !reflect.DeepEqual(tmp, fieldInfo.reflectVal.Interface()) {
			eq = false
		}
	}

	if !eq {
		errBuf.WriteString(strings.TrimSuffix(fieldInfoBuf.String(), ", ") + " " + ExplainEn + " they should be equal" + ErrEndFlag)
	}
}

// valid 验证
func (v *validCommon) valid(errBuf *strings.Builder) {
	// 判断下是否有值, 有就说明有 either 验证
	if len(v.valid2FieldsMap) == 0 {
		return
	}

	for validName, fieldInfos := range v.valid2FieldsMap {
		validKey, _, _ := ParseValidNameKV(validName)
		switch validKey {
		case Either:
			v.either(errBuf, fieldInfos)
		case BothEq:
			v.bothEq(errBuf, fieldInfos)
		}
	}
}
