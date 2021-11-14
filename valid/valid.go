package valid

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// vStruct 验证结构体
type vStruct struct {
	isCheckStructSlice bool   // 标记是否已验证是否为结构体切片
	targetTag          string // 结构体中的待指定的验证的 tag
	endFlag            string // 用于分割 err
	errBuf             *strings.Builder
	existMap           map[int][]*name2Value // 已存在的, 用于 either tag
}

// name2Value
type name2Value struct {
	structName string
	filedName  string
	val        reflect.Value
}

// NewVStruct 验证结构体, 默认目标 tagName 为 "valid"
func NewVStruct(targetTag ...string) *vStruct {
	obj := syncValidPool.Get().(*vStruct)
	tagName := defaultTargetTag
	if len(targetTag) > 0 {
		tagName = targetTag[0]
	}
	obj.targetTag = tagName
	obj.endFlag = errEndFlag
	obj.errBuf = new(strings.Builder)
	return obj
}

// Validate 验证执行体
func (v *vStruct) Validate(structName string, in interface{}, isValidSlice ...bool) *vStruct {
	// 辅助 errMsg, 用于嵌套时拼接上一级的结构体名
	if structName != "" {
		structName = structName + "."
	}

	ry := reflect.TypeOf(in)
	if ry.Kind() != reflect.Ptr {
		// 这里主要防止验证的切片为非结构体切片
		if len(isValidSlice) > 0 && isValidSlice[0] {
			return v
		}
		v.errBuf.WriteString("structName: " + structName + ry.Name() + " must ptr" + v.endFlag)
		return v
	}

	ry = ry.Elem()
	// 如果不是结构体就退出
	if ry.Kind() != reflect.Struct {
		// 这里主要防止验证的切片为非结构体切片
		if len(isValidSlice) > 0 && isValidSlice[0] {
			return v
		}
		v.errBuf.WriteString("in params \"" + structName + ry.Name() + "\" is not struct")
		if ry.Name() == "" {
			v.errBuf.WriteString(", params should is *struct")
		}
		v.errBuf.WriteString(v.endFlag)
		return v
	}

	// 取值
	rv := reflect.ValueOf(in).Elem()
	for filedNum := 0; filedNum < rv.NumField(); filedNum++ {
		tv := rv.Field(filedNum)
		// 不能导出就跳过
		if !tv.CanInterface() {
			continue
		}

		ty := ry.Field(filedNum)
		tag := ty.Tag.Get(v.targetTag)
		// 没有 tag 直接跳过
		if tag == "" {
			continue
		}

		// 根据 tag 中的验证内容进行验证
		for _, validName := range strings.Split(tag, ",") {
			validKey, _ := ParseValidNameKV(validName)
			fn, err := GetValidFn(validKey)
			if err != nil {
				v.errBuf.WriteString(err.Error())
				continue
			}

			// 开始验证
			if fn == nil && validKey == "required" { // 必填
				v.required(structName+ry.Name(), ty.Name, tv)
			} else if fn == nil && validKey == "either" { // 多选一
				v.initEither(tag, structName+ry.Name(), ty.Name, tv)
			} else {
				fn(v.errBuf, validName, structName+ry.Name(), ty.Name, tv)
			}
		}
	}
	return v
}

// required 验证 required
func (v *vStruct) required(structName, filedName string, tv reflect.Value) {
	if tv.IsZero() { // 验证必填
		// 生成如: "TestOrderDetailSlice.Price" is required
		v.errBuf.WriteString("\"" + structName + "." + filedName + "\" is required" + v.endFlag)
	} else if tv.Kind() == reflect.Ptr || (tv.Kind() == reflect.Struct && structName != "Time") { // 结构体
		v.Validate(structName, tv.Interface())
	} else if tv.Kind() == reflect.Slice { // 切片
		for i := 0; i < tv.Len(); i++ {
			v.Validate(fmt.Sprintf("%s-%d", structName, i), tv.Index(i).Interface(), true)
		}
	}
}

// initEither 为验证 either 进行准备
func (v *vStruct) initEither(validName, structName, filedName string, tv reflect.Value) {
	_, eitherNum := ParseValidNameKV(validName)
	if eitherNum == "" {
		v.errBuf.WriteString(eitherValErr.Error() + v.endFlag)
		return
	}

	if v.existMap == nil {
		v.existMap = make(map[int][]*name2Value, 5)
	}

	num, _ := strconv.Atoi(eitherNum)
	if _, ok := v.existMap[num]; !ok {
		v.existMap[num] = make([]*name2Value, 0, 2)
	}
	v.existMap[num] = append(v.existMap[num], &name2Value{structName: structName, filedName: filedName, val: tv})
}

// validEither 验证 either
func (v *vStruct) either() {
	// 判断下是否有值, 有就说明有 either 验证
	if len(v.existMap) == 0 {
		return
	}
	for _, objs := range v.existMap {
		l := len(objs)
		if l <= 1 { // 如果只有 1 个就没有必要向下执行了
			continue
		}
		isZeroLen := 0
		zeroInfoStr := "" // 拼接空的 structName, fliedName
		for _, obj := range objs {
			if obj.val.IsZero() {
				isZeroLen++
				zeroInfoStr += "\"" + obj.structName + "\"." + "\"" + obj.filedName + "\", "
			}
		}

		// 判断下是否全部为空
		if l == isZeroLen {
			zeroInfoStr = strings.TrimRight(zeroInfoStr, ", ")
			v.errBuf.WriteString(zeroInfoStr + " they shouldn't all be empty" + v.endFlag)
		}
	}
}

// GetErrMsg 获取 err
func (v *vStruct) GetErrMsg() error {
	defer v.free()

	// 验证下 either
	v.either()

	if v.errBuf.Len() == 0 {
		return nil
	}

	// 这里需要去掉最后一个 endFlag
	return errors.New(strings.TrimRight(v.errBuf.String(), v.endFlag))
}

// free 释放
func (v *vStruct) free() {
	v.isCheckStructSlice = false
	v.targetTag = ""
	v.endFlag = ""
	v.errBuf.Reset()
	syncValidPool.Put(v)
}

// =========================== 常用方法进行封装 =======================================

// ValidateStruct 验证结构体
func ValidateStruct(in interface{}, targetTag ...string) error {
	return NewVStruct(targetTag...).Validate("", in).GetErrMsg()
}
