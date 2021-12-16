package valid

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func ExampleRequired() {
	type Tmp struct {
		Name string `valid:"required"`
		Age  int32  `valid:"ge=0"`
	}
	v := &Tmp{Name: "", Age: 10}
	fmt.Println(ValidateStruct(v))

	// Output:
	// "Tmp.Name" is required
}

func ExampleTo() {
	type Tmp struct {
		Name string `valid:"to=1~3"`
		Age  int32  `valid:"to=0"`
		Addr string `valid:"le=3"`
	}
	v := &Tmp{Name: "测试调", Age: 100, Addr: "tets"}
	fmt.Println(ValidateStruct(v))
	// Output:
	// "Tmp.Addr" is length more than or equal 3
}

func ExampleOto() {
	type Tmp struct {
		Name     string `valid:"oto=1~3"`
		Age      int32  `valid:"oto=0~100"`
		NickName string `valid:"gt=1"`
		Addr     string `valid:"lt=3"`
	}
	v := &Tmp{Name: "测试", Age: 0, NickName: "h1", Addr: "tets"}
	fmt.Println(ValidateStruct(v))
	// Output:
	// "Tmp.Addr" is length more than 3
}

func ExampleEq() {
	type Tmp struct {
		Name  string  `valid:"required,eq=3"`
		Age   int32   `valid:"required,eq=20"`
		Score float64 `valid:"eq=80"`
		Phone string  `valid:"eq=11"`
	}
	v := &Tmp{
		Name:  "xue",
		Age:   20,
		Score: 80,
		Phone: "1354004261",
	}
	fmt.Println(ValidateStruct(v))
	// Output:
	// "Tmp.Phone" length should equal 11
}

func ExampleDate() {
	type Tmp struct {
		Date     string `valid:"date"`
		Datetime string `valid:"datetime"`
	}
	v := &Tmp{Date: "2021-10", Datetime: "2021-10-11"}
	fmt.Println(ValidateStruct(v))
	// Output:
	// "Tmp.Date" is not date, eg: 2021-11-15
	// "Tmp.Datetime" is not datetime, eg: 2021-11-15 23:59:59
}

func ExampleEither() {
	type Tmp struct {
		Either1 int32 `valid:"either=1"`
		Either2 int32 `valid:"either=1"`
	}
	v := &Tmp{}
	fmt.Println(ValidateStruct(v))
	// Output:
	// "Tmp"."Either1", "Tmp"."Either2" they shouldn't all be empty
}

func ExampleIn() {
	type Tmp struct {
		SelectNum int32  `valid:"in=(1/2/3/4)"`
		SelectStr string `valid:"in=(a/b/c/d)"`
	}
	v := &Tmp{SelectNum: 1, SelectStr: "ac"}
	fmt.Println(ValidateStruct(v))
	// Output:
	// "Tmp.SelectStr" should in (a/b/c/d)
}

func ExamplePhone() {
	type Tmp struct {
		Phone string `valid:"phone"`
	}
	v := &Tmp{Phone: "1"}
	fmt.Println(ValidateStruct(v))

	// Output:
	// "Tmp.Phone" is not phone
}

func ExampleEmail() {
	type Tmp struct {
		Email string `valid:"email"`
	}
	v := &Tmp{Email: "xuesongtao512@qq.com"}
	fmt.Println(ValidateStruct(v))

	// Output:
	// nil
}

func ExampleIdCard() {
	type Tmp struct {
		IDCard string `valid:"idcard"`
	}
	v := &Tmp{IDCard: "511321"}
	fmt.Println(ValidateStruct(v))

	// Output:
	// "Tmp.IDCard" is not idcard
}

func ExampleInt() {
	type Tmp struct {
		IntString string  `valid:"int"`
		IntNum    int     `valid:"int"`
		FloatNum  float32 `valid:"int"`
	}

	v := &Tmp{
		IntString: "11",
		IntNum:    1,
		FloatNum:  1.0,
	}
	fmt.Println(ValidateStruct(v))

	// Output:
	// "Tmp.FloatNum" is not integer
}

func ExampleFloat() {
	type Tmp struct {
		FloatString string  `valid:"float"`
		IntNum      int     `valid:"float"`
		FloatNum32  float32 `valid:"float"`
		FloatNum64  float64 `valid:"float"`
	}

	v := &Tmp{
		FloatString: "1.1",
		IntNum:      10,
		FloatNum32:  12.5,
		FloatNum64:  1.0,
	}
	fmt.Println(ValidateStruct(v))

	// Output:
	// "Tmp.IntNum" is not float
}

func ExampleSetCustomerValidFn() {
	type Tmp struct {
		Name string `valid:"required"`
		Age  string `valid:"num"`
	}

	isNumFn := func(errBuf *strings.Builder, validName, structName, filedName string, tv reflect.Value) {
		ok, _ := regexp.MatchString("^\\d+$", tv.String())
		if !ok {
			errBuf.WriteString(fmt.Sprintf("%q is not num", structName+"."+filedName))
			return
		}
	}

	SetCustomerValidFn("num", isNumFn)
	v := Tmp{Name: "12", Age: "1ha"}
	fmt.Println(ValidateStruct(&v))

	// Output:
	// "Tmp.Age" is not num
}

func ExampleRule() {
	type Tmp struct {
		Name string
		Age  string
		ClassName string
	}
	v := Tmp{Name: "xue", Age: "12a"}
	rule := RM{"name,age,classname": "required", "age": "int"}
	if err := ValidStructForRule(rule, &v); err != nil {
		fmt.Println(err)
	}

	// Output:
	// "Tmp.Age" is not integer; "Tmp.ClassName" is required
}


func ExampleRule2() {
	type Tmp struct {
		Name string
		Age  string
	}
	v := Tmp{Name: "xue"}
	ruleMap := NewRule()
	if v.Name == "xue" {
		ruleMap.Set("age", "required")
	}
	if err := ValidStructForRule(ruleMap, &v); err != nil {
		fmt.Println(err)
	}

	// Output:
	// "Tmp.Age" is required
}
