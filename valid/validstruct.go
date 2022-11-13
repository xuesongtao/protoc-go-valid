package valid

import (
	"errors"
	"reflect"
	"strings"
)

var (
	validOnlyOuterObj = reflect.TypeOf("validOnlyOuterObj") // 标记只验证最外层的结构体
)

// VStruct 验证结构体
type VStruct struct {
	targetTag       string // 结构体中的待指定的验证的 tag
	errBuf          *strings.Builder
	ruleMap         map[reflect.Type]RM      // 验证规则, k: 为结构体 reflect.Type, v: 为该结构体的规则
	valid2FieldsMap map[string][]*name2Value // 已存在的, 用于辅助 either, bothexist, botheq tag
	validFn         map[string]CommonValidFn // 存放自定义的验证函数, 可以做到调用完就被清理
}

// structType
type structType struct {
	name       string            // 名字
	fieldInfos []structFieldInfo // 偏移量对应的字段信息内容
}

// structFieldInfo
type structFieldInfo struct {
	export     bool   // 是否可导出
	offset     int    // 偏移量
	name       string // 字段名
	validNames string // 验证规则
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
	// obj.errBuf.Grow(1 << 7)
	return obj
}

// free 释放
func (v *VStruct) free() {
	v.errBuf.Reset()
	v.ruleMap = nil
	v.valid2FieldsMap = nil
	v.validFn = nil
	syncValidStructPool.Put(v)
}

// SetRule 指定结构体设置验证规则, 不传则验证最外层的结构体
// obj 只支持一个参数, 多个无效, 此参数 待验证结构体
func (v *VStruct) SetRule(rule RM, obj ...interface{}) *VStruct {
	var ty reflect.Type
	l := len(obj)
	if l == 0 {
		ty = validOnlyOuterObj // 只验证最外层 struct
	} else if l == 1 {
		ty = RemoveTypePtr(reflect.TypeOf(obj[0]))
		if ty == timeReflectType {
			return v
		}
	} else {
		return v
	}

	if v.ruleMap == nil {
		v.ruleMap = make(map[reflect.Type]RM, 1)
	}
	v.ruleMap[ty] = rule
	return v
}

// getCusRule 根据字段名获取自定义验证规则
func (v *VStruct) getCusRule(ty reflect.Type) RM {
	if v.ruleMap == nil {
		return nil
	}
	return v.ruleMap[ty]
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
			val := reflectValue.Index(i)
			if i == 0 {
				structName = val.Type().String()
			}
			v.validate(structName+"["+ToStr(i)+"]", val, true)
		}
		return v.getError()
	case reflect.Map:
		iter := reflectValue.MapRange()
		for iter.Next() {
			v.validate("map["+ToStr(iter.Key())+"]", iter.Value(), true)
		}
		return v.getError()
	}
	return v.validate("", reflectValue, false).getError()
}

// SetValidFn 自定义设置验证函数
func (v *VStruct) SetValidFn(validName string, fn CommonValidFn) *VStruct {
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
// isValidGatherObj 是否验证集合对象, 包含: slice/array/map
func (v *VStruct) validate(structName string, value reflect.Value, isValidGatherObj ...bool) *VStruct {
	tv := RemoveValuePtr(value)
	ty := tv.Type()
	// fmt.Printf("ty: %v, structName: %q\n", ty, structName)
	// 如果不是结构体就退出
	if tv.Kind() != reflect.Struct {
		// 这里主要防止验证的切片为非结构体切片, 如 []int{1, 2, 3}, 这里会出现1, 为非指针所有需要退出
		if len(isValidGatherObj) > 0 && isValidGatherObj[0] {
			return v
		}
		v.errBuf.WriteString("src param \"" + structName + "." + ty.Name() + "\" is not struct" + ErrEndFlag)
		return v
	}

	cacheStructType := v.getCacheStructType(ty)
	totalFieldNum := len(cacheStructType.fieldInfos)
	var cusRM RM
	if structName == "" { // 只有最外层的结构体此值为空
		structName = cacheStructType.name
		// 在调用 SetRule 时没有设置验证对象时, 默认验证最外层结构体
		cusRM = v.getCusRule(validOnlyOuterObj)
		if len(cusRM) == 0 { // 设置了验证对象
			cusRM = v.getCusRule(ty)
		}
	} else { // 递归验证嵌套对象
		cusRM = v.getCusRule(ty)
	}
	// fmt.Printf("cusRM: %+v\n", cusRM)
	for fieldNum := 0; fieldNum < totalFieldNum; fieldNum++ {
		fieldInfo := cacheStructType.fieldInfos[fieldNum]
		// 判断下是否可导出
		if !fieldInfo.export {
			continue
		}

		// 如果设置了规则就覆盖 tag 中的验证内容
		if rule := cusRM.Get(fieldInfo.name); rule != "" {
			fieldInfo.validNames = rule
		}

		// fmt.Printf("name: %s, rule: %s\n", fieldInfo.name, fieldInfo.validNames)
		// 没有 validNames 直接跳过
		if fieldInfo.validNames == "" {
			continue
		}

		fieldValue := tv.Field(fieldInfo.offset)
		// 根据 tag 中的验证内容进行验证
		for _, validName := range ValidNamesSplit(fieldInfo.validNames) {
			if validName == "" {
				continue
			}

			validKey, _, cusMsg := ParseValidNameKV(validName)
			fn, err := v.getValidFn(validKey)
			if err != nil {
				v.errBuf.WriteString(GetJoinFieldErr(structName, fieldInfo.name, err))
				continue
			}

			// fmt.Printf("structName: %s, structFieldName: %s, tv: %v\n", cacheStructType.name, fieldInfo.name, fieldValue)
			// 开始验证
			// VStruct 内的验证方法
			if fn == nil {
				switch validKey {
				case Required:
					v.required(structName, fieldInfo.name, cusMsg, fieldValue)
				case Exist:
					v.exist(true, structName, fieldInfo.name, cusMsg, fieldValue)
				case Either, BothEq:
					v.initValid2FieldsMap(validName, structName, fieldInfo.name, cusMsg, fieldValue)
				}
				continue
			}

			// VStruct 外拓展的验证方法
			if fieldValue.IsZero() { // 空就直接跳过
				continue
			}
			fn(v.errBuf, validName, structName, fieldInfo.name, fieldValue)
		}
	}
	return v
}

// getCacheStructType 获取缓存中的 reflect.Type
func (v *VStruct) getCacheStructType(ty reflect.Type) structType {
	if obj, ok := cacheStructType.Load(ty); ok {
		return obj.(structType)
	}

	l := ty.NumField()
	obj := structType{name: ty.Name()}
	obj.fieldInfos = make([]structFieldInfo, l)
	for fieldNum := 0; fieldNum < l; fieldNum++ {
		fieldInfo := ty.Field(fieldNum)
		if fieldInfo.Type == timeReflectType {
			continue
		}
		info := structFieldInfo{
			export:     IsExported(fieldInfo.Name),
			offset:     fieldNum,
			name:       fieldInfo.Name,
			validNames: fieldInfo.Tag.Get(v.targetTag),
		}
		obj.fieldInfos[fieldNum] = info
	}
	cacheStructType.Store(ty, obj)
	return obj
}

// required 验证 required
func (v *VStruct) required(structName, fieldName, cusMsg string, tv reflect.Value) {
	ok := true
	// 如果集合类型先判断下长度
	switch tv.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
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
		if tv.Type() == timeReflectType {
			return
		}
		v.validate(structName+"."+fieldName, tv, false)
	case reflect.Slice, reflect.Array:
		for i := 0; i < tv.Len(); i++ {
			v.validate(structName+"."+fieldName+"["+ToStr(i)+"]", tv.Index(i), true)
		}
	case reflect.Map:
		iter := tv.MapRange()
		for iter.Next() {
			v.validate(structName+"."+fieldName+"["+ToStr(iter.Key())+"]", iter.Value(), true)
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
	v.valid2FieldsMap[validName] = append(v.valid2FieldsMap[validName], &name2Value{objName: structName, fieldName: fieldName, cusMsg: cusMsg, reflectVal: tv})
}

// either 判断两者不能都为空
func (v *VStruct) either(fieldInfos []*name2Value) {
	l := len(fieldInfos)
	if l == 1 { // 如果只有 1 个就没有必要向下执行了
		info := fieldInfos[0]
		v.errBuf.WriteString(GetJoinFieldErr(info.objName, info.fieldName, eitherValErr))
		return
	}
	isZeroLen := 0
	fieldInfoStr := "" // 拼接空的 structName, fliedName
	for _, fieldInfo := range fieldInfos {
		fieldInfoStr += "\"" + fieldInfo.objName + "." + fieldInfo.fieldName + "\", "
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
		v.errBuf.WriteString(GetJoinFieldErr(info.objName, info.fieldName, bothEqValErr))
		return
	}

	var (
		tmp          interface{}
		fieldInfoStr string // 拼接空的 structName, fliedName
		eq           = true
	)
	for i, fieldInfo := range fieldInfos {
		fieldInfoStr += "\"" + fieldInfo.objName + "." + fieldInfo.fieldName + "\", "
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
