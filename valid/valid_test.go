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
	// t.Log(GetTimeFmt(DateTimeFmt))
	// t.Log(GetTimeFmt(DateFmt))
}

func TestGetTimeFmt(t *testing.T) {
	res := GetTimeFmt(DateFmt)
	if res != "2006-01-02" {
		t.Error("res is not ok")
	}
	res = GetTimeFmt(DateTimeFmt)
	if res != "2006-01-02 15:04:05" {
		t.Error("res is not ok")
	}
}

func TestValidManyStruct(t *testing.T) {
	type Tmp1 struct {
		Name string `valid:"required"`
	}

	type Tmp struct {
		Ip string `valid:"required,ipv4"`
		T  []Tmp1 `valid:"required"`
	}

	v := &Tmp{
		Ip: "256.12.22.4",
		T:  []Tmp1{{Name: ""}},
	}
	datas := append([]*Tmp{}, v, v)
	sureMsg := `"*valid.Tmp[0].Ip" input "256.12.22.4", explain: it is not ipv4; "*valid.Tmp[0].T[0].Name" input "", explain: it is required; "*valid.Tmp[1].Ip" input "256.12.22.4", explain: it is not ipv4; "*valid.Tmp[1].T[0].Name" input "", explain: it is required`
	err := ValidateStruct(datas)
	if !equal(err.Error(), sureMsg) {
		t.Error(noEqErr)
	}
}

func TestValidStructMap(t *testing.T) {
	t.Run("map value is struct", func(t *testing.T) {
		type Tmp struct {
			Name string `valid:"required"`
			Age  int    `valid:"required,to=1~100"`
		}
		tmp := &Tmp{
			Name: "测试json",
			Age:  101,
		}
		a := map[string]*Tmp{
			"1": tmp,
		}
		sureMsg := `"map[1].Age" input "101", explain: it is more than 100 num-size`
		err := Struct(a)
		if err != nil && !equal(err.Error(), sureMsg) {
			t.Error(err)
		}
	})

	t.Run("struct have map filed", func(t *testing.T) {
		type Tmp struct {
			Name     string          `valid:"required"`
			Age      int             `valid:"required,to=1~100"`
			Students map[string]*Tmp `valid:"exist"`
		}
		tmp := &Tmp{
			Name: "测试json",
			Age:  20,
			Students: map[string]*Tmp{
				"同学1": {
					Name: "同学1",
					Age:  -1,
				},
				"同学2": {
					Name: "同学2",
					Age:  10,
				},
				"同学3": {
					Name: "同学3",
					Age:  20,
				},
			},
		}
		sureMsg := `"Tmp.Students[同学1].Age" input "-1", explain: it is less than 1 num-size`
		err := Struct(tmp)
		if err != nil && !equal(err.Error(), sureMsg) {
			t.Error(err)
		}
	})
}

func TestValidManyStructRule(t *testing.T) {
	type Tmp1 struct {
		Name string
	}

	type Tmp struct {
		Ip string
		T  []Tmp1
	}
	rmap := map[interface{}]RM{
		// key 必须为 指针
		&Tmp{}:  NewRule().Set("Ip,T", Required).Set("Ip", GenValidKV(VIp, "", "ip 格式不正确")),
		&Tmp1{}: NewRule().Set("Name", GenValidKV(Required, "", "姓名必填")),
	}
	// t.Logf("rmap: %+v", rmap)
	v := &Tmp{
		Ip: "256.12.22.400",
		T:  []Tmp1{{Name: ""}, {Name: "2"}},
	}
	sureMsg := `"Tmp.Ip" input "256.12.22.400", 说明: ip 格式不正确; "Tmp.T[0].Name" input "", 说明: 姓名必填`
	err := NestedStructForRule(v, rmap)
	if err != nil && !equal(err.Error(), sureMsg) {
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
	sureMsg := `"TestOrder.AppName" input "测试", explain: it is less than 5 str-length; "TestOrder.TestOrderDetailPtr.GoodsName" input "玻尿酸", explain: it is more than 2 str-length; "TestOrder.TestOrderDetailSlice[0].GoodsName" input "", explain: it is required; "TestOrder.TestOrderDetailSlice[1].BuyerNames" input "", explain: it is required; "TestOrder.TestOrderDetailSlice[2].TmpTest3" input "", explain: it is required; "TestOrder.TestOrderDetailSlice[2].BuyerNames" input "", explain: it is required; "TestOrder.TestOrderDetailSlice[3].BuyerNames" input "", explain: it is required`
	if !equal(err.Error(), sureMsg) {
		t.Error(noEqErr)
	}
}

func TestProtoPb(t *testing.T) {
	type Man struct {
		Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty" valid:"required|姓名必填,to=1~3"` // 姓名 @tag valid:"required|姓名必填,to=1~3"
		Age  int32  `protobuf:"varint,2,opt,name=age,proto3" json:"age,omitempty" valid:"to=1~150"`              // 年龄 @tag valid:"to=1~150"
		Tmp  string `protobuf:"bytes,3,opt,name=tmp,proto3" json:"tmp,omitempty" valid:"he"`                     // 临时 @tag valid:"he"
	}
	type User struct {
		M     *Man   `protobuf:"bytes,1,opt,name=m,proto3" json:"m,omitempty" valid:"required"`      // 人 @tag valid:"required"
		Phone string `protobuf:"bytes,2,opt,name=phone,proto3" json:"phone,omitempty" valid:"phone"` // 手机 @tag valid:"phone"
	}
	
	u := &User{
		M: &Man{
			Name: "",
			Age:  0,
		},
		Phone: "13540042615",
	}
	err := ValidateStruct(u)
	if err == nil {
		return
	}

	suerMsg := `"User.M.Name" input "", 说明: 姓名必填; "User.M.Tmp" valid "he" is not exist, You can call SetValidFn`
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

func TestValidMap(t *testing.T) {
	t.Run("map[string]interface{}", func(t *testing.T) {
		testMap := map[string]interface{}{"name": "", "addr": "chendu", "age": -1}
		rm := NewRule().Set("name,addr", Required).Set("age", GenValidKV(VTo, "1~100", "年龄必须在1~100"))
		err := Map(testMap, rm)
		sureMsg := `"name" input "", explain: it is required; "age" input "-1", explain: 年龄必须在1~100`
		if err != nil && !equal(err.Error(), sureMsg) {
			t.Error(err)
		}
	})

	t.Run("map[string]string", func(t *testing.T) {
		testMap := map[string]string{"name": "", "addr": "chendu"}
		rm := NewRule().Set("name,addr", Required)
		err := Map(testMap, rm)
		sureMsg := `"map[name]" input "", explain: it is required`
		if err != nil && !equal(err.Error(), sureMsg) {
			t.Error(err)
		}
	})

	t.Run("map[string]int", func(t *testing.T) {
		testMap := map[string]int{"no1": 1, "no2": 2}
		rm := NewRule().Set("no1", GenValidKV(VEq, "2", "no1必须等于1")).Set("no2", GenValidKV(VEq, "2", "no2必须等于2"))
		err := Map(testMap, rm)
		sureMsg := `"map[no1]" input "1", 说明: no1必须等于1`
		if err != nil && !equal(err.Error(), sureMsg) {
			t.Error(err)
		}
	})

	t.Run("[]map[string]string", func(t *testing.T) {
		testMap := []map[string]string{{"name": "", "addr": "chendu"}, {"name": "test", "addr": ""}}
		rm := NewRule().Set("name,addr", Required)
		err := Map(testMap, rm)
		sureMsg := `"[0]map[name]" input "", explain: it is required; "[1]map[addr]" input "", explain: it is required`
		if err != nil && !equal(err.Error(), sureMsg) {
			t.Error(err)
		}
	})
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
