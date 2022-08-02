package valid

import (
	"fmt"
	"reflect"
	"testing"

	"gitee.com/xuesongtao/protoc-go-valid/test"
	"github.com/go-playground/validator/v10"
	"github.com/gookit/validate"
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
	type Users struct {
		Phone  string `validate:"required"`
		Passwd string `validate:"required,max=20,min=6"`
		Code   string `validate:"required,len=6"`
	}

	users := &Users{
		Phone:  "1326654487",
		Passwd: "123",
		Code:   "123456",
	}
	validate := validator.New()
	err := validate.Struct(users)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err) //Key: 'Users.Passwd' Error:Field validation for 'Passwd' failed on the 'min' tag
			return
		}
	}
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
	Name string `alipay:"required" validate:"required"`
}

func TestValidManyStruct(t *testing.T) {
	type Tmp struct {
		Ip string     `valid:"required,ipv4" validate:"required"`
		T  []TmpTest3 `valid:"required" validate:"required"`
	}

	v := &Tmp{
		// Ip: "61.240.17.210",
		Ip: "256.12.22.4",
	}
	datas := append([]*Tmp{}, v, v, v)
	sureMsg := `"*valid.Tmp-0.Tmp.Ip" input "256.12.22.4", explain: it is not ipv4; "*valid.Tmp-0.Tmp.T" input "", explain: it is required; "*valid.Tmp-1.Tmp.Ip" input "256.12.22.4", explain: it is not ipv4; "*valid.Tmp-1.Tmp.T" input "", explain: it is required; "*valid.Tmp-2.Tmp.Ip" input "256.12.22.4", explain: it is not ipv4; "*valid.Tmp-2.Tmp.T" input "", explain: it is required`
	err := ValidateStruct(datas)
	if !equal(err.Error(), sureMsg) {
		t.Error(noEqErr)
	}
}

func TestValidManyStruct2(t *testing.T) {
	type Tmp struct {
		Ip string     `valid:"required,ipv4" validate:"required"`
		T  []TmpTest3 `valid:"required" validate:"required"`
	}

	v := &Tmp{
		// Ip: "61.240.17.210",
		Ip: "256.12.22.4",
	}
	datas := append([]*Tmp{}, v, v, v)

	// 不支持
	validObj := validate.Struct(datas)
	validObj.Validate()
	for _, err := range validObj.Errors {
		t.Log(err)
	}
}

func TestValidManyStruct3(t *testing.T) {
	type Tmp struct {
		Ip string     `valid:"required,ipv4" validate:"required"`
		T  []TmpTest3 `valid:"required" validate:"required"`
	}

	v := &Tmp{
		// Ip: "61.240.17.210",
		// Ip: "256.12.22.4",
	}
	// datas := append([]*Tmp{}, v, v, v)

	// 不支持valid 多个
	validObj := validator.New()
	err := validObj.Struct(v)
	if err != nil {
		t.Log(err)
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

func TestValidateOrder(t *testing.T) {
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
		AppName:              "集美测试",
		TotalFeeFloat:        2,
		TestOrderDetailPtr:   testOrderDetailPtr,
		TestOrderDetailSlice: testOrderDetails,
	}
	// TODO: 验证单个就直接退出了
	validObj := validate.Struct(u)
	validObj.Validate()
	for _, err := range validObj.Errors {
		t.Log(err)
	}
}

func TestValidatorOrder(t *testing.T) {
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
		AppName:              "集美测试",
		TotalFeeFloat:        2,
		TestOrderDetailPtr:   testOrderDetailPtr,
		TestOrderDetailSlice: testOrderDetails,
	}
	validObj := validator.New()

	// TODO: 不支持验证切片结构体
	err := validObj.Struct(u)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(err)
}

func TestProtoPb1(t *testing.T) {
	u := &test.User{
		M: &test.Man{
			Name: "",
			Age:  0,
		},
		Phone: "13540042615",
	}
	err := ValidateStruct(u)
	if err == nil {
		return
	}

	suerMsg := `"User.Man.Name" input "", 说明: 姓名必填; "User.Man.Tmp valid: "he" is not exist, You can call SetValidFn`
	if !equal(err.Error(), suerMsg) {
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
