package valid

import (
	"strings"
	"testing"
)

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
		_ = ValidateStruct(users)
	}

	// BenchmarkValidateForValid-8              1235048               959.1 ns/op           456 B/op         12 allocs/op
	// BenchmarkValidateForValid-8              1258246               950.5 ns/op           456 B/op         12 allocs/op
	// BenchmarkValidateForValid-8              1258429               952.3 ns/op           456 B/op         12 allocs/op
	// BenchmarkValidateForValid-8              1259571               952.8 ns/op           456 B/op         12 allocs/op
	// BenchmarkValidateForValid-8              1258640               951.5 ns/op           456 B/op         12 allocs/op
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
		_ = ValidateStruct(&u, "alipay")
	}

	// BenchmarkValid-8          162354              7325 ns/op            4635 B/op        108 allocs/op
	// BenchmarkValid-8          160914              7321 ns/op            4635 B/op        108 allocs/op
	// BenchmarkValid-8          162373              7303 ns/op            4635 B/op        108 allocs/op
	// BenchmarkValid-8          160164              7414 ns/op            4635 B/op        108 allocs/op
	// BenchmarkValid-8          161552              7494 ns/op            4635 B/op        108 allocs/op
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
