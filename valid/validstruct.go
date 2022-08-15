package valid

import (
	"errors"
	"reflect"
	"strings"
)

// VStruct 验证结构体
type VStruct struct {
	targetTag       string // 结构体中的待指定的验证的 tag
	errBuf          *strings.Builder
	ruleObj         RM                       // 验证规则
	valid2FieldsMap map[string][]*name2Value // 已存在的, 用于辅助 either, bothexist, botheq tag
	validFn         map[string]CommonValidFn // 存放自定义的验证函数, 可以做到调用完就被清理
}

// NewVStruct 验证结构体, 默认目标 tagName 为 "valid"
func NewVStruct(targetTag ...string) *VStruct {
	obj := syncValidStructPool.Get().(*VStruct)
	tagName := defaultTargetTag
	if len(targetTag) > 0 {
		tagName = targetTag[0]
	}
	obj.targetTag = tagName
	if obj.errBuf == nil { // 储存使用的时候 new 下, 后续都是从缓存中处理
		obj.errBuf = new(strings.Builder)
	}
	obj.errBuf.Grow(1 << 7)
	return obj
}

// free 释放
func (v *VStruct) free() {
	v.errBuf.Reset()
	v.ruleObj = nil
	v.valid2FieldsMap = nil
	v.validFn = nil
	syncValidStructPool.Put(v)
}

// SetRule 添加验证规则
func (v *VStruct) SetRule(ruleObj RM) Valider {
	v.ruleObj = ruleObj
	return v
}

// getCusRule 根据字段名获取自定义验证规则
func (v *VStruct) getCusRule(structFieldName string) string {
	if v.ruleObj == nil {
		return ""
	}
	return v.ruleObj.Get(structFieldName)
}

// Valid 验证
// 1. 支持单结构体验证
// 2. 支持切片/数组类型结构体验证
func (v *VStruct) Valid(src interface{}) error {
	if src == nil {
		return errors.New("src is nil")
	}

	reflectValue := reflect.ValueOf(src)
	switch reflectValue.Kind() {
	case reflect.Ptr:
		if reflectValue.IsNil() {
			return errors.New("src \"" + reflectValue.Type().String() + "\" is nil")
		}
	case reflect.Slice, reflect.Array:
		var structName string
		for i := 0; i < reflectValue.Len(); i++ {
			if i == 0 {
				structName = reflectValue.Index(i).Type().String()
			}
			v.validate(structName+"-"+ToStr(i), reflectValue.Index(i), true)
		}
		return v.getError()
	}
	return v.validate("", reflectValue).getError()
}

// SetValidFn 自定义设置验证函数
func (v *VStruct) SetValidFn(validName string, fn CommonValidFn) Valider {
	if v.validFn == nil {
		v.validFn = make(map[string]CommonValidFn)
	}
	v.validFn[validName] = fn
	return v
}

// getValidFn 获取验证函数
func (v *VStruct) getValidFn(validName string) (CommonValidFn, error) {
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
func (v *VStruct) validate(structName string, value reflect.Value, isValidSlice ...bool) *VStruct {
	// 辅助 errMsg, 用于嵌套时拼接上一级的结构体名
	if structName != "" {
		structName = structName + "."
	}

	tv := RemoveValuePtr(value)
	ty := tv.Type()

	// 如果不是结构体就退出
	if ty.Kind() != reflect.Struct {
		// 这里主要防止验证的切片为非结构体切片, 如 []int{1, 2, 3}, 这里会出现1, 为非指针所有需要退出
		if len(isValidSlice) > 0 && isValidSlice[0] {
			return v
		}
		v.errBuf.WriteString("src param \"" + structName + ty.Name() + "\" is not struct" + ErrEndFlag)
		return v
	}

	totalFieldNum := tv.NumField()
	for fieldNum := 0; fieldNum < totalFieldNum; fieldNum++ {
		structField := ty.Field(fieldNum)
		// 判断下是否可导出
		if !IsExported(structField.Name) {
			continue
		}
		validNames := structField.Tag.Get(v.targetTag)

		// 如果设置了规则就覆盖 tag 中的验证内容
		if rule := v.getCusRule(structField.Name); rule != "" {
			validNames = rule
		}

		// 没有 validNames 直接跳过
		if validNames == "" {
			continue
		}

		fieldValue := tv.Field(fieldNum)
		// 根据 tag 中的验证内容进行验证
		for _, validName := range strings.Split(validNames, ",") {
			if validName == "" {
				continue
			}

			validKey, _, cusMsg := ParseValidNameKV(validName)
			fn, err := v.getValidFn(validKey)
			if err != nil {
				v.errBuf.WriteString(GetJoinFieldErr(structName+ty.Name(), structField.Name, err))
				continue
			}

			// fmt.Printf("structName: %s, structFieldName: %s, tv: %v\n", structName+ty.Name(), structField.Name, fieldValue)
			// 开始验证
			// VStruct 内的验证方法
			if fn == nil {
				switch validKey {
				case Required:
					v.required(structName+ty.Name(), structField.Name, cusMsg, fieldValue)
				case Exist:
					v.exist(true, structName+ty.Name(), structField.Name, cusMsg, fieldValue)
				case Either, BothEq:
					v.initValid2FieldsMap(validName, structName+ty.Name(), structField.Name, cusMsg, fieldValue)
				}
				continue
			}

			// VStruct 外拓展的验证方法
			if fieldValue.IsZero() { // 空就直接跳过
				continue
			}
			fn(v.errBuf, validName, structName+ty.Name(), structField.Name, fieldValue)
		}
	}
	return v
}

// required 验证 required
func (v *VStruct) required(structName, fieldName, cusMsg string, tv reflect.Value) {
	ok := true
	// 如果集合类型先判断下长度
	switch tv.Kind() {
	case reflect.Slice, reflect.Array:
		if tv.Len() == 0 {
			ok = false
		}
	}

	if !ok || tv.IsZero() { // 验证必填
		if cusMsg != "" {
			v.errBuf.WriteString(GetJoinValidErrStr(structName, fieldName, "", cusMsg))
			return
		}
		// 生成如: "TestOrderDetailSlice.Price" is required
		v.errBuf.WriteString(GetJoinValidErrStr(structName, fieldName, "", ExplainEn, "it is", Required))
		return
	}

	// 有值的话再判断下嵌套的类型
	v.exist(false, structName, fieldName, cusMsg, tv)
}

// exist 存在验证, 用于验证嵌套结构, 切片
func (v *VStruct) exist(isValidTvKind bool, structName, fieldName, cusMsg string, tv reflect.Value) {
	// 如果空的就没必要验证了
	if tv.IsZero() {
		return
	}
	switch tv.Kind() {
	case reflect.Ptr, reflect.Struct:
		if structName == "Time" && tv.Type() == timeReflectType {
			return
		}
		v.validate(structName, tv)
	case reflect.Slice, reflect.Array:
		for i := 0; i < tv.Len(); i++ {
			v.validate(structName+"-"+ToStr(i), tv.Index(i), true)
		}
	default:
		if isValidTvKind {
			if cusMsg != "" {
				v.errBuf.WriteString(GetJoinValidErrStr(structName, fieldName, tv.String(), cusMsg))
				return
			}
			v.errBuf.WriteString(GetJoinValidErrStr(structName, fieldName, tv.String(), ExplainEn, "it is nonsupport", Exist))
		}
	}
}

// initValid2FieldsMap 为验证 either/bothexist/botheq 进行准备
func (v *VStruct) initValid2FieldsMap(validName, structName, fieldName, cusMsg string, tv reflect.Value) {
	if v.valid2FieldsMap == nil {
		v.valid2FieldsMap = make(map[string][]*name2Value, 5)
	}

	if _, ok := v.valid2FieldsMap[validName]; !ok {
		v.valid2FieldsMap[validName] = make([]*name2Value, 0, 2)
	}
	v.valid2FieldsMap[validName] = append(v.valid2FieldsMap[validName], &name2Value{structName: structName, fieldName: fieldName, cusMsg: cusMsg, reflectVal: tv})
}

// either 判断两者不能都为空
func (v *VStruct) either(fieldInfos []*name2Value) {
	l := len(fieldInfos)
	if l == 1 { // 如果只有 1 个就没有必要向下执行了
		info := fieldInfos[0]
		v.errBuf.WriteString(GetJoinFieldErr(info.structName, info.fieldName, eitherValErr))
		return
	}
	isZeroLen := 0
	fieldInfoStr := "" // 拼接空的 structName, fliedName
	for _, fieldInfo := range fieldInfos {
		fieldInfoStr += "\"" + fieldInfo.structName + "." + fieldInfo.fieldName + "\", "
		if fieldInfo.reflectVal.IsZero() {
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
func (v *VStruct) bothEq(fieldInfos []*name2Value) {
	l := len(fieldInfos)
	if l == 1 { // 如果只有 1 个就没有必要向下执行了
		info := fieldInfos[0]
		v.errBuf.WriteString(GetJoinFieldErr(info.structName, info.fieldName, bothEqValErr))
		return
	}

	var (
		tmp          interface{}
		fieldInfoStr string // 拼接空的 structName, fliedName
		eq           = true
	)
	for i, fieldInfo := range fieldInfos {
		fieldInfoStr += "\"" + fieldInfo.structName + "." + fieldInfo.fieldName + "\", "
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
		fieldInfoStr = strings.TrimSuffix(fieldInfoStr, ", ")
		v.errBuf.WriteString(fieldInfoStr + " " + ExplainEn + " they should be equal" + ErrEndFlag)
	}
}

// againValid 再一次验证
func (v *VStruct) againValid() {
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
func (v *VStruct) getError() error {
	defer v.free()

	v.againValid()

	if v.errBuf.Len() == 0 {
		return nil
	}

	// 这里需要去掉最后一个 ErrEndFlag
	return errors.New(strings.TrimSuffix(v.errBuf.String(), ErrEndFlag))
}
