package valid

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	// "github.com/gookit/validate"
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

type TestOrder struct {
	AppName              string                  `alipay:"to=2~10" validate:"minLen:2|maxLen:10"` // 应用名
	TotalFeeFloat        float64                 `alipay:"to=2~5" validate:"min:2|max:5"`         // 订单总金额，单位为分，详见支付金额
	TestOrderDetailPtr   *TestOrderDetailPtr     `alipay:"required" validate:"required"`          // 商品详细描述
	TestOrderDetailSlice []*TestOrderDetailSlice `alipay:"required" validate:"required"`          // 商品详细描述
}

type TestOrderDetailPtr struct {
	TmpTest3  *TmpTest3 `alipay:"required" validate:"required"`
	GoodsName string    `alipay:"to=1~2" validate:"minLen:1|maxLen:2"`
}

type TestOrderDetailSlice struct {
	TmpTest3   *TmpTest3 `alipay:"required" validate:"required"`
	GoodsName  string    `alipay:"required" validate:"required"`
	BuyerNames []string  `alipay:"required" validate:"required"`
}

type TmpTest3 struct {
	Name string `alipay:"required" validate:"required"`
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
		AppName:              "集美测试",
		TotalFeeFloat:        2,
		TestOrderDetailPtr:   testOrderDetailPtr,
		TestOrderDetailSlice: testOrderDetails,
	}
	err := ValidateStruct(u, "alipay")
	if err == nil {
		return
	}
	sureMsg := `"TestOrder.TestOrderDetailPtr.GoodsName" input "玻尿酸" strLength more than 2; "TestOrder-0.TestOrderDetailSlice.GoodsName" input "" is required; "TestOrder-1.TestOrderDetailSlice.BuyerNames" input "" is required; "TestOrder-2.TestOrderDetailSlice.TmpTest3" input "" is required; "TestOrder-2.TestOrderDetailSlice.BuyerNames" input "" is required; "TestOrder-3.TestOrderDetailSlice.BuyerNames" input "" is required`
	if !equal(err.Error(), sureMsg) {
		t.Error(noEqErr)
	}
}

// func TestProtoPb1(t *testing.T) {
// 	u := &test.User{
// 		M: &test.Man{
// 			Name: "xue",
// 			Age:  0,
// 		},
// 		Phone: "13540042615",
// 	}
// 	err := ValidateStruct(u)
// 	if err == nil {
// 		return
// 	}

// 	suerMsg := `valid: "he" is not exist, You can call SetValidFn`
// 	if !equal(err.Error(), suerMsg) {
// 		t.Error(noEqErr)
// 	}
// }

// func TestValidateOrder(t *testing.T) {
// 	testOrderDetailPtr := &TestOrderDetailPtr{
// 		TmpTest3:  &TmpTest3{Name: "测试"},
// 		GoodsName: "玻尿酸",
// 	}
// 	// testOrderDetailPtr = nil

// 	testOrderDetails := []*TestOrderDetailSlice{
// 		{TmpTest3: &TmpTest3{Name: "测试1"}, BuyerNames: []string{"test1", "hello2"}},
// 		{TmpTest3: &TmpTest3{Name: "测试2"}, GoodsName: "隆鼻"},
// 		{GoodsName: "丰胸"},
// 		{TmpTest3: &TmpTest3{Name: "测试4"}, GoodsName: "隆鼻"},
// 	}
// 	// testOrderDetails = nil

// 	u := &TestOrder{
// 		AppName:              "集美测试",
// 		TotalFeeFloat:        2,
// 		TestOrderDetailPtr:   testOrderDetailPtr,
// 		TestOrderDetailSlice: testOrderDetails,
// 	}
// 	validObj := validate.Struct(u)
// 	validObj.Validate()
// 	for _, err := range validObj.Errors {
// 		t.Log(err)
// 	}
// }

func TestGetJoinValidErrStr(t *testing.T) {
	t.Skip("GetJoinValidErrStr")
	res := GetJoinValidErrStr("User", "Name", "xue", "len is less than 3")
	if !equal(res, `"User.Name" input "xue" len is less than 3;`) {
		t.Error(noEqErr)
	}
}

func BenchmarkValid(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		_ = ValidateStruct(&u, "alipay")
	}

	// BenchmarkValid-8          162354              7325 ns/op            4635 B/op        108 allocs/op
	// BenchmarkValid-8          160914              7321 ns/op            4635 B/op        108 allocs/op
	// BenchmarkValid-8          162373              7303 ns/op            4635 B/op        108 allocs/op
	// BenchmarkValid-8          160164              7414 ns/op            4635 B/op        108 allocs/op
	// BenchmarkValid-8          161552              7494 ns/op            4635 B/op        108 allocs/op
}

// func BenchmarkValidate(b *testing.B) {
// 	testOrderDetailPtr := &TestOrderDetailPtr{
// 		TmpTest3:  &TmpTest3{Name: "测试"},
// 		GoodsName: "玻尿酸",
// 	}
// 	// testOrderDetailPtr = nil

// 	testOrderDetails := []*TestOrderDetailSlice{
// 		{TmpTest3: &TmpTest3{Name: "测试1"}, BuyerNames: []string{"test1", "hello2"}},
// 		{TmpTest3: &TmpTest3{Name: "测试2"}, GoodsName: "隆鼻"},
// 		{GoodsName: "丰胸"},
// 		{TmpTest3: &TmpTest3{Name: "测试4"}, GoodsName: "隆鼻"},
// 	}
// 	// testOrderDetails = nil

// 	u := &TestOrder{
// 		AppName:              "集美测试",
// 		TotalFeeFloat:        2,
// 		TestOrderDetailPtr:   testOrderDetailPtr,
// 		TestOrderDetailSlice: testOrderDetails,
// 	}
// 	for i := 0; i < b.N; i++ {
// 		validObj := validate.Struct(u)
// 		validObj.Validate()
// 		_ = validObj.Errors.Error()
// 	}

// 	// BenchmarkValidate-8        30902             38716 ns/op           33267 B/op        430 allocs/op
// 	// BenchmarkValidate-8        30868             38861 ns/op           33263 B/op        430 allocs/op
// 	// BenchmarkValidate-8        29767             39055 ns/op           33262 B/op        430 allocs/op
// 	// BenchmarkValidate-8        30717             38774 ns/op           33268 B/op        430 allocs/op
// 	// BenchmarkValidate-8        31004             38489 ns/op           33264 B/op        430 allocs/op
// }

func BenchmarkValidIf(b *testing.B) {
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

	vFn := func(info *TestOrder) {
		// 说明: 写入的内容为随意写
		errBuf := new(strings.Builder)
		if len(info.AppName) >= 2 {
			errBuf.WriteString("info.AppName len is short \n")
		}

		if len(info.AppName) <= 10 {
			errBuf.WriteString("info.AppName len is long \n")
		}

		if info.TotalFeeFloat >= 2 {
			errBuf.WriteString("info.TotalFeeFloat should more than 2 \n")
		}

		if info.TotalFeeFloat <= 5 {
			errBuf.WriteString("info.TotalFeeFloat should more than 2 \n")
		}

		if info.TestOrderDetailPtr.TmpTest3 == nil {
			errBuf.WriteString("info.TestOrderDetailPtr.TmpTest3 is nil \n")
		} else {
			if info.TestOrderDetailPtr.TmpTest3.Name == "" {
				errBuf.WriteString("info.TestOrderDetailPtr.TmpTest3.Name is null \n")
			}
		}

		if len(info.TestOrderDetailPtr.GoodsName) >= 2 {
			errBuf.WriteString("info.TestOrderDetailPtr.GoodsName is null \n")
		}

		if len(info.TestOrderDetailPtr.GoodsName) <= 5 {
			errBuf.WriteString("info.TestOrderDetailPtr.GoodsName is null \n")
		}

		if info.TestOrderDetailSlice == nil {
			errBuf.WriteString("info.TestOrderDetailSlice is null \n")
		} else {
			for _, v := range info.TestOrderDetailSlice {
				if v.TmpTest3 == nil {
					errBuf.WriteString("info.TestOrderDetailPtr.GoodsName is null \n")
				} else {
					if v.TmpTest3.Name == "" {
						errBuf.WriteString("info.TestOrderDetailPtr.GoodsName is null \n")
					}
				}

				if len(v.BuyerNames) == 0 {
					errBuf.WriteString("info.TestOrderDetailPtr.GoodsName is null \n")
				}

				if v.GoodsName == "" {
					errBuf.WriteString("info.TestOrderDetailPtr.GoodsName is null \n")
				}
			}
		}
	}
	for i := 0; i < b.N; i++ {
		vFn(u)
	}

	// BenchmarkValidIf-8       4908489               242.1 ns/op          1232 B/op          5 allocs/op
	// BenchmarkValidIf-8       4932530               240.9 ns/op          1232 B/op          5 allocs/op
	// BenchmarkValidIf-8       4938175               252.1 ns/op          1232 B/op          5 allocs/op
	// BenchmarkValidIf-8       4774190               251.2 ns/op          1232 B/op          5 allocs/op
	// BenchmarkValidIf-8       4901599               253.2 ns/op          1232 B/op          5 allocs/op
}
