package valid

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

type TestOrder struct {
	AppName string `alipay:"to=5~10" validate:"min=5,max=10"` // 应用名
	// TotalFeeFloat        float64                 `alipay:"to=2~5" validate:"min:2|max:5"` // 订单总金额，单位为分，详见支付金额
	TotalFeeFloat        float64                 `alipay:"to=2~5" validate:"min=2,max=5"` // 订单总金额，单位为分，详见支付金额
	TestOrderDetailPtr   *TestOrderDetailPtr     `alipay:"required" validate:"required"`  // 商品详细描述
	TestOrderDetailSlice []*TestOrderDetailSlice `alipay:"required" validate:"required"`  // 商品详细描述
}

type TestOrderDetailPtr struct {
	TmpTest3  *TmpTest3 `alipay:"required" validate:"required"`
	GoodsName string    `alipay:"to=1~2" validate:"min=1,max=2"`
}

type TestOrderDetailSlice struct {
	TmpTest3   *TmpTest3 `alipay:"required" validate:"required"`
	GoodsName  string    `alipay:"required" validate:"required"`
	BuyerNames []string  `alipay:"required" validate:"required"`
}

type TmpTest3 struct {
	Name string `alipay:"required" valid:"required" validate:"required"`
}

func BenchmarkStringSplitValid(b *testing.B) {
	var s []string
	for i := 0; i < b.N; i++ {
		s = ValidNamesSplit("required,phone,test")
	}
	_ = s
}
func BenchmarkStringSplitValid1(b *testing.B) {
	var s []string
	for i := 0; i < b.N; i++ {
		s = ValidNamesSplit("required,phone,test|'test'")
	}
	_ = s
}

func BenchmarkStringSplit(b *testing.B) {
	var s []string
	for i := 0; i < b.N; i++ {
		s = strings.Split("required,phone,test", ",")
	}
	_ = s
}

func BenchmarkReNoComplice(b *testing.B) {
	a := "123456"
	for i := 0; i < b.N; i++ {
		_, _ = regexp.MatchString(`\d+`, a)
	}
}

func BenchmarkReComplice(b *testing.B) {
	a := "123456"
	for i := 0; i < b.N; i++ {
		_ = IntRe.MatchString(a)
	}
}

func BenchmarkStrJoin1(b *testing.B) {
	s := ""
	for i := 0; i < 100; i++ {
		s += fmt.Sprint(i) + ","
	}
	_ = s
}
func BenchmarkStrJoin2(b *testing.B) {
	buf := newStrBuf()
	for i := 0; i < 100; i++ {
		buf.WriteString(fmt.Sprint(i) + ",")
	}
	_ = buf.String()
}

// go test -benchmem -run=^$ -bench ^BenchmarkValidateForValid gitee.com/xuesongtao/protoc-go-valid/valid -v -count=5

func BenchmarkValidateForValid(b *testing.B) {
	b.ResetTimer()
	type Users struct {
		Phone  string `valid:"required"`
		Passwd string `valid:"required,to=6~20"`
		Code   string `validate:"required,eq=6"`
	}

	users := &Users{
		Phone:  "1326654487",
		Passwd: "123",
		Code:   "123456",
	}

	for i := 0; i < b.N; i++ {
		_ = Struct(users)
	}

	// BenchmarkValidateForValid-8              1614230               744.7 ns/op           416 B/op          9 allocs/op
	// BenchmarkValidateForValid-8              1617544               739.9 ns/op           416 B/op          9 allocs/op
	// BenchmarkValidateForValid-8              1618682               740.3 ns/op           416 B/op          9 allocs/op
	// BenchmarkValidateForValid-8              1621915               739.2 ns/op           416 B/op          9 allocs/op
	// BenchmarkValidateForValid-8              1612825               739.7 ns/op           416 B/op          9 allocs/op
}

// go test -benchmem -run=^$ -bench ^BenchmarkComplexValid gitee.com/xuesongtao/protoc-go-valid/valid -v -count=5

func BenchmarkComplexValid(b *testing.B) {
	b.ResetTimer()
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
		_ = StructForFn(&u, nil, "alipay")
	}

	// BenchmarkComplexValid-8           209080              5954 ns/op            3977 B/op         66 allocs/op
	// BenchmarkComplexValid-8           205692              5717 ns/op            3977 B/op         66 allocs/op
	// BenchmarkComplexValid-8           205198              5700 ns/op            3977 B/op         66 allocs/op
	// BenchmarkComplexValid-8           206556              5813 ns/op            3977 B/op         66 allocs/op
	// BenchmarkComplexValid-8           205345              5756 ns/op            3977 B/op         66 allocs/op
}

func BenchmarkComplexValidIf(b *testing.B) {
	b.ResetTimer()
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
