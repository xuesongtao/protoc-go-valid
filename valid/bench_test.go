package valid

import (
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gookit/validate"
)

func TestValidateForValid(t *testing.T) {
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
	validObj := NewVStruct()
	err := validObj.Valid(users)
	if err != nil {
		t.Error(err)
	}

	err = validObj.Valid(users)
	if err != nil {
		t.Error(err)
	}
}

func TestValidateForValidate(t *testing.T) {
	type Users struct {
		Phone  string `validate:"required"`
		Passwd string `validate:"required|minLen:6|maxLen:20"`
		Code   string `validate:"required|eq:6"`
	}

	users := &Users{
		Phone:  "1326654487",
		Passwd: "123",
		Code:   "123456",
	}

	validObj := validate.Struct(users)
	if !validObj.Validate() {
		t.Log(validObj.Errors)
	}
}

func TestValidateForValidator(t *testing.T) {
	type Users struct {
		Phone  string `validate:"required"`
		Passwd string `validate:"required,min=6,max=20"`
		Code   string `validate:"required,len=6"`
	}

	users := &Users{
		Phone:  "1326654487",
		Passwd: "123",
		Code:   "123456",
	}

	validObj := validator.New()
	err := validObj.Struct(users)
	if err != nil {
		t.Error(err)
	}

	err = validObj.Struct(users)
	if err != nil {
		t.Error(err)
	}

}

// go test -benchmem -run=^$ -bench ^BenchmarkValidateForValid gitee.com/xuesongtao/protoc-go-valid/valid -v -count=5

func BenchmarkValidateForValid(b *testing.B) {
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

func BenchmarkValidateForValidate(b *testing.B) {
	type Users struct {
		Phone  string `validate:"required"`
		Passwd string `validate:"required|minLen:6|maxLen:20"`
		Code   string `validate:"required|eq:6"`
	}

	users := &Users{
		Phone:  "1326654487",
		Passwd: "123",
		Code:   "123456",
	}

	for i := 0; i < b.N; i++ {
		_ = validate.Struct(users).Validate()
	}

	// BenchmarkValidateForValidate-8             51856             23022 ns/op           23087 B/op        160 allocs/op
	// BenchmarkValidateForValidate-8             52057             23147 ns/op           23087 B/op        160 allocs/op
	// BenchmarkValidateForValidate-8             51775             23082 ns/op           23085 B/op        160 allocs/op
	// BenchmarkValidateForValidate-8             47670             23004 ns/op           23087 B/op        160 allocs/op
	// BenchmarkValidateForValidate-8             51775             22992 ns/op           23086 B/op        160 allocs/op
}

func BenchmarkValidateForValidator(b *testing.B) {
	type Users struct {
		Phone  string `validate:"required"`
		Passwd string `validate:"required,min=6,max=20"`
		Code   string `validate:"required,eq=6"`
	}

	users := &Users{
		Phone:  "1326654487",
		Passwd: "123",
		Code:   "123456",
	}

	validObj := validator.New()
	for i := 0; i < b.N; i++ {
		_ = validObj.Struct(users)
	}

	// BenchmarkValidateForValidator-8          1962073               611.1 ns/op           472 B/op         12 allocs/op
	// BenchmarkValidateForValidator-8          1839142               621.7 ns/op           472 B/op         12 allocs/op
	// BenchmarkValidateForValidator-8          1965544               614.8 ns/op           472 B/op         12 allocs/op
	// BenchmarkValidateForValidator-8          1956302               620.4 ns/op           472 B/op         12 allocs/op
	// BenchmarkValidateForValidator-8          1944094               618.7 ns/op           472 B/op         12 allocs/op
}

// go test -benchmem -run=^$ -bench ^BenchmarkValid gitee.com/xuesongtao/protoc-go-valid/valid -v -count=5

func BenchmarkValid(b *testing.B) {
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

func BenchmarkValidate(b *testing.B) {
	b.ResetTimer()
	b.Skip()
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
	for i := 0; i < b.N; i++ {
		validObj := validate.Struct(u)
		validObj.Validate()
		_ = validObj.Errors.Error()
	}

	// BenchmarkValidate-8        30902             38716 ns/op           33267 B/op        430 allocs/op
	// BenchmarkValidate-8        30868             38861 ns/op           33263 B/op        430 allocs/op
	// BenchmarkValidate-8        29767             39055 ns/op           33262 B/op        430 allocs/op
	// BenchmarkValidate-8        30717             38774 ns/op           33268 B/op        430 allocs/op
	// BenchmarkValidate-8        31004             38489 ns/op           33264 B/op        430 allocs/op
}

func BenchmarkValidator(b *testing.B) {
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
	// TODO: 不支持验证切片结构体
	validObj := validator.New()
	for i := 0; i < b.N; i++ {
		_ = validObj.Struct(u)
	}

	// BenchmarkValidate-8        30902             38716 ns/op           33267 B/op        430 allocs/op
	// BenchmarkValidate-8        30868             38861 ns/op           33263 B/op        430 allocs/op
	// BenchmarkValidate-8        29767             39055 ns/op           33262 B/op        430 allocs/op
	// BenchmarkValidate-8        30717             38774 ns/op           33268 B/op        430 allocs/op
	// BenchmarkValidate-8        31004             38489 ns/op           33264 B/op        430 allocs/op
}

func BenchmarkValidIf(b *testing.B) {
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
