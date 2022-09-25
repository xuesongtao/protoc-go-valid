package valid

import (
	"fmt"
	"reflect"
	"testing"
)

const (
	noEqErr = "src, dest is not eq"
)

func equal(dest, src interface{}) bool {
	ok := reflect.DeepEqual(dest, src)
	if !ok {
		fmt.Printf("dest: %v\n", dest)
		fmt.Printf("src: %v\n", src)
	}
	return ok
}

func TestTmp(t *testing.T) {
	type Tmp struct {
		Name string
		Age  int
	}
	data := &Tmp{
		Name: "1",
		Age:  2,
	}
	ruleObj := NewRule()
	ruleObj.Set("Name,Age", Required)
	err := StructForFn(&data, ruleObj)
	if err != nil {
		t.Fatal(err)
	}

	data = &Tmp{}
	err = Struct(&data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(IncludeZhRe.MatchString("测试1hell"))
}

type TestOrder struct {
	AppName string `alipay:"to=5~10" validate:"min=5,max=10"` // 应用名
	// TotalFeeFloat        float64                 `alipay:"to=2~5" validate:"min:2|max:5"` // 订单总金额，单位为分，详见支付金额
	TotalFeeFloat        float64                 `alipay:"to=2~5" validate:"min=2,max=5"` // 订单总金额，单位为分，详见支付金额
	TestOrderDetailPtr   *TestOrderDetailPtr     `alipay:"required" validate:"required"`  // 商品详细描述
	TestOrderDetailSlice []*TestOrderDetailSlice `alipay:"required" validate:"required"`  // 商品详细描述
}

type TestOrderDetailPtr struct {
	TmpTest3 *TmpTest3 `alipay:"required" validate:"required"`
	// GoodsName string    `alipay:"to=1~2" validate:"minLen:1|maxLen:2"`
	GoodsName string `alipay:"to=1~2" validate:"min=1,max=2"`
}

type TestOrderDetailSlice struct {
	TmpTest3   *TmpTest3 `alipay:"required" validate:"required"`
	GoodsName  string    `alipay:"required" validate:"required"`
	BuyerNames []string  `alipay:"required" validate:"required"`
}

type TmpTest3 struct {
	Name string `alipay:"required" valid:"required" validate:"required"`
}

func TestValidManyStruct(t *testing.T) {
	type Tmp struct {
		Ip string     `valid:"required,ipv4"`
		T  []TmpTest3 `valid:"required"`
	}

	v := &Tmp{
		Ip: "256.12.22.4",
		T:  []TmpTest3{{Name: ""}},
	}
	datas := append([]*Tmp{}, v, v)
	sureMsg := `"*valid.Tmp-0.Tmp.Ip" input "256.12.22.4", explain: it is not ipv4; "*valid.Tmp-0.Tmp-0.TmpTest3.Name" input "", explain: it is required; "*valid.Tmp-1.Tmp.Ip" input "256.12.22.4", explain: it is not ipv4; "*valid.Tmp-1.Tmp-0.TmpTest3.Name" input "", explain: it is required`
	err := ValidateStruct(datas)
	if !equal(err.Error(), sureMsg) {
		t.Error(noEqErr)
	}
}

func TestValidateForValid(t *testing.T) {
	type Users struct {
		Phone  string `valid:"required"`
		Passwd string `valid:"required,to=6~20"`
		Code   string `valid:"required,eq=6"`
	}

	users := &Users{
		Phone:  "1326654487",
		Passwd: "123",
		Code:   "123456",
	}
	err := Struct(users)
	sureMsg := `"Users.Passwd" input "123", explain: it is less than 6 str-length`
	if !equal(err.Error(), sureMsg) {
		t.Error(noEqErr)
	}
}

func TestValidOrder(t *testing.T) {
	testOrderDetailPtr := &TestOrderDetailPtr{
		TmpTest3:  &TmpTest3{Name: "测试"},
		GoodsName: "玻尿酸",
	}
	// testOrderDetailPtr = nil

	testOrderDetails := []*TestOrderDetailSlice{
		{TmpTest3: &TmpTest3{Name: "测试1"}, BuyerNames: []string{"test1", "hello2"}},
		{TmpTest3: &TmpTest3{Name: "测试2"}, GoodsName: "隆鼻"},
		{GoodsName: "丰胸"},
		{TmpTest3: &TmpTest3{Name: "测试4"}, GoodsName: "隆鼻"},
	}
	// testOrderDetails = nil

	u := &TestOrder{
		AppName:              "测试",
		TotalFeeFloat:        2,
		TestOrderDetailPtr:   testOrderDetailPtr,
		TestOrderDetailSlice: testOrderDetails,
	}
	err := ValidateStruct(u, "alipay")
	if err == nil {
		return
	}
	sureMsg := `"TestOrder.AppName" input "测试", explain: it is less than 5 str-length; "TestOrder.TestOrderDetailPtr.GoodsName" input "玻尿酸", explain: it is more than 2 str-length; "TestOrder-0.TestOrderDetailSlice.GoodsName" input "", explain: it is required; "TestOrder-1.TestOrderDetailSlice.BuyerNames" input "", explain: it is required; "TestOrder-2.TestOrderDetailSlice.TmpTest3" input "", explain: it is required; "TestOrder-2.TestOrderDetailSlice.BuyerNames" input "", explain: it is required; "TestOrder-3.TestOrderDetailSlice.BuyerNames" input "", explain: it is required`
	if !equal(err.Error(), sureMsg) {
		t.Error(noEqErr)
	}
}

func TestGetJoinValidErrStr(t *testing.T) {
	t.Skip("GetJoinValidErrStr")
	res := GetJoinValidErrStr("User", "Name", "xue", "len is less than 3")
	if !equal(res, `"User.Name" input "xue" len is less than 3;`) {
		t.Error(noEqErr)
	}
}

func TestValidVar(t *testing.T) {
	t.Run("required", func(t *testing.T) {
		err := Var("hello world", Required)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("no support", func(t *testing.T) {
		err := Var("hello world", Exist, Either, BothEq)
		sureMsg := `valid "exist" is no support; valid "either" is no support; valid "botheq" is no support`
		if !equal(err.Error(), sureMsg) {
			t.Error(noEqErr)
		}
	})

	t.Run("to", func(t *testing.T) {
		err := Var(101, Required, GenValidKV(VTo, "1~100", "年龄1~100"))
		sureMsg := `input "101", 说明: 年龄1~100`
		if !equal(err.Error(), sureMsg) {
			t.Error(noEqErr)
		}
	})

	t.Run("in", func(t *testing.T) {
		err := Var("12", Required, GenValidKV(VIn, "11/2/3"))
		sureMsg := `input "12", explain: it should in (11/2/3)`
		if !equal(err.Error(), sureMsg) {
			t.Error(noEqErr)
		}

		err = Var("device", Required, GenValidKV(VInclude, "test/devia"))
		sureMsg = `input "device", explain: it should include (test/devia)`
		if !equal(err.Error(), sureMsg) {
			t.Error(noEqErr)
		}
	})

	t.Run("phone", func(t *testing.T) {
		err := Var("135400426170", Required, VPhone)
		sureMsg := `input "135400426170", explain: it is not phone`
		if !equal(err.Error(), sureMsg) {
			t.Error(noEqErr)
		}

	})

	t.Run("unique", func(t *testing.T) {
		err := Var([]string{"test", "test1", "test"}, Required, VUnique)
		sureMsg := `input "[test,test1,test]", explain: they're not unique`
		if !equal(err.Error(), sureMsg) {
			t.Error(noEqErr)
		}
	})
}

func TestValidUrl(t *testing.T) {
	t.Run("required", func(t *testing.T) {
		url := "http://test.com?name=test&age=10"
		ruleObj := NewRule()
		ruleObj.Set("name", Required, GenValidKV(VTo, "5~10|姓名需在5-10之间"))
		err := Url(url, ruleObj)
		sureMsg := `"name" input "test", 说明: 姓名需在5-10之间`
		if !equal(err.Error(), sureMsg) {
			t.Error(noEqErr)
		}
	})

	t.Run("no support", func(t *testing.T) {
		ruleObj := NewRule()
		ruleObj.Set("name", Exist)
		err := Url("http%3A%2F%2Ftest.com%3Fname%3Dtest%26age%3D10", ruleObj)
		sureMsg := `valid "exist" is no support`
		if !equal(err.Error(), sureMsg) {
			t.Error(noEqErr)
		}
	})

	t.Run("botheq", func(t *testing.T) {
		url := "http://test.com?name=test&age=10&nickname=test1"
		ruleObj := NewRule()
		ruleObj.Set("name", Required, GenValidKV(VTo, "5~10|姓名需在5-10之间"), GenValidKV(BothEq, "botheq=0"))
		ruleObj.Set("nickname", Required, GenValidKV(BothEq, "botheq=0"))
		err := Url(url, ruleObj)
		sureMsg := `"name" input "test", 说明: 姓名需在5-10之间; "name", "nickname" explain: they should be equal`
		if !equal(err.Error(), sureMsg) {
			t.Error(noEqErr)
		}
	})

	t.Run("in", func(t *testing.T) {
		// http://test.com?type=(1/2/3)
		url := "http%3A%2F%2Ftest.com%3Ftype%3D(1%2F2%2F3)"
		ruleObj := NewRule()
		ruleObj.Set("type", Required, GenValidKV(VIn, "1/2/3"))
		err := Url(url, ruleObj)
		sureMsg := `"type" input "(1/2/3)", explain: it should in (1/2/3)`
		if !equal(err.Error(), sureMsg) {
			t.Error(noEqErr)
		}
	})

	t.Run("phone", func(t *testing.T) {
		// http://test.com?phone=13540042619
		url := "http%3A%2F%2Ftest.com%3Fphone%3D135400426199"
		ruleObj := NewRule()
		ruleObj.Set("phone", Required, GenValidKV(VPhone, "", "不是手机号"))
		err := Url(url, ruleObj)
		sureMsg := `"phone" input "135400426199", 说明: 不是手机号`
		if !equal(err.Error(), sureMsg) {
			t.Error(noEqErr)
		}
	})
}
