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
	targetTag string              // 结构体中的待指定的验证的 tag
	ruleMap   map[reflect.Type]RM // 验证规则, key: 为结构体 reflect.Type, value: 为该结构体的规则
	errBuf    *strings.Builder
	vc        *validCommon // 组合验证
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
	obj.errBuf = newStrBuf()
	obj.vc = &validCommon{}
	return obj
}

// free 释放
func (v *VStruct) free() {
	putStrBuf(v.errBuf)
	v.ruleMap = nil
	v.vc = nil
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
// 3. 支持map类型结构体验证
func (v *VStruct) Valid(src interface{}) error {
	if src == nil {
		return errors.New("src is nil")
	}

	reflectValue := RemoveValuePtr(reflect.ValueOf(src))
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
	v.vc.setValidFn(validName, fn)
	return v
}

// getValidFn 获取验证函数
func (v *VStruct) getValidFn(validName string) (CommonValidFn, error) {
	return v.vc.getValidFn(validName)
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
		v.errBuf.WriteString(GetJoinFieldErr(structName, ty.Name(), "is not struct"))
		return v
	}

	cacheStructType := v.getCacheStructType(ty)
	totalFieldNum := len(cacheStructType.fieldInfos)
	cusRM := v.getCusRule(ty)
	if structName == "" { // 只有最外层的结构体此值为空
		structName = cacheStructType.name
		// 在调用 SetRule 时没有设置验证对象时, 默认验证最外层结构体
		if len(cusRM) == 0 { // 设置了验证对象
			cusRM = v.getCusRule(validOnlyOuterObj)
		}
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
					v.vc.initValid2FieldsMap(&name2Value{
						validName:  validName,
						objName:    structName,
						fieldName:  fieldInfo.name,
						cusMsg:     cusMsg,
						reflectVal: fieldValue,
					})
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

// getError 获取 err
func (v *VStruct) getError() error {
	defer v.free()

	v.vc.valid(v.errBuf)
	if v.errBuf.Len() == 0 {
		return nil
	}

	// 这里需要去掉最后一个 ErrEndFlag
	return errors.New(strings.TrimSuffix(v.errBuf.String(), ErrEndFlag))
}
