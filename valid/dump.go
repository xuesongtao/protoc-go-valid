package valid

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

type dumpStruct struct {
	buf             *strings.Builder
	nullStructFiled reflect.StructField
}

func NewDumpStruct() *dumpStruct {
	return &dumpStruct{
		buf: &strings.Builder{},
	}
}

func (d *dumpStruct) HandleDumpStruct(v interface{}, isSlice ...bool) *dumpStruct {
	tv := removeValuePtr(reflect.ValueOf(v))
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

	d.buf.WriteByte('{')
	maxIndex := tv.NumField()
	for i := 0; i < maxIndex; i++ {
		d.loopHandleKV(ty.Field(i), tv.Field(i))

		// 去掉最后一个逗号
		if i < maxIndex-1 {
			d.buf.WriteString(", ")
		}
	}
	d.buf.WriteByte('}')
	return d
}

func (d *dumpStruct) Get() string {
	defer d.buf.Reset()
	return d.buf.String()
}

func (d *dumpStruct) loopHandleKV(s reflect.StructField, tv reflect.Value, isNeedFileName ...bool) {
	// fmt.Printf("s: %+v\n", s)
	// fmt.Printf("tv: %+v\n", filedValue)

	needFiledName := true
	if len(isNeedFileName) > 0 {
		needFiledName = isNeedFileName[0]
	}

	// 写入字段名
	if needFiledName {
		d.buf.WriteString("\"" + s.Name + "\": ")
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
		d.buf.Write(strconv.AppendInt([]byte{}, tv.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		d.buf.Write(strconv.AppendUint([]byte{}, tv.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		d.buf.Write(strconv.AppendFloat([]byte{}, tv.Float(), 'f', -1, 64))
	case reflect.Ptr: // 指针结构体
		d.HandleDumpStruct(tv.Interface())
	case reflect.Struct: // 结构体
		d.HandleDumpStruct(tv.Interface())
	case reflect.Slice, reflect.Array: // 切片
		d.buf.WriteByte('[')
		sliceLen := tv.Len()
		for i := 0; i < sliceLen; i++ {
			d.HandleDumpStruct(tv.Index(i).Interface(), true)
			if i < sliceLen-1 {
				d.buf.WriteString(", ")
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
			d.buf.WriteString(": ")
			d.loopHandleKV(d.nullStructFiled, mapObj.Value(), false)
			if tmpIndex < mapLen-1 {
				d.buf.WriteString(", ")
			}
			tmpIndex++
		}
		d.buf.WriteByte('}')
	default: // 其他不处理, 如: func/chan/interface
		d.buf.WriteString("unknown")
	}
}

// =========================== 常用方法进行封装 =======================================

// GetDumpStructStr 获取待 dump 的结构体字符串
func GetDumpStructStr(v interface{}) string {
	return NewDumpStruct().HandleDumpStruct(v).Get()
}

func GetDumpStructStrForJson(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
