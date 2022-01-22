package valid

import (
	"strings"
	"testing"

	"gitee.com/xuesongtao/protoc-go-valid/test"
)

type TestOrder struct {
	AppName              string                  `alipay:"to=2~10"`  //应用名
	TotalFeeFloat        float64                 `alipay:"to=2~5"`   //订单总金额，单位为分，详见支付金额
	TestOrderDetailPtr   *TestOrderDetailPtr     `alipay:"required"` // 商品详细描述
	TestOrderDetailSlice []*TestOrderDetailSlice `alipay:"required"` // 商品详细描述
}

type TestOrderDetailPtr struct {
	TmpTest3  *TmpTest3 `alipay:"required"`
	GoodsName string    `alipay:"to=1~2"`
}

type TestOrderDetailSlice struct {
	TmpTest3   *TmpTest3 `alipay:"required"`
	GoodsName  string    `alipay:"required"`
	BuyerNames []string  `alipay:"required"`
}

type TmpTest3 struct {
	Name string `alipay:"required"`
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
	t.Log(ValidateStruct(u, "alipay"))
}

func TestProtoPb1(t *testing.T) {
	u := &test.User{
		M: &test.Man{
			Name: "xue",
			Age:  0,
		},
		Phone: "13540042615",
	}
	t.Log(ValidateStruct(u))
}

func TestProtoPb2(t *testing.T) {
	m := &test.Man{
		Name: "xue",
		Age:  0,
	}
	t.Log(ValidateStruct(m))
}

func TestTmp(t *testing.T) {
	t.Log(GetJoinValidErrStr("ab", "cd", "1", "is not ok", "tt"))
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
}

func BenchmarkIfValid(b *testing.B) {
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
}
