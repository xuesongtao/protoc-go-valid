package valid

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// go test -timeout 30s -run ^Example gitee.com/xuesongtao/protoc-go-valid/valid -v -count=1

func ExampleRequired() {
	type Tmp struct {
		Name  string  `valid:"required|姓名必填"`
		Age   int32   `valid:"ge=0"`
		Hobby []int32 `valid:"required|爱好必填,le=3|爱好不能超过 3 个"`
	}
	v := &Tmp{Name: "", Age: 10, Hobby: []int32{1, 2, 3, 4}}
	fmt.Println(Struct(v))

	// Output:
	// "Tmp.Name" input "", 说明: 姓名必填; "Tmp.Hobby" input "4", 说明: 爱好不能超过 3 个
}

func ExampleExist() {
	type Man struct {
		Name string `valid:"required"`
		Age  int32  `valid:"le=100"`
	}

	type Student struct {
		M Man `valid:"exist"`
	}

	type Teather struct {
		M   *Man       `valid:"exist"`
		Stu []*Student `valid:"exist"`
	}
	teather := Teather{
		M: &Man{
			Name: "",
			Age:  90,
		},
		Stu: []*Student{{Man{Name: "test1", Age: 120}}},
	}
	fmt.Println(Struct(teather))

	// Output:
	// "Teather.M.Name" input "", explain: it is required; "Teather.Stu[0].M.Age" input "120", explain: it is more than 100 num-size
}

func ExampleTo() {
	type Tmp struct {
		Name string `valid:"to=1~3|姓名长度为 1-3 个字符"`
		Age  int32  `valid:"to=0~99|年龄应该在 0-99 之间"`
		Addr string `valid:"to=3~10"`
	}
	v := &Tmp{Name: "测试调1", Age: 100, Addr: "tets"}
	fmt.Println(Struct(v))

	// Output:
	// "Tmp.Name" input "测试调1", 说明: 姓名长度为 1-3 个字符; "Tmp.Age" input "100", 说明: 年龄应该在 0-99 之间
}

func ExampleGe() {
	type Tmp struct {
		Name string `valid:"ge=1"`
		Age  int32  `valid:"ge=0|应该大于 0"`
	}
	v := &Tmp{Name: "测试调", Age: -1}
	fmt.Println(Struct(v))

	// Output:
	// "Tmp.Age" input "-1", 说明: 应该大于 0
}

func ExampleLe() {
	type Tmp struct {
		Name string `valid:"le=2"`
		Age  int32  `valid:"le=0"`
	}
	v := &Tmp{Name: "测试调", Age: 1}
	fmt.Println(Struct(v))

	// Output:
	// "Tmp.Name" input "测试调", explain: it is more than 2 str-length; "Tmp.Age" input "1", explain: it is more than 0 num-size
}

func ExampleOTo() {
	type Tmp struct {
		Name     string `valid:"oto=1~3"`
		Age      int32  `valid:"oto=1~100"`
		NickName string `valid:"oto=0~10"`
		Addr     string `valid:"oto=1~3|家庭地址长度应该在 大于 1 且小于 3"`
	}
	v := &Tmp{Name: "测试", Age: 0, NickName: "h1", Addr: "tet"}
	fmt.Println(Struct(v))

	// Output:
	// "Tmp.Addr" input "tet", 说明: 家庭地址长度应该在 大于 1 且小于 3
}

func ExampleGt() {
	type Tmp struct {
		Name string `valid:"gt=2"`
		Age  int32  `valid:"gt=0"`
	}
	v := &Tmp{Name: "测试", Age: -1}
	fmt.Println(Struct(v))

	// Output:
	// "Tmp.Name" input "测试", explain: it is less than or equal 2 str-length; "Tmp.Age" input "-1", explain: it is less than or equal 0 num-size
}

func ExampleLt() {
	type Tmp struct {
		Name string `valid:"lt=2"`
		Age  int32  `valid:"lt=40"`
	}
	v := &Tmp{Name: "测试", Age: 99}
	fmt.Println(Struct(v))

	// Output:
	// "Tmp.Name" input "测试", explain: it is more than or equal 2 str-length; "Tmp.Age" input "99", explain: it is more than or equal 40 num-size
}

func ExampleEq() {
	type Tmp struct {
		Name  string  `valid:"required,eq=3"`
		Age   int32   `valid:"required,eq=20|年龄应该等于 20"`
		Score float64 `valid:"eq=80"`
		Phone string  `valid:"eq=11"`
	}
	v := &Tmp{
		Name:  "xue",
		Age:   21,
		Score: 80,
		Phone: "1354004261",
	}
	fmt.Println(Struct(v))

	// Output:
	// "Tmp.Age" input "21", 说明: 年龄应该等于 20; "Tmp.Phone" input "1354004261", explain: it should equal 11 str-length
}

func ExampleNoEq() {
	type Tmp struct {
		Name  string  `valid:"required,noeq=3"`
		Age   int32   `valid:"required,noeq=20|年龄不应该等于 20"`
		Score float64 `valid:"noeq=80"`
		Phone string  `valid:"noeq=11"`
	}
	v := &Tmp{
		Name:  "xue",
		Age:   20,
		Score: 80,
		Phone: "1354004261",
	}
	fmt.Println(Struct(v))

	// Output:
	// "Tmp.Name" input "xue", explain: it is not equal 3 str-length; "Tmp.Age" input "20", 说明: 年龄不应该等于 20; "Tmp.Score" input "80", explain: it is not equal 80 num-size
}

func ExampleDate() {
	type Tmp struct {
		Year       string `valid:"year"`
		Year2Month string `valid:"year2month=/"`
		Date       string `valid:"date=/"`
		Date2      string `valid:"date=''"`
		Datetime   string `valid:"datetime|应该为 xxxx-xx-xx xx:xx:xx 的时间格式"`
	}
	v := &Tmp{Year: "2001", Year2Month: "2000/01", Date: "2021/01/22", Date2: "20210122", Datetime: "2021-01-11 23:22"}
	fmt.Println(Struct(v))

	// Output:
	// "Tmp.Datetime" input "2021-01-11 23:22", 说明: 应该为 xxxx-xx-xx xx:xx:xx 的时间格式
}

func ExampleDatetime() {
	type Tmp struct {
		Datetime1 string `valid:"required,datetime|应该为 xxxx-xx-xx xx:xx:xx 的时间格式"`
		Datetime2 string `valid:"required,datetime='/'|应该为 xxxx/xx/xx xx:xx:xx 的时间格式"`
		Datetime3 string `valid:"required,datetime='/, ,/'|应该为 xxxx/xx/xx xx/xx/xx 的时间格式"`
		Datetime4 string `valid:"required,datetime=', ,'|应该为 xxxxxxxx xxxxxx 的时间格式"`
		Datetime5 string `valid:"required,datetime=',,'|应该为 xxxxxxxxxxxxxx 的时间格式"`
	}
	v := &Tmp{
		Datetime1: "2022/11-09 10:05:00",
		Datetime2: "2022/11/09 10:05:00",
		Datetime3: "2022/11/09 10/05/00",
		Datetime4: "20221109 100500",
		Datetime5: "20221109100500",
	}
	fmt.Println(Struct(v))

	// Output:
	// "Tmp.Datetime1" input "2022/11-09 10:05:00", 说明: 应该为 xxxx-xx-xx xx:xx:xx 的时间格式
}

func ExampleEither() {
	type Tmp struct {
		Either1 int32 `valid:"either=1"`
		Either2 int32 `valid:"either=1"`
	}
	v := &Tmp{}
	fmt.Println(Struct(v))

	// Output:
	// "Tmp.Either1", "Tmp.Either2" explain: they shouldn't all be empty
}

func ExampleBothEq() {
	type Tmp struct {
		BothEq1 int32 `valid:"botheq=1"`
		BothEq2 int32 `valid:"botheq=1"`
		BothEq3 int32 `valid:"botheq=1"`
	}
	v := &Tmp{
		BothEq1: 1,
		BothEq2: 1,
		BothEq3: 10,
	}
	fmt.Println(Struct(v))

	// Output:
	// "Tmp.BothEq1", "Tmp.BothEq2", "Tmp.BothEq3" explain: they should be equal
}

func ExampleIn() {
	type Tmp struct {
		SelectNum  int32  `valid:"in=(1/2/3/4)"`
		SelectStr  string `valid:"in=(a/b/c/d)|应该在 a/b/c/d 里选择"`
		SelectStr1 string `valid:"in=(a/'/d')|应该在 a/'/d' 里选择"` // 转义
	}
	v := &Tmp{SelectNum: 1, SelectStr: "ac", SelectStr1: "/d"}
	fmt.Println(Struct(v))

	// Output:
	// "Tmp.SelectStr" input "ac", 说明: 应该在 a/b/c/d 里选择
}

func ExampleInclude() {
	type Tmp struct {
		SelectStr string `valid:"include=(hello/test)"`
	}
	v := &Tmp{SelectStr: "hel"}
	fmt.Println(Struct(v))

	// Output:
	// "Tmp.SelectStr" input "hel", explain: it should include (hello/test)
}

func ExamplePhone() {
	type Tmp struct {
		Phone string `valid:"phone"`
	}
	v := &Tmp{Phone: "1"}
	fmt.Println(Struct(v))

	// Output:
	// "Tmp.Phone" input "1", explain: it is not phone
}

func ExampleEmail() {
	type Tmp struct {
		Email string `valid:"email"`
	}
	v := &Tmp{Email: "xuesongtao512qq.com"}
	fmt.Println(Struct(v))

	// Output:
	// "Tmp.Email" input "xuesongtao512qq.com", explain: it is not email
}

func ExampleIDCard() {
	type Tmp struct {
		IDCard string `valid:"idcard"`
	}
	v := &Tmp{IDCard: "511321"}
	fmt.Println(Struct(v))

	// Output:
	// "Tmp.IDCard" input "511321", explain: it is not idcard
}

func ExampleIp() {
	type Tmp struct {
		Ip   string `valid:"required,ip|ip不正确"`
		Ipv4 string `valid:"required,ipv4|ipv4不正确"`
		Ipv6 string `valid:"required,ipv6|ipv6不正确"`
	}

	v := &Tmp{
		Ip:   "192.168.1.0",
		Ipv4: "192.168.0.113",
		Ipv6: "2001:3CA1:10F:1A:121B:0:0:",
	}
	fmt.Println(Struct(&v))

	// OutPut:
	// "Tmp.Ipv6" input "2001:3CA1:10F:1A:121B:0:0:", 说明: ipv6不正确
}

func ExampleInt() {
	type Tmp struct {
		IntString string `valid:"int|请输入整数类"`
		IntNum    int    `valid:"int"`
	}

	v := &Tmp{
		IntString: "11.121",
		IntNum:    1,
	}
	fmt.Println(Struct(&v))

	// Output:
	// "Tmp.IntString" input "11.121", 说明: 请输入整数类
}

func ExampleInts() {
	type Tmp struct {
		IntStr      string   `valid:"ints"`
		IntStr1     string   `valid:"ints=-"`
		IntSlice    []int    `valid:"ints"`
		IntSliceStr []string `valid:"ints"`
	}

	v := &Tmp{
		IntStr:      "1,2,3",
		IntStr1:     "1,2,3",
		IntSliceStr: []string{"1", "hello"},
		IntSlice:    []int{1, 2},
	}
	fmt.Println(Struct(&v))

	// Output:
	// "Tmp.IntStr1" input "1,2,3", explain: it is not separated by "-" num; "Tmp.IntSliceStr" input "[1, hello]", explain: slice/array element is not all num
}

func ExampleFloat() {
	type Tmp struct {
		FloatString string  `valid:"float|请输入浮点数"`
		FloatNum32  float32 `valid:"float"`
		FloatNum64  float64 `valid:"float"`
	}

	v := &Tmp{
		FloatString: "1",
		FloatNum32:  12.5,
		FloatNum64:  1.0,
	}

	fmt.Println(Struct(v))

	// Output:
	// "Tmp.FloatString" input "1", 说明: 请输入浮点数
}

func ExampleRe() {
	type Tmp struct {
		Name string `valid:"required|必填,re='[a-z]+'|姓名必须为英文"`
		Age  string `valid:"re='\\d{2}'|年龄必须为 2 位数"`
		Addr string `valid:"required|地址必须,re='[\u4e00-\u9fa5]'|地址必须为中文"`
		Pwd  string `valid:"required,re='^[A-Za-z0-9]{8,16}$'|密码必须为8-16为字母和数字组合"`
	}

	v := &Tmp{
		Name: "测试",
		Age:  "1",
		Addr: "四川成都",
		Pwd:  "1234567",
	}

	fmt.Println(Struct(v))

	// Output:
	// "Tmp.Name" input "测试", 说明: 姓名必须为英文; "Tmp.Age" input "1", 说明: 年龄必须为 2 位数; "Tmp.Pwd" input "1234567", 说明: 密码必须为8-16为字母和数字组合
}

func ExampleUnique() {
	type Tmp struct {
		Name      string   `valid:"required|必填,re='[a-z]+'|姓名必须为英文"`
		Hobby     string   `valid:"required|必填,unique|爱好唯一"`
		ClassName []string `valid:"required|必填,unique|班级名唯一"`
		LikeNum   []int    `valid:"required|必填,unique|幸运数唯一"`
	}

	v := &Tmp{
		Name:      "测试",
		Hobby:     "打篮球,踢足球,打篮球",
		ClassName: []string{"美术班", "算术班", "足球班"},
		LikeNum:   []int{1, 2, 3, 1},
	}

	fmt.Println(Struct(v))

	// Output:
	// "Tmp.Name" input "测试", 说明: 姓名必须为英文; "Tmp.Hobby" input "打篮球,踢足球,打篮球", 说明: 爱好唯一; "Tmp.LikeNum" input "[1,2,3,1]", 说明: 幸运数唯一
}

func ExampleJson() {
	type Tmp struct {
		Name  string `valid:"required"`
		Json1 string `valid:"required,json|json格式不正确"`
		Json2 string `valid:"required,json"`
	}
	tmp := &Tmp{
		Name:  "测试json",
		Json1: `[{"id":1,"name":"test","age":10,"cls_name":"初一","addr":"四川成都"},{"id":2,"name":"test","age":10,"cls_name":"初二","addr":"四川成都"}]`,
		Json2: `{"name":"A-Tao","hobby":["play basketball","write golang code"],"class_name":"community university","addr":"ChengDu"}`,
	}
	err := Struct(tmp)
	fmt.Println(err)

	// Output:
	// <nil>
}

func ExamplePrefix() {
	type Tmp struct {
		Name string `valid:"required,prefix=test|前缀必须为test"`
	}
	tmp := &Tmp{
		Name: "ceshi",
	}
	err := Struct(tmp)
	fmt.Println(err)

	// Output:
	// "Tmp.Name" input "ceshi", 说明: 前缀必须为test
}

func ExampleSuffix() {
	type Tmp struct {
		Name string `valid:"required,suffix=test|后缀必须为test"`
	}
	tmp := &Tmp{
		Name: "mytest",
	}
	err := Struct(tmp)
	fmt.Println(err)

	// Output:
	// <nil>
}

func ExampleStructForFns() {
	type Tmp struct {
		Name string
		Addr string
	}

	isLower := func(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
		if err := CheckFieldIsStr(objName, fieldName, tv); err != nil {
			errBuf.WriteString(err.Error())
			return
		}
		if strings.ToLower(tv.String()) != tv.String() {
			errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), ExplainEn, "it is not lower"))
		}
	}

	isZh := func(errBuf *strings.Builder, validName, objName, fieldName string, tv reflect.Value) {
		if err := CheckFieldIsStr(objName, fieldName, tv); err != nil {
			errBuf.WriteString(err.Error())
			return
		}
		if match, _ := regexp.MatchString("^[\u4e00-\u9fa5]+$", tv.String()); !match {
			errBuf.WriteString(GetJoinValidErrStr(objName, fieldName, tv.String(), ExplainEn, "it should is zh"))
		}
	}
	fnMap := Name2FnMap{"islower": isLower, "iszh": isZh}
	rm := RM{"Name": "islower", "Addr": "iszh"}

	v := &Tmp{
		Name: "Test",
		Addr: "四川",
	}
	err := StructForFns(v, rm, fnMap)
	fmt.Println(err)

	// Output:
	// "Tmp.Name" input "Test", explain: it is not lower
}

func ExampleNestedStructForRule() {
	type Tmp1 struct {
		Name string
	}

	type Tmp struct {
		Ip string
		T  []Tmp1
	}
	rmap := map[interface{}]RM{
		// key 必须为 指针
		&Tmp{}:  NewRule().Set("Ip,T", Required).Set("Ip", GenValidKV(VIp, "", "ip 格式不正确")),
		&Tmp1{}: NewRule().Set("Name", GenValidKV(Required, "", "姓名必填")),
	}
	v := &Tmp{
		Ip: "256.12.22.400",
		T:  []Tmp1{{Name: ""}, {Name: "2"}},
	}
	err := NestedStructForRule(v, rmap)
	fmt.Println(err)

	// Output:
	// "Tmp.Ip" input "256.12.22.400", 说明: ip 格式不正确; "Tmp.T[0].Name" input "", 说明: 姓名必填
}

func ExampleJoinTag2Val() {
	val := JoinTag2Val(VIn, "1/2/3", "必须在 1,2,3 之中")
	fmt.Println(val)

	// Output:
	// in=(1/2/3)|必须在 1,2,3 之中
}

func ExampleValidStructForRule() {
	type Tmp struct {
		Name string
		Age  int
	}
	v := Tmp{Name: "xue", Age: 101}
	ruleObj := NewRule()
	if v.Name == "xue" {
		// "required|必填,le=100|年龄最大为 100"
		ruleObj.Set("Age", JoinTag2Val(Required, "", "必填"), JoinTag2Val(VLe, "100", "年龄最大为 100"))
	}
	if err := ValidStructForRule(ruleObj, &v); err != nil {
		fmt.Println(err)
	}

	// Output:
	// "Tmp.Age" input "101", 说明: 年龄最大为 100
}

func ExampleSetCustomerValidFn() {
	type Tmp struct {
		Name string `valid:"required"`
		Age  string `valid:"num"`
	}

	isNumFn := func(errBuf *strings.Builder, validName, structName, fieldName string, tv reflect.Value) {
		ok, _ := regexp.MatchString("^\\d+$", tv.String())
		if !ok {
			errBuf.WriteString(fmt.Sprintf("%q is not num", structName+"."+fieldName))
			return
		}
	}

	SetCustomerValidFn("num", isNumFn)
	v := Tmp{Name: "12", Age: "1ha"}
	fmt.Println(Struct(&v))

	// Output:
	// "Tmp.Age" is not num
}

func ExampleValidStructForMyValidFn() {
	type Tmp struct {
		Name string `valid:"required"`
		Age  string `valid:"num"`
	}

	isNumFn := func(errBuf *strings.Builder, validName, structName, fieldName string, tv reflect.Value) {
		ok, _ := regexp.MatchString("^\\d+$", tv.String())
		if !ok {
			errBuf.WriteString(fmt.Sprintf("%q is not num", structName+"."+fieldName))
			return
		}
	}

	v := Tmp{Name: "12", Age: "1ha"}

	fmt.Println(ValidStructForMyValidFn(v, "num", isNumFn))

	// Output:
	// "Tmp.Age" is not num
}

func ExampleVar() {
	err := Var(101, Required, GenValidKV(VTo, "1~100", "年龄1~100"))
	fmt.Println(err)

	// OutPut:
	// input "101", 说明: 年龄1~100
}

func ExampleUrl() {
	url := "http://test.com?name=test&age=10&nickname=test1"
	ruleObj := NewRule()
	ruleObj.Set("name", Required, GenValidKV(VTo, "5~10|姓名需在5-10之间"), GenValidKV(BothEq, "botheq=0"))
	ruleObj.Set("nickname", Required, GenValidKV(BothEq, "botheq=0"))
	err := Url(url, ruleObj)
	fmt.Println(err)

	// OutPut:
	// "name" input "test", 说明: 姓名需在5-10之间; "name", "nickname" explain: they should be equal
}

func ExampleGetOnlyExplainErr() {
	type Tmp struct {
		Name      string `valid:"to=2~4|姓名长度应该在2-4个长度"`
		IntString string `valid:"int|IntString请输入整数类"`
		IntNum    string `valid:"int|IntNum请输入整数"`
	}

	v := &Tmp{
		Name:      "xuesongtao",
		IntString: "11.121",
		IntNum:    "1a",
	}
	err := Struct(&v)
	fmt.Println(GetOnlyExplainErr(err.Error()))

	// OutPut:
	// 姓名长度应该在2-4个长度; IntString请输入整数类; IntNum请输入整数
}
