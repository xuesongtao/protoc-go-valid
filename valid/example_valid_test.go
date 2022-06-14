package valid

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// go test -run ^Example

func ExampleRequired() {
	type Tmp struct {
		Name  string  `valid:"required|姓名必填"`
		Age   int32   `valid:"ge=0"`
		Hobby []int32 `valid:"required|爱好必填,le=3|爱好不能超过 3 个"`
	}
	v := &Tmp{Name: "", Age: 10, Hobby: []int32{1, 2, 3, 4}}
	fmt.Println(ValidateStruct(v))

	// Output:
	// "Tmp.Name" input "" 说明: 姓名必填; "Tmp.Hobby" input "4" 说明: 爱好不能超过 3 个
}

func ExampleExist() {
	type Man struct {
		Name string `valid:"required"`
		Age  int32  `valid:"le=0"`
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
			Age:  0,
		},
		Stu: []*Student{{Man{Name: "test1", Age: 10}}},
	}
	fmt.Println(ValidateStruct(teather))

	// Output:
	// "Teather.Man.Name" input "" is required; "Teather-0.Student.Man.Age" input "10" numSize more than 0
}

func ExampleTo() {
	type Tmp struct {
		Name string `valid:"to=1~3|姓名长度为 1-3 个字符"`
		Age  int32  `valid:"to=0~99|年龄应该在 0-99 之间"`
		Addr string `valid:"to=3~10"`
	}
	v := &Tmp{Name: "测试调1", Age: 100, Addr: "tets"}
	fmt.Println(ValidateStruct(v))
	// Output:
	// "Tmp.Name" input "测试调1" 说明: 姓名长度为 1-3 个字符; "Tmp.Age" input "100" 说明: 年龄应该在 0-99 之间
}

func ExampleGe() {
	type Tmp struct {
		Name string `valid:"ge=1"`
		Age  int32  `valid:"ge=0|应该大于 0"`
	}
	v := &Tmp{Name: "测试调", Age: -1}
	fmt.Println(ValidateStruct(v))
	// Output:
	// "Tmp.Age" input "-1" 说明: 应该大于 0
}

func ExampleLe() {
	type Tmp struct {
		Name string `valid:"le=2"`
		Age  int32  `valid:"le=0"`
	}
	v := &Tmp{Name: "测试调", Age: 1}
	fmt.Println(ValidateStruct(v))
	// Output:
	// "Tmp.Name" input "测试调" strLength more than 2; "Tmp.Age" input "1" numSize more than 0
}

func ExampleOto() {
	type Tmp struct {
		Name     string `valid:"oto=1~3"`
		Age      int32  `valid:"oto=1~100"`
		NickName string `valid:"oto=0~10"`
		Addr     string `valid:"oto=1~3|家庭地址长度应该在 大于 1 且小于 3"`
	}
	v := &Tmp{Name: "测试", Age: 0, NickName: "h1", Addr: "tet"}
	fmt.Println(ValidateStruct(v))
	// Output:
	// "Tmp.Addr" input "tet" 说明: 家庭地址长度应该在 大于 1 且小于 3
}

func ExampleGt() {
	type Tmp struct {
		Name string `valid:"gt=2"`
		Age  int32  `valid:"gt=0"`
	}
	v := &Tmp{Name: "测试", Age: -1}
	fmt.Println(ValidateStruct(v))
	// Output:
	// "Tmp.Name" input "测试" strLength less than or equal 2; "Tmp.Age" input "-1" numSize less than or equal 0
}

func ExampleLt() {
	type Tmp struct {
		Name string `valid:"lt=2"`
		Age  int32  `valid:"lt=40"`
	}
	v := &Tmp{Name: "测试", Age: 99}
	fmt.Println(ValidateStruct(v))
	// Output:
	// "Tmp.Name" input "测试" strLength more than or equal 2; "Tmp.Age" input "99" numSize more than or equal 40
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
	fmt.Println(ValidateStruct(v))
	// Output:
	// "Tmp.Age" input "21" numSize 说明: 年龄应该等于 20; "Tmp.Phone" input "1354004261" strLength should equal 11
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
	fmt.Println(ValidateStruct(v))
	// Output:
	// "Tmp.Name" input "xue" strLength should no equal 3; "Tmp.Age" input "20" 说明: 年龄不应该等于 20; "Tmp.Score" input "80" numSize should no equal 80
}

func ExampleDate() {
	type Tmp struct {
		Year       string `valid:"year"`
		Year2Month string `valid:"year2month=/"`
		Date       string `valid:"date=/"`
		Datetime   string `valid:"datetime|应该为 xxxx-xx-xx xx:xx:xx 的时间格式"`
	}
	v := &Tmp{Year: "2001", Year2Month: "2000/01", Date: "2021/01/22", Datetime: "2021-01-11 23:22"}
	fmt.Println(ValidateStruct(v))
	// Output:
	// "Tmp.Datetime" input "2021-01-11 23:22" 说明: 应该为 xxxx-xx-xx xx:xx:xx 的时间格式
}

func ExampleEither() {
	type Tmp struct {
		Either1 int32 `valid:"either=1"`
		Either2 int32 `valid:"either=1"`
	}
	v := &Tmp{}
	fmt.Println(ValidateStruct(v))
	// Output:
	// "Tmp.Either1", "Tmp.Either2" they shouldn't all be empty
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
	fmt.Println(ValidateStruct(v))
	// Output:
	// "Tmp.BothEq1", "Tmp.BothEq2", "Tmp.BothEq3" they shouldn't is both equal
}

func ExampleIn() {
	type Tmp struct {
		SelectNum int32  `valid:"in=(1/2/3/4)"`
		SelectStr string `valid:"in=(a/b/c/d)|应该在 a/b/c/d 里选择"`
	}
	v := &Tmp{SelectNum: 1, SelectStr: "ac"}
	fmt.Println(ValidateStruct(v))
	// Output:
	// "Tmp.SelectStr" input "ac" 说明: 应该在 a/b/c/d 里选择
}

func ExampleInclude() {
	type Tmp struct {
		SelectStr string `valid:"include=(hello/test)|应该包含 hello/test"`
	}
	v := &Tmp{SelectStr: "hel"}
	fmt.Println(ValidateStruct(v))
	// Output:
	// "Tmp.SelectStr" input "hel" 说明: 应该包含 hello/test
}

func ExamplePhone() {
	type Tmp struct {
		Phone string `valid:"phone|不是正确的手机号码"`
	}
	v := &Tmp{Phone: "1"}
	fmt.Println(ValidateStruct(v))

	// Output:
	// "Tmp.Phone" input "1" 说明: 不是正确的手机号码
}

func ExampleEmail() {
	type Tmp struct {
		Email string `valid:"email"`
	}
	v := &Tmp{Email: "xuesongtao512qq.com"}
	fmt.Println(ValidateStruct(v))

	// Output:
	// "Tmp.Email" input "xuesongtao512qq.com" is not email
}

func ExampleIdCard() {
	type Tmp struct {
		IDCard string `valid:"idcard"`
	}
	v := &Tmp{IDCard: "511321"}
	fmt.Println(ValidateStruct(v))

	// Output:
	// "Tmp.IDCard" input "511321" is not idcard
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
	fmt.Println(ValidateStruct(&v))

	// Output:
	// "Tmp.IntString" input "11.121" 说明: 请输入整数类
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

	fmt.Println(ValidateStruct(v))

	// Output:
	// "Tmp.FloatString" input "1" 说明: 请输入浮点数
}

func ExampleValidStructForRule() {
	type Tmp struct {
		Name string
		Age  int
	}
	v := Tmp{Name: "xue", Age: 101}
	ruleObj := NewRule()
	if v.Name == "xue" {
		ruleObj.Set("Age", "required|必填,le=100|年龄最大为 100")
	}
	if err := ValidStructForRule(ruleObj, &v); err != nil {
		fmt.Println(err)
	}

	// Output:
	// "Tmp.Age" input "101" 说明: 年龄最大为 100
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

	// 弃用
	SetCustomerValidFn("num", isNumFn)
	v := Tmp{Name: "12", Age: "1ha"}
	fmt.Println(ValidateStruct(&v))

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
