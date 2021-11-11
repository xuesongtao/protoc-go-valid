package valid

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	syncValidPool       = sync.Pool{New: func() interface{} { return new(VStruct) }}
	defaultTargetTag    = "valid" // 默认的验证 tag
	vStructToTagErr     = errors.New("tag \"to\" is not ok, eg: to=1/to=6~30")
	vStructEitherTagErr = errors.New("tag \"either\" is not ok, eg: " +
		"type Test struct {\n" +
		"    OrderNo string `either=1`\n" +
		"    TradeNo sting `either=1`\n" +
		"}, errMsg: \"OrderNo\" either \"TradeNo\" they shouldn't all be empty")
)

type VStruct struct {
	targetTag string // 结构体中的 tag name
	endFlag   string // 用于分割 err
	errBuf    *strings.Builder
	existMap  map[int][]*name2Value // 已存在的, 用于 either tag
}

// name2Value
type name2Value struct {
	structName string
	filedName  string
	val        reflect.Value
}

// VStruct 验证结构体, 默认目标 tagName 为 defaultTargetTag, 现在支持的 tag 的值 如下:
// 1. required: 必填标识
// 2. to: 长度验证, 格式为 to=xxx,xxx(字段类型: 字符串则为长度, 字段类型: 数字型则为大小)
// 3. either: 多选一, 即多个中必须有一个必填, 格式为 either=xxx(通过数据进行标识)
// 4. phone: 手机号验证
// 说明: 如果想定义必填可以设置为: required, 如果为区间类型可以设置为: to=1~10
func NewVStruct(targetTag ...string) *VStruct {
	obj := syncValidPool.Get().(*VStruct)
	tagName := defaultTargetTag
	if len(targetTag) > 0 {
		tagName = targetTag[0]
	}
	obj.targetTag = tagName
	obj.endFlag = "\n "
	obj.errBuf = new(strings.Builder)
	return obj
}

// parseTag 解析 tag, 获取目标内容
func (v *VStruct) parseTag(dest, src string) string {
	complie := regexp.MustCompile(dest + "=(.*)")
	resSlice := complie.FindStringSubmatch(src)
	if len(resSlice) < 1 {
		return ""
	}
	// fmt.Println("resSlice: ", resSlice)
	// 因为 tag 是按逗号隔开, 这里通过逗号先分割直接取下标为 0 的就是 size tag
	resSlice = strings.Split(resSlice[1], ",")
	return resSlice[0]
}

// parseTagTo 解析 tag: to 中 min, max
func (v *VStruct) parseTagTo(toStr string) (min int, max int, err error) {
	// 通过分割符来判断是否为区间
	toSlice := strings.Split(toStr, "~")
	l := len(toSlice)
	// fmt.Println("toSlice: ", toSlice)
	switch l {
	case 1:
		min, err = strconv.Atoi(toSlice[0])
	case 2:
		if min, err = strconv.Atoi(toSlice[0]); err != nil {
			return
		}

		if max, err = strconv.Atoi(toSlice[1]); err != nil {
			return
		}
	default:
		err = vStructToTagErr
	}
	return
}

// Validate 验证执行体
func (v *VStruct) Validate(in interface{}) *VStruct {
	ry := reflect.TypeOf(in)
	if ry.Kind() != reflect.Ptr {
		v.errBuf.WriteString("structName: " + ry.Name() + " must ptr" + v.endFlag)
		return v
	}

	ry = ry.Elem()
	// 如果不是结构体就退出
	if ry.Kind() != reflect.Struct {
		v.errBuf.WriteString("in params \"" + ry.Name() + "\" is not struct")
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

		if strings.Index(tag, "required") > -1 { // 验证必填
			v.validRequired(ry.Name(), ty.Name, tv, ty.Type)
		} else if strings.Index(tag, "to") > -1 { // 验证长度
			v.validTo(tag, ry.Name(), ty.Name, tv)
		} else if strings.Index(tag, "either") > -1 { // 验证二选一
			v.initEither(tag, ry.Name(), ty.Name, tv)
		} else if strings.Index(tag, "phone") > -1 { // 验证手机号
			v.validPhone(tag, ry.Name(), ty.Name, tv)
		} else {
			v.errBuf.WriteString("\"" + ry.Name() + "." + ty.Name +  "\" tag \"" + tag + "\" no have, I am sorry")
		}
	}
	return v
}

// validRequired 验证 required
func (v *VStruct) validRequired(structName, filedName string, tv reflect.Value, ty reflect.Type) {
	if tv.IsZero() { // 验证必填
		// 生成如: "TestOrderDetailSlice.Price" is required
		v.errBuf.WriteString("\"" + structName + "." + filedName + "\" is required" + v.endFlag)
	} else if ty.Kind() == reflect.Ptr || (ty.Kind() == reflect.Struct && ty.Name() != "Time") { // 结构体
		if ty.Name() == structName { // 防止出现死循环
			return
		}
		v.Validate(tv.Interface())
	} else if ty.Kind() == reflect.Slice { // 切片
		for i := 0; i < tv.Len(); i++ {
			v.Validate(tv.Index(i).Interface())
		}
	}
}

// validTo 验证 to
func (v *VStruct) validTo(tag, structName, filedName string, tv reflect.Value) {
	min, max, err := v.parseTagTo(v.parseTag("to", tag))
	if err != nil {
		v.errBuf.WriteString(err.Error() + v.endFlag)
		return
	}

	var (
		unitStr                = "size" // 大小的单位
		isLessThan, isMoreThan bool
	)
	// fmt.Printf("min: %v, max: %v\n", min, max)
	switch tv.Kind() {
	case reflect.String:
		unitStr = "len"
		inLen := len([]rune(tv.String()))
		if min > 0 && inLen < min {
			isLessThan = true
		}

		if max > 0 && inLen > max {
			isMoreThan = true
		}
	case reflect.Float32, reflect.Float64:
		val := tv.Float()
		if min > 0 && val < float64(min) {
			isLessThan = true
		}

		if max > 0 && val > float64(max) {
			isMoreThan = true
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := tv.Int()
		if min > 0 && val < int64(min) {
			isLessThan = true
		}

		if max > 0 && val > int64(max) {
			isMoreThan = true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val := tv.Uint()
		if min > 0 && val < uint64(min) {
			isLessThan = true
		}

		if max > 0 && val > uint64(max) {
			isMoreThan = true
		}
	}

	if isLessThan {
		// 生成如: "TestOrder.AppName" is len less than 2
		v.errBuf.WriteString("\"" + structName + "." + filedName + "\" is " + unitStr + " less than " + fmt.Sprintf("%d", min) + v.endFlag)
	}

	if isMoreThan {
		// 生成如: "TestOrder.AppName" is len more than 30
		v.errBuf.WriteString("\"" + structName + "." + filedName + "\" is " + unitStr + " more than " + fmt.Sprintf("%d", max) + v.endFlag)
	}
}

// initEither 为验证 either 进行准备
func (v *VStruct) initEither(tag, structName, filedName string, tv reflect.Value) {
	if v.existMap == nil {
		v.existMap = make(map[int][]*name2Value, 5)
	}

	eitherNum := v.parseTag("either", tag)
	if eitherNum == "" {
		v.errBuf.WriteString(vStructEitherTagErr.Error() + v.endFlag)
		return
	}
	num, _ := strconv.Atoi(eitherNum)
	if _, ok := v.existMap[num]; !ok {
		v.existMap[num] = make([]*name2Value, 0, 5)
	}
	v.existMap[num] = append(v.existMap[num], &name2Value{structName: structName, filedName: filedName, val: tv})
}

// validEither 验证 either
func (v *VStruct) validEither() {
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

// validPhone 验证手机号
func (v *VStruct) validPhone(tag, structName, filedName string, tv reflect.Value) {
	matched, _ := regexp.MatchString("^1[3,4,5,7,8,9]\\d{9}$", tv.String())
	if !matched {
		v.errBuf.WriteString(structName + "." + filedName + " is not phone")
	}
}

func (v *VStruct) GetErrMsg() error {
	if v.errBuf.Len() == 0 {
		return nil
	}
	defer v.free()

	// 验证下 either
	v.validEither()

	// 这里需要去掉最后一个 endFlag
	return errors.New(strings.TrimRight(v.errBuf.String(), v.endFlag))
}

// free 释放
func (v *VStruct) free() {
	v.targetTag = ""
	v.endFlag = ""
	v.errBuf.Reset()
	v.existMap = nil
	syncValidPool.Put(v)
}

// =========================== 常用方法进行封装 =======================================

// ValidateStruct 验证结构体
func ValidateStruct(in interface{}, targetTag ...string) error {
	return NewVStruct(targetTag...).Validate(in).GetErrMsg()
}
