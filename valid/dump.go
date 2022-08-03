package valid

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

// dumpStruct
type dumpStruct struct {
	buf             *strings.Builder
	nullStructFiled reflect.StructField
	numBytes        [64]byte
}

// NewDumpStruct
func NewDumpStruct() *dumpStruct {
	return &dumpStruct{
		buf: &strings.Builder{},
	}
}

// HandleDumpStruct 处理待打印的 struct
func (d *dumpStruct) HandleDumpStruct(v reflect.Value, isSlice ...bool) *dumpStruct {
	tv := reflect.Indirect(v)
	if !tv.IsValid() {
		d.buf.WriteString("null")
		return d
	}
	ty := tv.Type()

	// 不是结构体就不处理
	if tv.Kind() != reflect.Struct {
		if len(isSlice) > 0 && isSlice[0] { // 非结构体切片
			d.loopHandleKV(d.nullStructFiled, tv, false)
		}
		return d
	}

	// 结构体
	d.buf.WriteByte('{')
	maxIndex := tv.NumField()
	if maxIndex == 0 {
		return d
	}

	// 需要添加第一个逗号
	needAddComma := false
	structField := ty.Field(0)
	if IsExported(structField.Name) {
		d.loopHandleKV(structField, tv.Field(0))
		needAddComma = true
	}
	for i := 1; i < maxIndex; i++ {
		structField = ty.Field(i)
		if !IsExported(structField.Name) {
			continue
		}
		// 第一次需要判断下
		if needAddComma {
			d.buf.WriteString(",")
		}
		needAddComma = true
		d.loopHandleKV(structField, tv.Field(i))
	}
	d.buf.WriteByte('}')
	return d
}

// Get 获取打印的结果
func (d *dumpStruct) Get() string {
	defer d.buf.Reset()
	return d.buf.String()
}

// loopHandleKV 循环处理
func (d *dumpStruct) loopHandleKV(s reflect.StructField, tv reflect.Value, isNeedFileName ...bool) {
	// fmt.Printf("s: %+v\n", s)
	// fmt.Printf("tv: %+v\n", filedValue)

	needFiledName := true
	if len(isNeedFileName) > 0 {
		needFiledName = isNeedFileName[0]
	}

	// 写入字段名
	if needFiledName {
		d.buf.WriteString("\"" + s.Name + "\":")
	}

	// 不处理的内容
	if s.Name == "Time" && s.Type == timeReflectType {
		d.buf.WriteString("\"time is not handle\"")
		return
	}

	switch tv.Kind() {
	case reflect.String: // 字符串
		d.buf.WriteString("\"" + tv.String() + "\"")
	case reflect.Bool:
		boolStr := "false"
		if tv.Bool() {
			boolStr = "true"
		}
		d.buf.WriteString("\"" + boolStr + "\"")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		d.buf.Write(strconv.AppendInt(d.numBytes[:0], tv.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		d.buf.Write(strconv.AppendUint(d.numBytes[:0], tv.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		d.buf.Write(strconv.AppendFloat(d.numBytes[:0], tv.Float(), 'f', -1, 64))
	case reflect.Ptr, reflect.Struct, reflect.Interface:
		d.HandleDumpStruct(tv)
	case reflect.Slice, reflect.Array: // 切片
		d.buf.WriteByte('[')
		sliceLen := tv.Len()
		for i := 0; i < sliceLen; i++ {
			d.HandleDumpStruct(tv.Index(i), true)
			if i < sliceLen-1 {
				d.buf.WriteString(",")
			}
		}
		d.buf.WriteByte(']')
	case reflect.Map: // map
		d.buf.WriteByte('{')
		mapObj := tv.MapRange()
		mapLen := tv.Len()
		tmpIndex := 0
		for mapObj.Next() {
			// 把 key 处理成字符串
			d.buf.WriteByte('"')
			d.loopHandleKV(d.nullStructFiled, mapObj.Key(), false)
			d.buf.WriteByte('"')
			d.buf.WriteString(":")
			d.loopHandleKV(d.nullStructFiled, mapObj.Value(), false)
			if tmpIndex < mapLen-1 {
				d.buf.WriteString(",")
			}
			tmpIndex++
		}
		d.buf.WriteByte('}')
	default: // 其他不处理, 如: func/chan
		d.buf.WriteString("\"unknown\"")
	}
}

// =========================== 常用方法进行封装 =======================================

// GetDumpStructStr 获取待 dump 的结构体字符串, 支持 json 格式化
// 只会把字段名解析成 key, 值为 value
// GetDumpStructStrForJson 性能较好, 只是在有 json_tag 的时候会用 json_tag 作为 key
func GetDumpStructStr(v interface{}) string {
	return NewDumpStruct().HandleDumpStruct(reflect.ValueOf(v)).Get()
}

// GetDumpStructStrForJson 先 json 序列化, 再获取
func GetDumpStructStrForJson(v interface{}) string {
	buf := new(strings.Builder)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(v)
	return strings.TrimSuffix(buf.String(), "\n")
}
