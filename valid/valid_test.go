package valid

import (
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

func TestEither(t *testing.T) {
	type Tmp struct {
		Either1 int32 `valid:"either=1"`
		Either2 int32 `valid:"either=1"`
	}
	v := &Tmp{}
	t.Log(ValidateStruct(v))
}

func TestIn(t *testing.T) {
	type Tmp struct {
		SelectNum int32  `valid:"in=(1/2/3/4)"`
		SelectStr string `valid:"in=(a/b/c/d)"`
	}
	v := &Tmp{SelectNum: 1, SelectStr: "a"}
	t.Log(ValidateStruct(v))
}

func TestPhone(t *testing.T) {
	type Tmp struct {
		Phone string `valid:"phone"`
	}
	v := &Tmp{Phone: "12344"}
	t.Log(ValidateStruct(v))
}

func TestEmail(t *testing.T) {
	type Tmp struct {
		Email string `valid:"email"`
	}
	v := &Tmp{Email: "xuesongtao512@qq.com"}
	t.Log(ValidateStruct(v))
}

func TestIdCard(t *testing.T) {
	type Tmp struct {
		IDCard string `valid:"idcard"`
	}
	v := &Tmp{IDCard: "511321"}
	t.Log(ValidateStruct(v))
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
