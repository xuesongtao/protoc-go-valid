package Validate

import (
	"protoc-go-cjvalid/test"
	"testing"
	"time"
)

type TestOrder struct {
	AppName              string                  `alipay:"to=2" wechat:"required"`     //应用名
	GoodsName            string                  `alipay:"required" wechat:"required"` //商品名
	TestOrderDetailPtr   *TestOrderDetailPtr     `alipay:"required" wechat:"required"` // 商品详细描述
	TestOrderDetailSlice []*TestOrderDetailSlice `alipay:"required" wechat:"required"` // 商品详细描述
	OutTradeNo           string                  `alipay:"either=1" wechat:"required"` // 商户订单号
	OrderNo              string                  `alipay:"either=1" wechat:"required"` // 商户订单号
	TotalFee             int8                    `alipay:"required,to=1~10"`           //订单总金额，单位为分，详见支付金额
	TotalFeeFloat        float64                 `alipay:"required,to=2~5"`            //订单总金额，单位为分，详见支付金额
	TimeExpire           time.Time               `alipay:"required"`                   //交易过期时间
}

type TestOrderDetailPtr struct {
	TmpTest3  *TmpTest3 `alipay:"required"`
	GoodsId   string    `alipay:"required"`
	GoodsName string    `alipay:"required"`
	Quantity  int
	Price     int
}

type TestOrderDetailSlice struct {
	TmpTest3  *TmpTest3 `alipay:"required"`
	GoodsId   string    `alipay:"required"`
	GoodsName string    `alipay:"required"`
	Quantity  int       `alipay:"required"`
	Price     int       `alipay:"required"`
}

type TmpTest3 struct {
	Name string `alipay:"required"`
}

func TestCjValid(t *testing.T) {
	testOrderDetailPtr := &TestOrderDetailPtr{
		TmpTest3:  &TmpTest3{Name: "测试1"},
		GoodsId:   "10001",
		GoodsName: "玻尿酸",
	}
	testOrderDetailPtr = nil

	testOrderDetails := []*TestOrderDetailSlice{
		{TmpTest3: &TmpTest3{Name: "测试2"}, GoodsId: "10002", GoodsName: "隆鼻"},
		{TmpTest3: &TmpTest3{Name: "测试3"}, GoodsId: "10003", GoodsName: "丰胸"},
	}
	testOrderDetails = nil

	u := &TestOrder{
		AppName:   "集12",
		GoodsName: "集美",
		// OutTradeNo:           "MZ1234567890",
		// OrderNo:              "123",
		TestOrderDetailPtr:   testOrderDetailPtr,
		TestOrderDetailSlice: testOrderDetails,
		TotalFee:             100,
		TotalFeeFloat:        25,
	}
	t.Log(NewVStruct("alipay").Validate(u).GetErrMsg())
}

func TestTmp(t *testing.T) {
	u := &test.User{
		M: &test.Man{
			Name: "xue",
			Age:  0,
		},
		Phone: "",
	}
	t.Log(ValidateStruct(u))
}

func TestTmp1(t *testing.T) {
	m := &test.Man{
		Name: "xue",
		Age:  0,
	}
	t.Log(ValidateStruct(m))
}
